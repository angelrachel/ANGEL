"""Non-destructive, canary-based authorization comparison for BOLA/IDOR testing."""

from __future__ import annotations

from collections.abc import Callable
from dataclasses import dataclass
from typing import Any


@dataclass(frozen=True)
class Actor:
    name: str
    tenant: str
    role: str


@dataclass(frozen=True)
class Canary:
    identifier: str
    tenant: str
    marker: str


@dataclass(frozen=True)
class ResponseSnapshot:
    status: int
    body: Any
    changed_state: bool = False


@dataclass(frozen=True)
class AuthorizationFinding:
    actor: str
    actor_tenant: str
    resource: str
    resource_tenant: str
    operation: str
    expected: str
    observed_status: int
    marker_present: bool
    confidence: str
    evidence: dict[str, Any]


RequestFn = Callable[[Actor, Canary], ResponseSnapshot]


def validate_cross_tenant_read(
    actor: Actor, own: Canary, foreign: Canary, request: RequestFn
) -> AuthorizationFinding | None:
    own_response = request(actor, own)
    if own_response.status < 200 or own_response.status >= 300:
        raise ValueError("control canary must be readable before comparison")
    foreign_response = request(actor, foreign)
    marker_present = foreign.marker in str(foreign_response.body)
    if foreign_response.changed_state:
        raise ValueError("read validator received a state-changing response")
    if marker_present:
        return AuthorizationFinding(
            actor.name,
            actor.tenant,
            foreign.identifier,
            foreign.tenant,
            "read",
            "deny",
            foreign_response.status,
            True,
            "high",
            {
                "control_status": own_response.status,
                "cross_tenant_status": foreign_response.status,
                "marker": foreign.marker,
            },
        )
    return None


def validate_cross_tenant_write(actor: Actor, foreign: Canary, request: RequestFn) -> AuthorizationFinding | None:
    response = request(actor, foreign)
    if response.changed_state:
        return AuthorizationFinding(
            actor.name,
            actor.tenant,
            foreign.identifier,
            foreign.tenant,
            "write",
            "deny",
            response.status,
            False,
            "high",
            {"cross_tenant_status": response.status, "state_changed": True},
        )
    return None


__all__ = [
    "Actor",
    "AuthorizationFinding",
    "Canary",
    "ResponseSnapshot",
    "validate_cross_tenant_read",
    "validate_cross_tenant_write",
]
