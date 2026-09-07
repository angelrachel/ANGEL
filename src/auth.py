"""Minimal role-based authentication for the ANGEL operator API."""

from __future__ import annotations

import base64
import hashlib
import hmac
import json
import os
import time
from dataclasses import dataclass
from typing import Any


class AuthError(ValueError):
    """Raised when an operator token is invalid or insufficient."""


@dataclass(frozen=True)
class Principal:
    subject: str
    role: str
    expires_at: int

    def can(self, permission: str) -> bool:
        permissions = {
            "admin": {"read", "write", "queue", "evidence", "report"},
            "operator": {"read", "write", "queue", "evidence", "report"},
            "auditor": {"read", "evidence", "report"},
            "viewer": {"read", "report"},
        }
        return permission in permissions.get(self.role, set())


def _secret() -> bytes:
    return os.environ.get("ANGEL_AUTH_SECRET", "development-auth-secret").encode()


def issue_token(subject: str, role: str, ttl: int = 3600, now: int | None = None) -> str:
    if role not in {"admin", "operator", "auditor", "viewer"}:
        raise AuthError("unknown role")
    issued = int(time.time()) if now is None else now
    payload = {"sub": subject, "role": role, "exp": issued + ttl}
    encoded = base64.urlsafe_b64encode(
        json.dumps(payload, sort_keys=True, separators=(",", ":")).encode()
    ).decode().rstrip("=")
    signature = hmac.new(_secret(), encoded.encode(), hashlib.sha256).hexdigest()
    return f"{encoded}.{signature}"


def authenticate(token: str, now: int | None = None) -> Principal:
    try:
        encoded, signature = token.split(".", 1)
        expected = hmac.new(_secret(), encoded.encode(), hashlib.sha256).hexdigest()
        if not hmac.compare_digest(expected, signature):
            raise AuthError("invalid token signature")
        padded = encoded + "=" * (-len(encoded) % 4)
        payload: dict[str, Any] = json.loads(base64.urlsafe_b64decode(padded.encode()))
        principal = Principal(str(payload["sub"]), str(payload["role"]), int(payload["exp"]))
    except (ValueError, KeyError, TypeError, json.JSONDecodeError) as exc:
        raise AuthError("malformed token") from exc
    current = int(time.time()) if now is None else now
    if principal.expires_at <= current:
        raise AuthError("token expired")
    if not principal.subject or not principal.can("read"):
        raise AuthError("invalid principal")
    return principal


__all__ = ["AuthError", "Principal", "authenticate", "issue_token"]
