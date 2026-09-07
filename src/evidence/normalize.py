"""Normalize captured HTTP evidence before persistence."""

from __future__ import annotations

import hashlib
import json
from typing import Any

from .chain import redact

SENSITIVE_HEADERS = frozenset({"authorization", "cookie", "set-cookie", "x-api-key"})


def normalize_http_evidence(
    *,
    method: str,
    url: str,
    request_headers: dict[str, str],
    response_status: int,
    response_headers: dict[str, str],
    body: str,
) -> tuple[dict[str, Any], str]:
    if method.upper() not in {"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}:
        raise ValueError("unsupported HTTP method")
    if not 100 <= response_status <= 599 or not url.strip():
        raise ValueError("invalid HTTP evidence metadata")
    def clean(headers: dict[str, str]) -> dict[str, str]:
        return {
            key: ("[REDACTED]" if key.lower() in SENSITIVE_HEADERS else str(value))
            for key, value in sorted(headers.items())
        }
    evidence = {
        "request": {"method": method.upper(), "url": url, "headers": clean(request_headers)},
        "response": {"status": response_status, "headers": clean(response_headers), "body": body[:65536]},
    }
    safe = json.loads(redact(json.dumps(evidence, sort_keys=True, separators=(",", ":"))))
    record_hash = hashlib.sha256(json.dumps(safe, sort_keys=True, separators=(",", ":")).encode()).hexdigest()
    return safe, record_hash


__all__ = ["normalize_http_evidence"]
