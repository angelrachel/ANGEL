from __future__ import annotations

from src.api_intel import SchemaDrift
from src.finding_adapter import schema_drift_to_finding
from src.graphql_intel import GraphQLDrift


def test_schema_drift_finding_adapter() -> None:
    api_finding = schema_drift_to_finding(SchemaDrift("/users", "GET", "security_drift", "Auth changed"))
    gql_finding = schema_drift_to_finding(GraphQLDrift("authorization_drift", "Query.users", "Auth changed"))
    assert api_finding.severity == "high"
    assert gql_finding.asset == "Query.users"
    assert api_finding.evidence_refs == ["schema-drift:/users"]
