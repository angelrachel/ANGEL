"""Report adapters for GraphQL and authorization-matrix results."""

from __future__ import annotations

from .auth_matrix import MatrixFinding
from .graphql_intel import GraphQLDrift
from .reporting import Finding


def matrix_to_finding(result: MatrixFinding) -> Finding:
    return Finding(
        title=f"Authorization matrix mismatch: {result.operation}",
        severity="high" if result.confidence == "high" else "medium",
        confidence=result.confidence,
        asset=result.resource,
        summary=f"{result.actor} ({result.role}) observed {result.observed} where policy expected {result.expected}.",
        impact=f"Authorization behavior differs for tenant {result.resource_tenant}.",
        reproduction=[
            f"Use actor {result.actor}.",
            f"Attempt {result.operation} on canary {result.resource}.",
            f"Observe HTTP {result.status}.",
        ],
        remediation="Enforce authorization server-side for every operation and resource tenant.",
        evidence_refs=[f"matrix:{result.actor}:{result.resource}:{result.operation}"],
    )


def graphql_drift_to_finding(result: GraphQLDrift) -> Finding:
    severity = "high" if result.category == "authorization_drift" else "medium"
    return Finding(
        title=f"GraphQL schema {result.category}",
        severity=severity,
        confidence="medium",
        asset=result.subject,
        summary=result.details,
        impact="The observed GraphQL contract differs from the expected authorization or schema model.",
        reproduction=[
            "Compare the expected and observed introspection snapshots.",
            f"Review drift at {result.subject}.",
        ],
        remediation="Keep schema documentation and resolver authorization requirements synchronized.",
        evidence_refs=[f"graphql:{result.category}:{result.subject}"],
    )


__all__ = ["graphql_drift_to_finding", "matrix_to_finding"]
