#!/usr/bin/env python3
"""Read-only repository integrity checks for ANGEL.

This checker intentionally does not execute attack modules or modify repository files.
It reports deterministic hygiene/integration risks for local and CI use.
"""
from __future__ import annotations

import ast
import hashlib
import json
import pathlib
import re
import subprocess
import sys
from dataclasses import asdict, dataclass

ROOT = pathlib.Path(__file__).resolve().parents[1]
SOURCE_ROOTS = (ROOT / "src", ROOT / "tests")
IGNORED_DIRS = {".git", "__pycache__", ".mypy_cache", ".pytest_cache", ".ruff_cache", ".venv", "venv"}
LOCAL_MODULE_RE = re.compile(r"^(?:src(?:\.|$)|[A-Za-z_][A-Za-z0-9_]*(?:\.|$))")


@dataclass
class Finding:
    severity: str
    category: str
    path: str
    line: int | None
    message: str


def rel(path: pathlib.Path) -> str:
    return path.relative_to(ROOT).as_posix()


def tracked_files() -> set[str]:
    result = subprocess.run(
        ["git", "-C", str(ROOT), "ls-files"], capture_output=True, text=True, check=True
    )
    return {line for line in result.stdout.splitlines() if line}


def python_modules() -> tuple[dict[str, pathlib.Path], set[str]]:
    modules: dict[str, pathlib.Path] = {}
    packages: set[str] = set()
    for base in SOURCE_ROOTS:
        if not base.exists():
            continue
        for directory in [base, *base.rglob("*")]:
            if directory.is_dir() and not any(part in IGNORED_DIRS for part in directory.parts):
                packages.add(".".join(directory.relative_to(ROOT).parts))
        for path in base.rglob("*.py"):
            if any(part in IGNORED_DIRS for part in path.parts):
                continue
            module = ".".join(path.relative_to(ROOT).with_suffix("").parts)
            if module.endswith(".__init__"):
                module = module[: -len(".__init__")]
            modules[module] = path
    return modules, packages


def resolve_local_module(module: str, modules: dict[str, pathlib.Path], packages: set[str]) -> bool:
    candidates = [module, module.replace(".", "/")]
    return any(candidate in modules or candidate in packages for candidate in candidates)


def collect() -> list[Finding]:
    findings: list[Finding] = []
    tracked = tracked_files()
    modules, packages = python_modules()

    for path in ROOT.rglob("*"):
        if not path.is_file() or any(part in IGNORED_DIRS for part in path.parts):
            continue
        relative = rel(path)
        if path.stat().st_size == 0 and relative in tracked:
            findings.append(Finding("high", "empty-tracked-file", relative, None, "Tracked file is empty."))
        if path.suffix in {".pyc", ".pyo"} and relative in tracked:
            findings.append(Finding("high", "generated-tracked-file", relative, None, "Generated Python bytecode is tracked."))
        if path.name in {"server", "agent", "server.exe", "agent.exe"} and relative in tracked:
            digest = hashlib.sha256(path.read_bytes()).hexdigest()
            findings.append(Finding("high", "tracked-binary", relative, None, f"Tracked executable artifact; sha256={digest}."))
        if path.name.endswith(".tfstate") or path.name.endswith(".tfstate.backup"):
            findings.append(Finding("high", "terraform-state", relative, None, "Terraform state must not be committed."))

    for path in (ROOT / "src").rglob("*.py") if (ROOT / "src").exists() else []:
        if any(part in IGNORED_DIRS for part in path.parts):
            continue
        try:
            tree = ast.parse(path.read_text(encoding="utf-8"), filename=str(path))
        except SyntaxError as exc:
            findings.append(Finding("high", "python-syntax", rel(path), exc.lineno, str(exc)))
            continue
        for node in ast.walk(tree):
            module: str | None = None
            line = getattr(node, "lineno", None)
            if isinstance(node, ast.Import):
                for alias in node.names:
                    module = alias.name
                    if module.startswith("src") and not resolve_local_module(module, modules, packages):
                        findings.append(Finding("high", "broken-local-import", rel(path), line, f"Missing local module: {module}"))
            elif isinstance(node, ast.ImportFrom):
                if node.level:
                    base_parts = path.relative_to(ROOT).with_suffix("").parts[:-1]
                    prefix = list(base_parts[: max(0, len(base_parts) - node.level + 1)])
                    if node.module:
                        prefix.extend(node.module.split("."))
                    module = ".".join(prefix)
                elif node.module:
                    module = node.module
                if module and (module.startswith("src") or module.startswith("tests")) and not resolve_local_module(module, modules, packages):
                    findings.append(Finding("high", "broken-local-import", rel(path), line, f"Missing local module: {module}"))

    manifest = ROOT / "requirements.txt"
    if manifest.exists():
        for index, raw in enumerate(manifest.read_text(encoding="utf-8").splitlines(), 1):
            item = raw.strip()
            if item and not item.startswith("#") and not re.search(r"(==|@|sha256=)", item):
                findings.append(Finding("medium", "unpinned-python-dependency", rel(manifest), index, f"Dependency is not hash/version pinned: {item}"))

    for filename in ("LICENSE", "LICENSE.md", "COPYING", "NOTICE"):
        if (ROOT / filename).exists():
            break
    else:
        findings.append(Finding("medium", "missing-provenance-document", "", None, "No root LICENSE, COPYING, or NOTICE file found."))

    return findings


def main() -> int:
    findings = collect()
    payload = {
        "repository": str(ROOT),
        "tracked_files": len(tracked_files()),
        "finding_count": len(findings),
        "findings": [asdict(item) for item in findings],
    }
    print(json.dumps(payload, indent=2, sort_keys=True))
    return 1 if any(item.severity == "high" for item in findings) else 0


if __name__ == "__main__":
    raise SystemExit(main())
