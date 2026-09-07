from __future__ import annotations

import json

import pytest

from src.evidence.chain import EvidenceChain, redact
from src.reporting import Finding, Report
from src.scope import ScopePolicy, ScopeViolation


def test_scope_allows_in_scope_subdomain_and_blocks_out_of_scope() -> None:
    policy = ScopePolicy({"example.test"}, paths=("/api/",), max_requests_per_minute=2)
    policy.authorize("https://api.example.test/api/status", now=100)
    with pytest.raises(ScopeViolation):
        policy.authorize("https://other.test/api/status", now=100)
    with pytest.raises(ScopeViolation):
        policy.authorize("https://api.example.test/admin", now=100)


def test_scope_blocks_private_expiry_and_rate_limit() -> None:
    policy = ScopePolicy({"127.0.0.1"}, max_requests_per_minute=1, expires_at=100)
    with pytest.raises(ScopeViolation):
        policy.authorize("http://127.0.0.1/", now=99)
    public = ScopePolicy({"example.test"}, max_requests_per_minute=1)
    public.authorize("https://example.test/", now=100)
    with pytest.raises(ScopeViolation):
        public.authorize("https://example.test/", now=100)
    expired = ScopePolicy({"example.test"}, expires_at=100)
    with pytest.raises(ScopeViolation):
        expired.authorize("https://example.test/", now=101)


def test_evidence_redacts_secrets_and_verifies_chain() -> None:
    assert "[REDACTED]" in redact("Authorization: Bearer abc123 Cookie: sid=xyz")
    chain = EvidenceChain()
    first = chain.append(
        "request", "tester", {"authorization": "Bearer abc", "url": "https://example.test"}, created_at=10
    )
    chain.append("response", "tester", {"status": 200}, created_at=11)
    assert first.sequence == 1
    assert chain.verify()
    chain.records[1] = chain.records[1].__class__(**{**chain.records[1].__dict__, "actor": "tampered"})
    assert not chain.verify()


def test_report_exports() -> None:
    report = Report("ANGEL Report", "authorized-test", executive_summary="One confirmed finding.")
    report.add(
        Finding(
            "Cross-tenant access",
            "high",
            "high",
            "api.example.test",
            "Object access mismatch",
            "Tenant boundary is not enforced",
            ["Create canary", "Replay paired request"],
            "Enforce server-side ownership",
            ["ev-1"],
        )
    )
    markdown = report.to_markdown()
    assert "Cross-tenant access" in markdown
    assert json.loads(report.to_json())["findings"][0]["severity"] == "high"


def test_report_validation_rejects_malformed_findings() -> None:
    with pytest.raises(ValueError, match="severity"):
        Finding("Title", "urgent", "high", "asset", "Summary", "Impact", ["step"], "Fix")
    with pytest.raises(ValueError, match="confidence"):
        Finding("Title", "high", "certain", "asset", "Summary", "Impact", ["step"], "Fix")
    with pytest.raises(ValueError, match="engagement"):
        Report("Title", "")
