"""Fail-closed policy and approval primitives for orchestrated work."""

from __future__ import annotations

from dataclasses import dataclass
from enum import StrEnum


class RiskLevel(StrEnum):
    LOW = "low"
    MEDIUM = "medium"
    HIGH = "high"


class PolicyError(ValueError):
    """Raised when a policy decision cannot be made safely."""


@dataclass(frozen=True)
class ToolPermission:
    name: str
    risk: RiskLevel
    requires_approval: bool


@dataclass(frozen=True)
class Approval:
    request_id: str
    actor: str
    tool: str
    approved: bool
    reason: str


class PermissionRegistry:
    def __init__(self) -> None:
        self._permissions: dict[str, ToolPermission] = {}

    def register(self, permission: ToolPermission) -> None:
        if not permission.name.strip():
            raise PolicyError("tool name is required")
        if permission.name in self._permissions:
            raise PolicyError("tool permission already exists")
        self._permissions[permission.name] = permission

    def get(self, tool: str) -> ToolPermission:
        try:
            return self._permissions[tool]
        except KeyError as exc:
            raise PolicyError("tool is not registered") from exc

    def authorize(self, tool: str, *, approved: bool = False) -> None:
        permission = self.get(tool)
        if permission.requires_approval and not approved:
            raise PolicyError("tool approval is required")


class ApprovalGate:
    def __init__(self) -> None:
        self._approvals: dict[str, Approval] = {}

    def decide(self, request_id: str, actor: str, tool: str, approved: bool, reason: str) -> Approval:
        if not request_id.strip() or not actor.strip() or not tool.strip() or not reason.strip():
            raise PolicyError("approval fields are required")
        if request_id in self._approvals:
            raise PolicyError("approval request already decided")
        approval = Approval(request_id, actor, tool, approved, reason)
        self._approvals[request_id] = approval
        return approval

    def get(self, request_id: str) -> Approval:
        try:
            return self._approvals[request_id]
        except KeyError as exc:
            raise PolicyError("approval request not found") from exc


__all__ = ["Approval", "ApprovalGate", "PermissionRegistry", "PolicyError", "RiskLevel", "ToolPermission"]
