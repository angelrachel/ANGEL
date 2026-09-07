"""Manifest and integrity verification for exported evidence chains."""

from __future__ import annotations

import hashlib
import json
from typing import Any

from .chain import EvidenceChain


class ManifestError(ValueError):
    """Raised when an evidence manifest is malformed or does not verify."""


def build_manifest(chain: EvidenceChain, *, engagement_id: str) -> dict[str, Any]:
    """Build a deterministic manifest without exposing secrets."""

    if not engagement_id.strip():
        raise ManifestError("engagement_id is required")
    exported = chain.export()
    serialized = json.dumps(exported, sort_keys=True, separators=(",", ":")).encode()
    return {
        "engagement_id": engagement_id,
        "record_count": len(exported),
        "first_record_hash": exported[0]["record_hash"] if exported else None,
        "last_record_hash": exported[-1]["record_hash"] if exported else None,
        "chain_sha256": hashlib.sha256(serialized).hexdigest(),
        "chain_valid": chain.verify(),
    }


def verify_manifest(chain: EvidenceChain, manifest: dict[str, Any]) -> bool:
    """Verify record count, chain hash, and chain integrity against a manifest."""

    if not isinstance(manifest, dict) or not isinstance(manifest.get("chain_sha256"), str):
        raise ManifestError("malformed evidence manifest")
    expected = build_manifest(chain, engagement_id=str(manifest.get("engagement_id", "")))
    return expected == manifest and chain.verify()


__all__ = ["ManifestError", "build_manifest", "verify_manifest"]
