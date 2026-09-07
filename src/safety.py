"""Safety guardrails for authorized, non-destructive security assessments.

This module deliberately exposes only passive or simulation capabilities. It is
used as a second policy layer so a caller cannot accidentally turn the control
plane into an exploitation, persistence, credential-theft, evasion, or
 destructive-impact tool.
"""

from __future__ import annotations

from dataclasses import dataclass
from enum import StrEnum


class SafetyError(ValueError):
    """Raised when a requested capability is outside the safe operating mode."""


class Capability(StrEnum):
    PASSIVE_RECON = "passive_recon"
    INVENTORY = "inventory"
    EVIDENCE = "evidence"
    REPORTING = "reporting"
    SIMULATION = "simulation"


@dataclass(frozen=True)
class SafetyDecision:
    capability: Capability
    allowed: bool
    reason: str


_ALLOWED = frozenset(Capability)


def authorize_capability(capability: str | Capability) -> SafetyDecision:
    """Authorize one explicitly safe capability and fail closed otherwise."""

    try:
        normalized = capability if isinstance(capability, Capability) else Capability(capability.strip().lower())
    except (AttributeError, ValueError) as exc:
        raise SafetyError(
            "capability is not allowed; use passive_recon, inventory, evidence, reporting, or simulation"
        ) from exc
    if normalized not in _ALLOWED:  # pragma: no cover - defensive if enum changes
        raise SafetyError("capability is not allowed")
    return SafetyDecision(normalized, True, "explicitly allowed safe capability")


def reject_destructive_request(request: str) -> None:
    """Reject obvious destructive or offensive requests before orchestration."""

    if not isinstance(request, str) or not request.strip():
        raise SafetyError("request must be a non-empty string")
    lowered = request.casefold()
    blocked_terms = (
        "credential theft",
        "keylogger",
        "ransomware",
        "wiper",
        "persistence",
        "process injection",
        "bypass edr",
        "disable antivirus",
        "golden ticket",
        "shellcode",
        "exfiltrate",
        "delete logs",
    )
    if any(term in lowered for term in blocked_terms):
        raise SafetyError("offensive, evasive, credential-theft, or destructive request rejected")


__all__ = ["Capability", "SafetyDecision", "SafetyError", "authorize_capability", "reject_destructive_request"]
