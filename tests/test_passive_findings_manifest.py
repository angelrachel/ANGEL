from __future__ import annotations

import pytest

from src.evidence.chain import EvidenceChain
from src.evidence.manifest import ManifestError, build_manifest, verify_manifest
from src.osint.passive_checks import PassiveSignal
from src.passive_findings import signals_to_findings


def test_passive_signals_become_deduplicated_informational_findings() -> None:
    signals = [
        PassiveSignal("framework", "django", "captured-response"),
        PassiveSignal("framework", "django", "captured-response"),
    ]
    findings = signals_to_findings(signals, "api.example.test")
    assert len(findings) == 1
    assert findings[0].severity == "info"
    assert findings[0].evidence_refs[0].startswith("passive:")


def test_manifest_round_trip_and_tamper_detection() -> None:
    chain = EvidenceChain()
    chain.append("request", "auditor", {"status": 200}, created_at=1)
    manifest = build_manifest(chain, engagement_id="eng-1")
    assert verify_manifest(chain, manifest)
    manifest["record_count"] = 99
    assert not verify_manifest(chain, manifest)


def test_manifest_requires_engagement_id() -> None:
    with pytest.raises(ManifestError, match="engagement_id"):
        build_manifest(EvidenceChain(), engagement_id="")
    with pytest.raises(ManifestError, match="malformed"):
        verify_manifest(EvidenceChain(), {})
