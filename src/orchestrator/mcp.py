"""Minimal permission-aware MCP-style adapter for local integrations."""

from __future__ import annotations

from collections.abc import Callable
from dataclasses import dataclass
from typing import Any

from .policy import PermissionRegistry, PolicyError


@dataclass(frozen=True)
class ToolCall:
    request_id: str
    tool: str
    arguments: dict[str, Any]


class MockMCPServer:
    def __init__(self, permissions: PermissionRegistry) -> None:
        self.permissions = permissions
        self._tools: dict[str, Callable[[dict[str, Any]], dict[str, Any]]] = {}

    def register_tool(self, name: str, handler: Callable[[dict[str, Any]], dict[str, Any]]) -> None:
        self.permissions.get(name)
        if name in self._tools:
            raise PolicyError("tool handler already exists")
        self._tools[name] = handler

    def call(self, request: ToolCall, *, approved: bool = False) -> dict[str, Any]:
        self.permissions.authorize(request.tool, approved=approved)
        try:
            handler = self._tools[request.tool]
        except KeyError as exc:
            raise PolicyError("tool handler not found") from exc
        result = handler(request.arguments)
        if not isinstance(result, dict):
            raise PolicyError("tool result must be an object")
        return {"request_id": request.request_id, "tool": request.tool, "result": result}


__all__ = ["MockMCPServer", "ToolCall"]
