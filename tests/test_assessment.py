from __future__ import annotations

from datetime import UTC, datetime, timedelta

import pytest

from src.assessment import AssessmentError, AssessmentMode, AssessmentPlan
from src.engagement import RulesOfEngagement


def engagement(active: bool = True) -> RulesOfEngagement:
    start = datetime.now(UTC) - timedelta(minutes=1) if active else datetime.now(UTC) - timedelta(days=2)
    return RulesOfEngagement(
        "eng-assessment",
        "owner@example.test",
        start,
        start + (timedelta(hours=1) if active else timedelta(minutes=1)),
        frozenset({"example.test"}),
        "on-call@example.test",
    )


def test_plan_normalizes_checks_and_serializes() -> None:
    plan = AssessmentPlan(
        engagement(),
        AssessmentMode.PASSIVE,
        ("asset_inventory", "response_headers", "asset_inventory"),
        "operator@example.test",
    )
    assert plan.checks == ("asset_inventory", "response_headers")
    assert plan.to_dict()["mode"] == "passive"


def test_plan_rejects_unknown_or_expired_checks() -> None:
    with pytest.raises(AssessmentError, match="not allowed"):
        AssessmentPlan(engagement(), AssessmentMode.SIMULATION, ("process_injection",), "operator")
    with pytest.raises(AssessmentError, match="not active"):
        AssessmentPlan(engagement(active=False), AssessmentMode.PASSIVE, ("asset_inventory",), "operator")
    with pytest.raises(AssessmentError, match="required"):
        AssessmentPlan(engagement(), AssessmentMode.PASSIVE, ("asset_inventory",), "")
