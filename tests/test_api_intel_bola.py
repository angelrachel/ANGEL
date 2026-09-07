from __future__ import annotations

import pytest

from src.api_intel import compare_openapi, inventory_openapi
from src.bola import Actor, Canary, ResponseSnapshot, validate_cross_tenant_read, validate_cross_tenant_write
from src.finding_adapter import authorization_to_finding


def test_openapi_inventory_and_drift() -> None:
    expected = {
        "security": [{"oauth": []}],
        "paths": {"/users": {"get": {"operationId": "users", "parameters": [{"name": "page"}]}}},
    }
    observed = {"paths": {"/users": {"get": {"operationId": "users"}}, "/admin": {"get": {}}}}
    endpoints = inventory_openapi(expected)
    assert endpoints[0].operation_id == "users"
    issues = compare_openapi(expected, observed)
    assert {issue.issue for issue in issues} == {"undocumented_endpoint", "parameter_drift", "security_drift"}


def test_bola_read_finding_and_adapter() -> None:
    actor = Actor("user-a", "tenant-a", "member")
    own = Canary("record-a", "tenant-a", "marker-a")
    foreign = Canary("record-b", "tenant-b", "marker-b")

    def request(current: Actor, resource: Canary) -> ResponseSnapshot:
        return ResponseSnapshot(200, {"marker": resource.marker})

    result = validate_cross_tenant_read(actor, own, foreign, request)
    assert result is not None and result.marker_present
    finding = authorization_to_finding(result)
    assert finding.severity == "high"
    assert "tenant-b" in finding.summary


def test_bola_stops_on_bad_control_and_detects_write() -> None:
    actor = Actor("user-a", "tenant-a", "member")
    own = Canary("record-a", "tenant-a", "marker-a")
    foreign = Canary("record-b", "tenant-b", "marker-b")
    with pytest.raises(ValueError):
        validate_cross_tenant_read(actor, own, foreign, lambda _a, _c: ResponseSnapshot(403, {}))
    write = validate_cross_tenant_write(actor, foreign, lambda _a, _c: ResponseSnapshot(200, {}, changed_state=True))
    assert write is not None and write.operation == "write"
    assert validate_cross_tenant_write(actor, foreign, lambda _a, _c: ResponseSnapshot(403, {}, False)) is None
