from __future__ import annotations

from src.auth_matrix import MatrixActor, MatrixObservation, MatrixResource, build_matrix, evaluate_matrix
from src.graphql_intel import compare_schema, inventory_operation, inventory_schema
from src.graphql_matrix_adapter import graphql_drift_to_finding, matrix_to_finding


def schema(auth: bool = True) -> dict:
    return {
        "data": {
            "__schema": {
                "types": [
                    {
                        "kind": "OBJECT",
                        "name": "Query",
                        "fields": [
                            {
                                "name": "account",
                                "requiresAuth": auth,
                                "args": [{"name": "id"}],
                                "type": {"name": "Account"},
                            }
                        ],
                    }
                ]
            }
        }
    }


def test_graphql_schema_and_operation_inventory() -> None:
    fields = inventory_schema(schema())
    assert fields[0].name == "account"
    operations = inventory_operation('query GetAccount { account(id: "x") { id } }')
    assert operations[0].name == "GetAccount"
    assert "account" in operations[0].fields


def test_graphql_schema_drift_adapter() -> None:
    issues = compare_schema(schema(), schema(False))
    assert issues[0].category == "authorization_drift"
    assert graphql_drift_to_finding(issues[0]).severity == "high"


def test_authorization_matrix_evaluation() -> None:
    actors = [MatrixActor("a", "tenant-a", "member"), MatrixActor("b", "tenant-b", "member")]
    resources = [MatrixResource("r-a", "tenant-a", "marker-a"), MatrixResource("r-b", "tenant-b", "marker-b")]
    cases = build_matrix(actors, resources, ("read",))

    def observe(case):
        if case.actor.tenant == case.resource.tenant:
            return MatrixObservation(200, marker_present=True)
        return MatrixObservation(200, marker_present=True)

    findings = evaluate_matrix(cases, observe)
    assert len(findings) == 2
    assert matrix_to_finding(findings[0]).confidence == "high"
