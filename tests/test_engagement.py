from __future__ import annotations

from datetime import UTC, datetime, timedelta

import pytest

from src.engagement import EngagementError, RulesOfEngagement


def make_roe() -> RulesOfEngagement:
    start = datetime.now(UTC) - timedelta(minutes=1)
    return RulesOfEngagement(
        engagement_id="eng-2026-001",
        authorized_by="security-owner@example.test",
        starts_at=start,
        ends_at=start + timedelta(hours=2),
        hosts=frozenset({"example.test"}),
        emergency_contact="on-call@example.test",
    )


def test_rules_of_engagement_is_active_and_scope_aware() -> None:
    roe = make_roe()
    assert roe.is_active()
    assert roe.allows_host("api.example.test")
    assert not roe.allows_host("example.test.evil.test")
    assert "credential_theft" in roe.prohibited_actions


def test_rules_of_engagement_rejects_naive_or_empty_authorization() -> None:
    start = datetime.now(UTC)
    with pytest.raises(EngagementError, match="timezone"):
        RulesOfEngagement(
            "eng",
            "owner",
            start.replace(tzinfo=None),
            start + timedelta(hours=1),
            frozenset({"example.test"}),
            "on-call",
        )
    with pytest.raises(EngagementError, match="required"):
        RulesOfEngagement("", "owner", start, start + timedelta(hours=1), frozenset({"example.test"}), "on-call")


def test_rules_of_engagement_rejects_url_as_host() -> None:
    start = datetime.now(UTC)
    with pytest.raises(EngagementError, match="DNS"):
        RulesOfEngagement(
            "eng", "owner", start, start + timedelta(hours=1), frozenset({"https://example.test"}), "on-call"
        )
