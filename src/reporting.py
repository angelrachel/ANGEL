"""Structured finding and report exports."""

from __future__ import annotations

import json
from dataclasses import asdict, dataclass, field
from typing import Any


@dataclass
class Finding:
    title: str
    severity: str
    confidence: str
    asset: str
    summary: str
    impact: str
    reproduction: list[str]
    remediation: str
    evidence_refs: list[str] = field(default_factory=list)


@dataclass
class Report:
    title: str
    engagement: str
    findings: list[Finding] = field(default_factory=list)
    executive_summary: str = ""

    def add(self, finding: Finding) -> None:
        self.findings.append(finding)

    def to_dict(self) -> dict[str, Any]:
        return asdict(self)

    def to_json(self) -> str:
        return json.dumps(self.to_dict(), indent=2, sort_keys=True)

    def to_markdown(self) -> str:
        lines = [f"# {self.title}", "", f"**Engagement:** {self.engagement}", ""]
        if self.executive_summary:
            lines += ["## Executive Summary", "", self.executive_summary, ""]
        lines += ["## Findings", ""]
        if not self.findings:
            lines.append("No findings recorded.")
        for index, finding in enumerate(self.findings, start=1):
            lines += [
                f"### {index}. {finding.title}",
                "",
                f"- **Severity:** {finding.severity}",
                f"- **Confidence:** {finding.confidence}",
                f"- **Asset:** `{finding.asset}`",
                "",
                f"**Summary:** {finding.summary}",
                "",
                f"**Impact:** {finding.impact}",
                "",
                "**Reproduction:**",
                *[f"{step}. {instruction}" for step, instruction in enumerate(finding.reproduction, start=1)],
                "",
                f"**Remediation:** {finding.remediation}",
                "",
                f"**Evidence:** {', '.join(finding.evidence_refs) if finding.evidence_refs else 'None'}",
                "",
            ]
        return "\n".join(lines).rstrip() + "\n"


__all__ = ["Finding", "Report"]
