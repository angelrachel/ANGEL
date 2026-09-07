"""Convert passive metadata signals into conservative report findings."""

from __future__ import annotations

from collections.abc import Iterable

from .osint.passive_checks import PassiveSignal
from .reporting import Finding


def signals_to_findings(signals: Iterable[PassiveSignal], asset: str) -> list[Finding]:
    """Create informational findings from already-collected passive signals.

    This function never claims exploitability. It records an observation and
    asks the owner to validate whether the exposure is intended.
    """

    if not asset.strip():
        raise ValueError("asset is required")
    findings: list[Finding] = []
    seen: set[tuple[str, str]] = set()
    for signal in signals:
        key = (signal.kind, signal.value)
        if key in seen:
            continue
        seen.add(key)
        findings.append(
            Finding(
                title=f"Passive {signal.kind} signal observed",
                severity="info",
                confidence="medium",
                asset=asset,
                summary=f"Observed {signal.kind} metadata: {signal.value}.",
                impact=(
                    "The observation may assist asset inventory or defensive configuration review; "
                    "it is not proof of a vulnerability."
                ),
                reproduction=[
                    f"Review the captured metadata from {signal.source}.",
                    f"Confirm the {signal.kind} value {signal.value} with the system owner.",
                ],
                remediation=(
                    "Confirm the asset and metadata are expected, minimize unnecessary disclosure, "
                    "and document the decision."
                ),
                evidence_refs=[f"passive:{signal.source}:{signal.kind}:{signal.value}"],
            )
        )
    return findings


__all__ = ["signals_to_findings"]
