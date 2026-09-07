from __future__ import annotations

import pytest

from src.orchestrator.mcp import MockMCPServer, ToolCall
from src.orchestrator.policy import PermissionRegistry, PolicyError, RiskLevel, ToolPermission


def test_mock_mcp_server_enforces_permissions() -> None:
    permissions = PermissionRegistry()
    permissions.register(ToolPermission("inventory", RiskLevel.LOW, False))
    server = MockMCPServer(permissions)
    server.register_tool("inventory", lambda arguments: {"count": arguments.get("count", 0)})
    result = server.call(ToolCall("req-1", "inventory", {"count": 2}))
    assert result["result"] == {"count": 2}
    permissions.register(ToolPermission("missing", RiskLevel.LOW, False))
    with pytest.raises(PolicyError, match="handler"):
        server.call(ToolCall("req-2", "missing", {}))


def test_mock_mcp_server_requires_approval() -> None:
    permissions = PermissionRegistry()
    permissions.register(ToolPermission("review", RiskLevel.MEDIUM, True))
    server = MockMCPServer(permissions)
    server.register_tool("review", lambda _: {"status": "ok"})
    with pytest.raises(PolicyError, match="approval"):
        server.call(ToolCall("req-1", "review", {}))
    assert server.call(ToolCall("req-1", "review", {}), approved=True)["result"]["status"] == "ok"
