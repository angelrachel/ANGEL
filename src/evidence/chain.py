"""Tamper-evident evidence records with secret redaction."""

from __future__ import annotations

import hashlib
import json
import re
import time
from dataclasses import asdict, dataclass
from typing import Any

_SECRET_PATTERNS = (
    re.compile(r"(?i)(authorization\s*:\s*bearer\s+)[^\s]+"),
    re.compile(r"(?i)(cookie\s*:\s*)[^\r\n]+"),
    re.compile(r"(?i)(api[_-]?key|token|password|secret)(\s*[:=]\s*)[^,\s}]+"),
)


def redact(value: str) -> str:
    result = value
    for pattern in _SECRET_PATTERNS:
        result = pattern.sub(lambda match: f"{match.group(1)}[REDACTED]" if match.lastindex else "[REDACTED]", result)
    return result


@dataclass(frozen=True)
class EvidenceRecord:
    sequence: int
    evidence_type: str
    actor: str
    content: dict[str, Any]
    created_at: int
    previous_hash: str
    record_hash: str


class EvidenceChain:
    def __init__(self) -> None:
        self.records: list[EvidenceRecord] = []

    def append(
        self,
        evidence_type: str,
        actor: str,
        content: dict[str, Any],
        *,
        created_at: int | None = None,
    ) -> EvidenceRecord:
        safe_content = json.loads(redact(json.dumps(content, sort_keys=True, separators=(",", ":"))))
        timestamp = int(time.time()) if created_at is None else created_at
        previous = self.records[-1].record_hash if self.records else "0" * 64
        sequence = len(self.records) + 1
        material = {
            "sequence": sequence,
            "type": evidence_type,
            "actor": actor,
            "content": safe_content,
            "created_at": timestamp,
            "previous_hash": previous,
        }
        digest = hashlib.sha256(json.dumps(material, sort_keys=True, separators=(",", ":")).encode()).hexdigest()
        record = EvidenceRecord(sequence, evidence_type, actor, safe_content, timestamp, previous, digest)
        self.records.append(record)
        return record

    def verify(self) -> bool:
        previous = "0" * 64
        for expected_sequence, record in enumerate(self.records, start=1):
            material = {
                "sequence": record.sequence,
                "type": record.evidence_type,
                "actor": record.actor,
                "content": record.content,
                "created_at": record.created_at,
                "previous_hash": record.previous_hash,
            }
            digest = hashlib.sha256(json.dumps(material, sort_keys=True, separators=(",", ":")).encode()).hexdigest()
            if record.sequence != expected_sequence or record.previous_hash != previous or record.record_hash != digest:
                return False
            previous = record.record_hash
        return True

    def export(self) -> list[dict[str, Any]]:
        return [asdict(record) for record in self.records]


__all__ = ["EvidenceChain", "EvidenceRecord", "redact"]
