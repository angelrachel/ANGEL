from __future__ import annotations

import pytest

from src.orchestrator.policy import ApprovalGate, PermissionRegistry, PolicyError, RiskLevel, ToolPermission


def test_permission_registry_is_fail_closed() -> None:
    registry = PermissionRegistry()
    registry.register(ToolPermission("passive_inventory", RiskLevel.LOW, False))
    registry.authorize("passive_inventory")
    with pytest.raises(PolicyError, match="registered"):
        registry.authorize("unknown")
    registry.register(ToolPermission("active_check", RiskLevel.MEDIUM, True))
    with pytest.raises(PolicyError, match="approval"):
        registry.authorize("active_check")
    registry.authorize("active_check", approved=True)


def test_approval_gate_is_immutable_per_request() -> None:
    gate = ApprovalGate()
    approval = gate.decide("req-1", "operator", "active_check", True, "authorized lab")
    assert gate.get("req-1") == approval
    with pytest.raises(PolicyError, match="already"):
        gate.decide("req-1", "operator", "active_check", False, "changed")
