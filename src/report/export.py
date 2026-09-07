"""Validated report export helpers."""

from __future__ import annotations

import json
from pathlib import Path

from ..reporting import Report


class ExportError(ValueError):
    """Raised when an unsupported report export is requested."""


def export_report(report: Report, destination: str | Path, output_format: str = "markdown") -> Path:
    selected = output_format.lower()
    if selected == "markdown":
        content = report.to_markdown()
        suffix = ".md"
    elif selected == "json":
        content = report.to_json() + "\n"
        suffix = ".json"
    else:
        raise ExportError("format must be markdown or json")
    path = Path(destination)
    if path.suffix.lower() != suffix:
        path = path.with_suffix(suffix)
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(content, encoding="utf-8")
    if selected == "json":
        json.loads(path.read_text(encoding="utf-8"))
    return path


__all__ = ["ExportError", "export_report"]
