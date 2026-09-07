from __future__ import annotations

import time

from src.orchestrator.workflow import CircuitBreaker, RetryPolicy, WorkflowRunner, WorkflowState


def test_workflow_success_and_retry() -> None:
    calls = 0

    def operation() -> str:
        nonlocal calls
        calls += 1
        if calls == 1:
            raise ValueError("transient")
        return "done"

    result = WorkflowRunner().run(operation, RetryPolicy(max_attempts=2, timeout_seconds=1))
    assert result.state == WorkflowState.SUCCEEDED
    assert result.attempts == 2
    assert result.value == "done"


def test_workflow_timeout_and_open_breaker() -> None:
    breaker = CircuitBreaker(failure_threshold=1, reset_timeout=60)
    runner = WorkflowRunner(breaker)
    result = runner.run(lambda: time.sleep(0.05), RetryPolicy(timeout_seconds=0.001))
    assert result.state == WorkflowState.TIMED_OUT
    blocked = runner.run(lambda: "not-run")
    assert blocked.state == WorkflowState.BLOCKED
    assert blocked.attempts == 0


def test_retry_policy_and_breaker_limits() -> None:
    try:
        RetryPolicy(max_attempts=0)
        raise AssertionError("expected validation")
    except ValueError:
        pass
    try:
        CircuitBreaker(failure_threshold=0)
        raise AssertionError("expected validation")
    except ValueError:
        pass
