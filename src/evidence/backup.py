"""Safe SQLite backup and retention helpers."""

from __future__ import annotations

import sqlite3
from pathlib import Path


class BackupError(RuntimeError):
    """Raised when a backup cannot be created or validated."""


def backup_database(source: str | Path, destination: str | Path) -> Path:
    source_path = Path(source)
    destination_path = Path(destination)
    if not source_path.is_file():
        raise BackupError("source database does not exist")
    destination_path.parent.mkdir(parents=True, exist_ok=True)
    try:
        with sqlite3.connect(source_path) as source_db, sqlite3.connect(destination_path) as target_db:
            source_db.backup(target_db)
            result = target_db.execute("PRAGMA integrity_check").fetchone()
    except sqlite3.Error as exc:
        raise BackupError("database backup failed") from exc
    if result != ("ok",):
        raise BackupError("backup integrity check failed")
    return destination_path


def prune_backups(directory: str | Path, keep: int = 5) -> list[Path]:
    if keep < 1:
        raise ValueError("keep must be positive")
    candidates = sorted(Path(directory).glob("*.db"), key=lambda path: path.stat().st_mtime, reverse=True)
    removed: list[Path] = []
    for path in candidates[keep:]:
        path.unlink()
        removed.append(path)
    return removed


__all__ = ["BackupError", "backup_database", "prune_backups"]
