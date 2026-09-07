"""Policy-driven authorization matrix comparison for authorized test fixtures."""

from __future__ import annotations

from collections.abc import Callable
from dataclasses import dataclass


@dataclass(frozen=True)
class MatrixActor:
    name: str
    tenant: str
    role: str


@dataclass(frozen=True)
class MatrixResource:
    name: str
    tenant: str
    marker: str


@dataclass(frozen=True)
class MatrixCase:
    actor: MatrixActor
    resource: MatrixResource
    operation: str
    expected: str


@dataclass(frozen=True)
class MatrixObservation:
    status: int
    marker_present: bool = False
    state_changed: bool = False


@dataclass(frozen=True)
class MatrixFinding:
    actor: str
    role: str
    actor_tenant: str
    resource: str
    resource_tenant: str
    operation: str
    expected: str
    observed: str
    status: int
    confidence: str


Observe = Callable[[MatrixCase], MatrixObservation]


def validate_matrix(cases: list[MatrixCase]) -> None:
    if not cases:
        raise ValueError("authorization matrix cannot be empty")
    seen: set[tuple[str, str, str]] = set()
    for case in cases:
        if case.expected not in {"allow", "deny"}:
            raise ValueError("matrix expected value must be allow or deny")
        if case.operation not in {"read", "write", "delete"}:
            raise ValueError("unsupported matrix operation")
        if not case.actor.name.strip() or not case.resource.name.strip():
            raise ValueError("matrix actor and resource names are required")
        key = (case.actor.name, case.resource.name, case.operation)
        if key in seen:
            raise ValueError("duplicate authorization matrix case")
        seen.add(key)


def build_matrix(
    actors: list[MatrixActor],
    resources: list[MatrixResource],
    operations: tuple[str, ...] = ("read", "write"),
) -> list[MatrixCase]:
    cases: list[MatrixCase] = []
    for actor in actors:
        for resource in resources:
            for operation in operations:
                expected = "allow" if actor.tenant == resource.tenant and operation == "read" else "deny"
                cases.append(MatrixCase(actor, resource, operation, expected))
    return cases


def evaluate_matrix(cases: list[MatrixCase], observe: Observe) -> list[MatrixFinding]:
    validate_matrix(cases)
    findings: list[MatrixFinding] = []
    for case in cases:
        observation = observe(case)
        allowed = 200 <= observation.status < 300
        denied_case_violated = case.expected == "deny" and (
            allowed or observation.marker_present or observation.state_changed
        )
        allowed_case_violated = case.expected == "allow" and not allowed
        violated = denied_case_violated or allowed_case_violated
        if not violated:
            continue
        findings.append(
            MatrixFinding(
                case.actor.name,
                case.actor.role,
                case.actor.tenant,
                case.resource.name,
                case.resource.tenant,
                case.operation,
                case.expected,
                "allow" if allowed else "deny",
                observation.status,
                "high"
                if case.expected == "deny" and (observation.marker_present or observation.state_changed)
                else "medium",
            )
        )
    return findings


__all__ = [
    "MatrixActor",
    "MatrixCase",
    "MatrixFinding",
    "MatrixObservation",
    "MatrixResource",
    "build_matrix",
    "evaluate_matrix",
    "validate_matrix",
]
