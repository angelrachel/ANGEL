"""Adapters from safe validators to report findings."""

from __future__ import annotations

from .api_intel import SchemaDrift
from .bola import AuthorizationFinding
from .graphql_intel import GraphQLDrift
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


def schema_drift_to_finding(result: SchemaDrift | GraphQLDrift) -> Finding:
    issue = result.issue if isinstance(result, SchemaDrift) else result.category
    subject = result.path if isinstance(result, SchemaDrift) else result.subject
    severity = "high" if "security" in issue or "authorization" in issue else "medium"
    return Finding(
        title=f"Schema drift: {issue}",
        severity=severity,
        confidence="high",
        asset=subject,
        summary=result.details,
        impact=(
            "The observed interface differs from the expected contract and may weaken "
            "authorization or availability assumptions."
        ),
        reproduction=[
            "Capture the expected schema contract.",
            "Capture the observed schema contract.",
            f"Compare {subject} and record the drift.",
        ],
        remediation="Review the contract change, update authorization tests, and deploy only an approved schema.",
        evidence_refs=[f"schema-drift:{subject}"],
    )


__all__ = ["authorization_to_finding", "schema_drift_to_finding"]
