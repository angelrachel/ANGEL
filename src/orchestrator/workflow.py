"""Deterministic workflow primitives for safe, bounded orchestration."""

from __future__ import annotations

from collections.abc import Callable
from concurrent.futures import ThreadPoolExecutor
from concurrent.futures import TimeoutError as FutureTimeout
from dataclasses import dataclass, field
from enum import StrEnum
from time import monotonic
from typing import Generic, TypeVar

T = TypeVar("T")


class WorkflowError(RuntimeError):
    """Raised when workflow execution cannot proceed."""


class WorkflowState(StrEnum):
    PENDING = "pending"
    RUNNING = "running"
    SUCCEEDED = "succeeded"
    FAILED = "failed"
    TIMED_OUT = "timed_out"
    BLOCKED = "blocked"


@dataclass(frozen=True)
class RetryPolicy:
    max_attempts: int = 1
    timeout_seconds: float = 30.0

    def __post_init__(self) -> None:
        if self.max_attempts < 1:
            raise ValueError("max_attempts must be positive")
        if self.timeout_seconds <= 0:
            raise ValueError("timeout_seconds must be positive")


@dataclass
class CircuitBreaker:
    failure_threshold: int = 3
    reset_timeout: float = 60.0
    failures: int = field(default=0, init=False)
    opened_at: float | None = field(default=None, init=False)

    def __post_init__(self) -> None:
        if self.failure_threshold < 1 or self.reset_timeout <= 0:
            raise ValueError("invalid circuit breaker limits")

    @property
    def is_open(self) -> bool:
        if self.opened_at is None:
            return False
        if monotonic() - self.opened_at >= self.reset_timeout:
            self.opened_at = None
            self.failures = 0
            return False
        return True

    def record_success(self) -> None:
        self.failures = 0
        self.opened_at = None

    def record_failure(self) -> None:
        self.failures += 1
        if self.failures >= self.failure_threshold:
            self.opened_at = monotonic()


@dataclass(frozen=True)
class WorkflowResult(Generic[T]):
    state: WorkflowState
    attempts: int
    value: T | None = None
    error: str | None = None


class WorkflowRunner(Generic[T]):
    """Runs one bounded callable and never executes more than max_attempts."""

    def __init__(self, breaker: CircuitBreaker | None = None) -> None:
        self.breaker = breaker or CircuitBreaker()

    def run(self, operation: Callable[[], T], policy: RetryPolicy | None = None) -> WorkflowResult[T]:
        selected = policy or RetryPolicy()
        if self.breaker.is_open:
            return WorkflowResult(WorkflowState.BLOCKED, 0, error="circuit breaker is open")
        last_error = "workflow failed"
        for attempt in range(1, selected.max_attempts + 1):
            executor = ThreadPoolExecutor(max_workers=1)
            future = executor.submit(operation)
            try:
                value = future.result(timeout=selected.timeout_seconds)
            except FutureTimeout:
                future.cancel()
                last_error = "workflow timed out"
                self.breaker.record_failure()
                state = WorkflowState.TIMED_OUT
            except Exception as exc:  # noqa: BLE001 - workflow boundary normalizes failures
                last_error = str(exc) or exc.__class__.__name__
                self.breaker.record_failure()
                state = WorkflowState.FAILED
            else:
                self.breaker.record_success()
                executor.shutdown(wait=True)
                return WorkflowResult(WorkflowState.SUCCEEDED, attempt, value=value)
            finally:
                executor.shutdown(wait=False, cancel_futures=True)
            if self.breaker.is_open:
                break
        return WorkflowResult(state, attempt, error=last_error)


__all__ = ["CircuitBreaker", "RetryPolicy", "WorkflowError", "WorkflowResult", "WorkflowRunner", "WorkflowState"]
