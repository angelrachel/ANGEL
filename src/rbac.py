"""In-memory RBAC registry used by API composition and lab deployments."""

from __future__ import annotations

from dataclasses import dataclass


class RBACError(ValueError):
    """Raised when a role or authorization decision is invalid."""


@dataclass(frozen=True)
class Principal:
    subject: str
    roles: frozenset[str]
    revoked: bool = False


class RBAC:
    def __init__(self) -> None:
        self._roles: dict[str, frozenset[str]] = {}
        self._principals: dict[str, Principal] = {}

    def define_role(self, role: str, permissions: set[str]) -> None:
        if not role.strip() or not permissions or any(not item.strip() for item in permissions):
            raise RBACError("role and permissions are required")
        if role in self._roles:
            raise RBACError("role already exists")
        self._roles[role] = frozenset(permissions)

    def upsert_principal(self, subject: str, roles: set[str]) -> Principal:
        if not subject.strip() or not roles or not roles.issubset(self._roles):
            raise RBACError("principal has invalid roles")
        principal = Principal(subject, frozenset(roles))
        self._principals[subject] = principal
        return principal

    def revoke(self, subject: str) -> None:
        try:
            principal = self._principals[subject]
        except KeyError as exc:
            raise RBACError("principal not found") from exc
        self._principals[subject] = Principal(principal.subject, principal.roles, True)

    def authorize(self, subject: str, permission: str) -> None:
        try:
            principal = self._principals[subject]
        except KeyError as exc:
            raise RBACError("principal not found") from exc
        if principal.revoked:
            raise RBACError("principal is revoked")
        granted = {item for role in principal.roles for item in self._roles[role]}
        if permission not in granted:
            raise RBACError("permission denied")


__all__ = ["Principal", "RBAC", "RBACError"]
