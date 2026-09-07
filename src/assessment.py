"""Auditable plans for passive and simulated security assessments."""

from __future__ import annotations

from dataclasses import dataclass
from enum import StrEnum
from typing import ClassVar

from .engagement import RulesOfEngagement


class AssessmentError(ValueError):
    """Raised when an assessment plan is unsafe or incomplete."""


class AssessmentMode(StrEnum):
    PASSIVE = "passive"
    SIMULATION = "simulation"


@dataclass(frozen=True)
class AssessmentPlan:
    engagement: RulesOfEngagement
    mode: AssessmentMode
    checks: tuple[str, ...]
    requested_by: str

    SAFE_CHECKS: ClassVar[frozenset[str]] = frozenset(
        {
            "asset_inventory",
            "certificate_metadata",
            "response_headers",
            "robots_metadata",
            "api_schema_review",
            "graphql_schema_review",
            "evidence_integrity",
            "control_plane_self_test",
        }
    )

    def __post_init__(self) -> None:
        if not self.requested_by.strip():
            raise AssessmentError("requested_by is required")
        normalized = tuple(dict.fromkeys(check.strip().lower() for check in self.checks if check.strip()))
        if not normalized:
            raise AssessmentError("at least one assessment check is required")
        unknown = set(normalized) - self.SAFE_CHECKS
        if unknown:
            raise AssessmentError(f"assessment checks are not allowed: {', '.join(sorted(unknown))}")
        if not self.engagement.is_active():
            raise AssessmentError("engagement authorization is not active")
        object.__setattr__(self, "checks", normalized)

    def to_dict(self) -> dict[str, object]:
        return {
            "engagement_id": self.engagement.engagement_id,
            "mode": self.mode.value,
            "checks": list(self.checks),
            "requested_by": self.requested_by,
            "hosts": sorted(self.engagement.hosts),
        }


__all__ = ["AssessmentError", "AssessmentMode", "AssessmentPlan"]
