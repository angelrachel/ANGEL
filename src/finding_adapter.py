"""Adapters from safe validators to report findings."""

from __future__ import annotations

from .bola import AuthorizationFinding
from .reporting import Finding


def authorization_to_finding(result: AuthorizationFinding) -> Finding:
    severity = "high" if result.operation == "read" and result.marker_present else "medium"
    return Finding(
        title=f"Cross-tenant {result.operation} authorization failure",
        severity=severity,
        confidence=result.confidence,
        asset=result.resource,
        summary=f"Actor {result.actor} accessed a resource owned by tenant {result.resource_tenant}.",
        impact="Tenant isolation is not enforced for the tested canary resource.",
        reproduction=[
            f"Use actor {result.actor} from tenant {result.actor_tenant}.",
            f"Request operation {result.operation} against canary {result.resource}.",
            f"Observe HTTP {result.observed_status} and the canary authorization result.",
        ],
        remediation="Enforce server-side authorization against resource ownership and tenant context.",
        evidence_refs=[f"authorization:{result.actor}:{result.resource}"],
    )


__all__ = ["authorization_to_finding"]
