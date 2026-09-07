"""Cooperative cancellation context for bounded workflows."""

from __future__ import annotations

import threading


class CancellationError(RuntimeError):
    """Raised when a cooperative operation observes cancellation."""


class CancellationContext:
    def __init__(self) -> None:
        self._event = threading.Event()

    def cancel(self) -> None:
        self._event.set()

    @property
    def cancelled(self) -> bool:
        return self._event.is_set()

    def raise_if_cancelled(self) -> None:
        if self.cancelled:
            raise CancellationError("operation cancelled")

    def wait(self, timeout: float | None = None) -> bool:
        return self._event.wait(timeout)


__all__ = ["CancellationContext", "CancellationError"]
