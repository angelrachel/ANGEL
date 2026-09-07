"""Validated rules of engagement for authorized assessments.

The control plane may consume this object before scheduling any check. It is
intentionally data-only: it does not grant access or execute network actions.
"""

from __future__ import annotations

from dataclasses import dataclass
from datetime import UTC, datetime


class EngagementError(ValueError):
    """Raised when an engagement authorization is incomplete or inconsistent."""


@dataclass(frozen=True)
class RulesOfEngagement:
    engagement_id: str
    authorized_by: str
    starts_at: datetime
    ends_at: datetime
    hosts: frozenset[str]
    emergency_contact: str
    prohibited_actions: frozenset[str] = frozenset(
        {
            "credential_theft",
            "persistence",
            "destructive_impact",
            "evasion",
            "data_exfiltration",
        }
    )

    def __post_init__(self) -> None:
        if not self.engagement_id.strip() or not self.authorized_by.strip():
            raise EngagementError("engagement_id and authorized_by are required")
        if not self.emergency_contact.strip():
            raise EngagementError("emergency_contact is required")
        if self.starts_at.tzinfo is None or self.ends_at.tzinfo is None:
            raise EngagementError("engagement timestamps must include a timezone")
        if self.ends_at <= self.starts_at:
            raise EngagementError("engagement must end after it starts")
        normalized_hosts = frozenset(host.lower().rstrip(".") for host in self.hosts if host.strip())
        if not normalized_hosts:
            raise EngagementError("at least one authorized host is required")
        if any("/" in host or " " in host for host in normalized_hosts):
            raise EngagementError("hosts must be DNS names or IP addresses, not URLs")
        object.__setattr__(self, "hosts", normalized_hosts)

    def is_active(self, now: datetime | None = None) -> bool:
        """Return whether the authorization window is currently active."""

        current = now or datetime.now(UTC)
        if current.tzinfo is None:
            raise EngagementError("current time must include a timezone")
        return self.starts_at <= current < self.ends_at

    def allows_host(self, host: str) -> bool:
        """Match an exact host or a subdomain of an authorized DNS suffix."""

        candidate = host.strip().lower().rstrip(".")
        return bool(candidate) and any(
            candidate == allowed or candidate.endswith("." + allowed) for allowed in self.hosts
        )

    def to_dict(self) -> dict[str, object]:
        return {
            "engagement_id": self.engagement_id,
            "authorized_by": self.authorized_by,
            "starts_at": self.starts_at.isoformat(),
            "ends_at": self.ends_at.isoformat(),
            "hosts": sorted(self.hosts),
            "emergency_contact": self.emergency_contact,
            "prohibited_actions": sorted(self.prohibited_actions),
        }


__all__ = ["EngagementError", "RulesOfEngagement"]
