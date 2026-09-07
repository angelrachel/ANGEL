from __future__ import annotations

import pytest

from src.api_rate_limit import RateLimiter, RateLimitError


def test_rate_limiter_enforces_window() -> None:
    limiter = RateLimiter(max_requests=2, window_seconds=10)
    limiter.check("client", now=100)
    limiter.check("client", now=101)
    with pytest.raises(RateLimitError):
        limiter.check("client", now=102)
    limiter.check("client", now=111)


def test_rate_limiter_validates_configuration() -> None:
    with pytest.raises(ValueError):
        RateLimiter(max_requests=0)
    with pytest.raises(ValueError):
        RateLimiter().check("", now=1)
