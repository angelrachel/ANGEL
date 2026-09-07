"""Small fixed-window rate limiter for the operator API."""

from __future__ import annotations

import threading
import time
from collections import defaultdict, deque


class RateLimitError(RuntimeError):
    """Raised when a caller exceeds its configured request budget."""


class RateLimiter:
    def __init__(self, max_requests: int = 120, window_seconds: int = 60) -> None:
        if max_requests < 1 or window_seconds < 1:
            raise ValueError("rate limit values must be positive")
        self.max_requests = max_requests
        self.window_seconds = window_seconds
        self._requests: dict[str, deque[float]] = defaultdict(deque)
        self._lock = threading.Lock()

    def check(self, key: str, now: float | None = None) -> None:
        if not key.strip():
            raise ValueError("rate limit key is required")
        current = time.monotonic() if now is None else now
        cutoff = current - self.window_seconds
        with self._lock:
            bucket = self._requests[key]
            while bucket and bucket[0] <= cutoff:
                bucket.popleft()
            if len(bucket) >= self.max_requests:
                raise RateLimitError("rate limit exceeded")
            bucket.append(current)


__all__ = ["RateLimitError", "RateLimiter"]
