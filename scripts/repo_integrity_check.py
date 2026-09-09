#!/usr/bin/env python3
"""Read-only repository hygiene checks for CI and release review."""
from __future__ import annotations

import subprocess
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]


def tracked_files() -> list[Path]:
    output = subprocess.check_output(["git", "ls-files", "-z"], cwd=ROOT)
    return [ROOT / item for item in output.decode().split("\0") if item]


def is_binary(path: Path) -> bool:
    try:
        data = path.read_bytes()[:4096]
    except OSError:
        return False
    return b"\0" in data or data.startswith((b"\x7fELF", b"MZ"))


def main() -> int:
    failures: list[str] = []
    for path in tracked_files():
        relative = path.relative_to(ROOT)
        if not path.exists():
            failures.append(f"tracked path missing from working tree: {relative}")
            continue
        if path.is_file() and path.stat().st_size == 0 and path.name not in {".gitkeep"}:
            failures.append(f"empty tracked file: {relative}")
        if path.is_file() and is_binary(path) and path.suffix not in {".png", ".jpg", ".jpeg", ".gif", ".ico", ".pdf"}:
            failures.append(f"tracked binary artifact: {relative}")
        if relative.name.endswith((".tfstate", ".tfstate.backup")):
            failures.append(f"tracked Terraform state: {relative}")

    if failures:
        print("repository integrity check failed:", file=sys.stderr)
        for failure in failures:
            print(f"- {failure}", file=sys.stderr)
        return 1
    print(f"repository integrity check passed ({len(tracked_files())} tracked files)")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
