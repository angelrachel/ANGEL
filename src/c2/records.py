"""Serializable records persisted by the ANGEL control plane."""

from __future__ import annotations

from dataclasses import dataclass
from typing import Any


@dataclass(frozen=True)
class ScopeRecord:
    id: int
    name: str
    hosts: list[str]
    paths: list[str]
    expires_at: int | None
    status: str


@dataclass(frozen=True)
class EvidenceRecordRow:
    id: int
    scope_id: int
    evidence_type: str
    actor: str
    payload: dict[str, Any]
    record_hash: str
    created_at: int


@dataclass(frozen=True)
class ReportRecord:
    id: int
    scope_id: int
    title: str
    payload: dict[str, Any]
    created_at: int


__all__ = ["ScopeRecord", "EvidenceRecordRow", "ReportRecord"]
