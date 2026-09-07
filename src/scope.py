"""Fail-closed scope and request authorization policy."""

from __future__ import annotations

import time
from dataclasses import dataclass, field
from ipaddress import ip_address
from urllib.parse import urlparse


class ScopeViolation(ValueError):  # noqa: N818
    """Raised when a request is outside the declared engagement scope."""


@dataclass
class ScopePolicy:
    hosts: set[str]
    paths: tuple[str, ...] = ("/",)
    expires_at: int | None = None
    max_requests_per_minute: int = 60
    block_private_ips: bool = True
    _requests: list[float] = field(default_factory=list, init=False, repr=False)

    def __post_init__(self) -> None:
        self.hosts = {host.lower().rstrip(".") for host in self.hosts if host.strip()}
        if not self.hosts:
            raise ValueError("scope requires at least one host")
        if self.max_requests_per_minute < 1:
            raise ValueError("request limit must be positive")

    def authorize(self, url: str, *, method: str = "GET", now: int | None = None) -> None:
        parsed = urlparse(url)
        host = (parsed.hostname or "").lower().rstrip(".")
        if parsed.scheme not in {"http", "https"} or not host:
            raise ScopeViolation("only HTTP(S) URLs with a host are allowed")
        if not self._host_allowed(host):
            raise ScopeViolation("host is outside scope")
        if not any(parsed.path.startswith(path.rstrip("*") or "/") for path in self.paths):
            raise ScopeViolation("path is outside scope")
        if self.block_private_ips:
            try:
                address = ip_address(host)
            except ValueError:
                address = None
            if address is not None and (address.is_private or address.is_loopback or address.is_link_local):
                raise ScopeViolation("private or link-local address is blocked")
        current = int(time.time()) if now is None else now
        if self.expires_at is not None and current > self.expires_at:
            raise ScopeViolation("scope has expired")
        if method.upper() not in {"GET", "HEAD", "OPTIONS", "POST", "PUT", "PATCH", "DELETE"}:
            raise ScopeViolation("HTTP method is not allowed")
        self._consume(current)

    def authorize_redirect(self, location: str, *, now: int | None = None) -> None:
        self.authorize(location, method="GET", now=now)

    def _host_allowed(self, host: str) -> bool:
        return any(host == allowed or host.endswith("." + allowed) for allowed in self.hosts)

    def _consume(self, now: int) -> None:
        cutoff = now - 60
        self._requests = [stamp for stamp in self._requests if stamp > cutoff]
        if len(self._requests) >= self.max_requests_per_minute:
            raise ScopeViolation("scope rate limit exceeded")
        self._requests.append(float(now))


__all__ = ["ScopePolicy", "ScopeViolation"]
