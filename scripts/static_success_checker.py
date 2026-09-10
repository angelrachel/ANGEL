#!/usr/bin/env python3
"""Report suspicious success claims and ignored command errors without executing code."""
from __future__ import annotations

import argparse
import re
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
DEFAULT_ROOTS = [ROOT / "src" / "assessment"]
IGNORED_FUNCTIONS = re.compile(r"(always|dummy|mock|stub|success|placeholder|notimplemented)", re.IGNORECASE)
FUNCTION = re.compile(r"^func(?:\s+\([^)]*\))?\s+([A-Za-z0-9_]+)\s*\(")
IGNORED_COMMAND = re.compile(r"\.Run\(\)")
SUCCESS_STATUS = re.compile(r'Status\s*:\s*"success"')
RETURN_TRUE = re.compile(r"return\s+true")


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--strict", action="store_true", help="return non-zero when findings exist")
    args = parser.parse_args()
    findings: list[tuple[str, int, str, str]] = []
    for root in DEFAULT_ROOTS:
        for path in sorted(root.rglob("*.go")):
            if path.name.endswith("_test.go"):
                continue
            function = ""
            for number, line in enumerate(path.read_text(errors="replace").splitlines(), 1):
                match = FUNCTION.match(line.strip())
                if match:
                    function = match.group(1)
                if IGNORED_COMMAND.search(line):
                    findings.append((str(path.relative_to(ROOT)), number, "ignored-command-error", line.strip()))
                if SUCCESS_STATUS.search(line) and "DeliveryResult" not in line:
                    findings.append((str(path.relative_to(ROOT)), number, "success-literal", line.strip()))
                if RETURN_TRUE.search(line) and IGNORED_FUNCTIONS.search(function):
                    findings.append((str(path.relative_to(ROOT)), number, "success-literal", line.strip()))
    for path, number, kind, line in findings:
        print(f"{path}:{number}: {kind}: {line}")
    print(f"static success checker: {len(findings)} finding(s)")
    return 1 if args.strict and findings else 0


if __name__ == "__main__":
    raise SystemExit(main())
