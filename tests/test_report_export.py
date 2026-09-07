from __future__ import annotations

from pathlib import Path

import pytest

from src.report.export import ExportError, export_report
from src.reporting import Report


def test_report_exports_to_markdown_and_json(tmp_path: Path) -> None:
    report = Report("Test", "authorized-test")
    markdown = export_report(report, tmp_path / "report", "markdown")
    payload = export_report(report, tmp_path / "report", "json")
    assert markdown.suffix == ".md" and markdown.read_text().startswith("# Test")
    assert payload.suffix == ".json" and '"title": "Test"' in payload.read_text()


def test_report_export_rejects_unknown_format(tmp_path: Path) -> None:
    with pytest.raises(ExportError, match="markdown"):
        export_report(Report("Test", "authorized-test"), tmp_path / "report", "pdf")
