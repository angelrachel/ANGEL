from __future__ import annotations

import pytest

from src.evidence.normalize import normalize_http_evidence


def test_http_evidence_normalization_is_redacted_and_deterministic() -> None:
    kwargs = {
        "method": "get",
        "url": "https://example.test/api",
        "request_headers": {"Authorization": "Bearer secret", "Accept": "application/json"},
        "response_status": 200,
        "response_headers": {"Set-Cookie": "sid=secret", "Content-Type": "application/json"},
        "body": '{"ok":true}',
    }
    first, first_hash = normalize_http_evidence(**kwargs)
    second, second_hash = normalize_http_evidence(**kwargs)
    assert first == second and first_hash == second_hash
    assert first["request"]["headers"]["Authorization"] == "[REDACTED]"
    assert len(first_hash) == 64


def test_http_evidence_validation() -> None:
    with pytest.raises(ValueError):
        normalize_http_evidence(
            method="TRACE",
            url="https://example.test",
            request_headers={},
            response_status=200,
            response_headers={},
            body="",
        )
