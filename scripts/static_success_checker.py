#!/usr/bin/env python3
"""Report suspicious success claims and ignored command errors without executing code."""
from __future__ import annotations

import argparse
import re
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
DEFAULT_ROOTS = [ROOT / "src" / "c2"]
PATTERNS = {
    "ignored-command-error": re.compile(r"\.Run\(\)"),
    "success-literal": re.compile(r'Status\s*:\s*"success"|return\s+true'),
}


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--strict", action="store_true", help="return non-zero when findings exist")
    args = parser.parse_args()
    findings: list[tuple[str, int, str, str]] = []
    for root in DEFAULT_ROOTS:
        for path in sorted(root.rglob("*.go")):
            for number, line in enumerate(path.read_text(errors="replace").splitlines(), 1):
                for kind, pattern in PATTERNS.items():
                    if pattern.search(line):
                        findings.append((str(path.relative_to(ROOT)), number, kind, line.strip()))
    for path, number, kind, line in findings:
        print(f"{path}:{number}: {kind}: {line}")
    print(f"static success checker: {len(findings)} finding(s)")
    return 1 if args.strict and findings else 0


if __name__ == "__main__":
    raise SystemExit(main())
