#!/usr/bin/env python3
"""Generate a conservative status inventory from source/test evidence."""
from __future__ import annotations

from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
SRC = ROOT / "src" / "c2"
OUT = ROOT / "docs" / "MODULE_INVENTORY.md"
HIGH_RISK = ("credential", "evasion", "persistence", "rootkit", "destruction", "lateral", "implant")


def classify(path: Path, files: list[Path]) -> str:
    relative = str(path.relative_to(SRC))
    has_tests = any(file.name.endswith("_test.go") for file in files)
    source_count = sum(file.suffix == ".go" for file in files)
    if any(token in relative for token in HIGH_RISK):
        return "partial/high-risk"
    if has_tests:
        return "foundation-tested"
    if source_count <= 1:
        return "skeleton/unverified"
    return "partial/unverified"


def main() -> None:
    rows: list[str] = []
    for directory in sorted(p for p in SRC.rglob("*") if p.is_dir()):
        files = [p for p in directory.iterdir() if p.is_file() and p.suffix == ".go"]
        if not files:
            continue
        rows.append(f"| `{directory.relative_to(ROOT)}` | {len(files)} | {sum(p.name.endswith('_test.go') for p in files)} | {classify(directory, files)} |")
    text = "# C2 Module Inventory\n\nGenerated from source/test evidence; a filename does not prove operational capability.\n\n| Package path | Go files | Test files | Conservative status |\n|---|---:|---:|---|\n" + "\n".join(rows) + "\n"
    OUT.write_text(text)
    print(f"wrote {OUT.relative_to(ROOT)} ({len(rows)} packages)")


if __name__ == "__main__":
    main()
