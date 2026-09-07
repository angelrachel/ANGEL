from __future__ import annotations

from pathlib import Path

import pytest

from src.c2.storage import Store
from src.evidence.backup import BackupError, BackupManager, backup_database, prune_backups


def test_backup_database_and_prune(tmp_path: Path) -> None:
    source = tmp_path / "source.db"
    Store(source).upsert_agent("a1", "host", "linux", "amd64")
    backup = backup_database(source, tmp_path / "backups" / "one.db")
    assert backup.is_file()
    for index in range(2, 5):
        backup_database(source, tmp_path / "backups" / f"{index}.db")
    removed = prune_backups(tmp_path / "backups", keep=2)
    assert len(removed) == 2
    assert len(list((tmp_path / "backups").glob("*.db"))) == 2


def test_backup_requires_source(tmp_path: Path) -> None:
    with pytest.raises(BackupError, match="does not exist"):
        backup_database(tmp_path / "missing.db", tmp_path / "out.db")
    with pytest.raises(ValueError):
        prune_backups(tmp_path, keep=0)


def test_backup_manager_creates_timestamped_snapshot(tmp_path: Path) -> None:
    source = tmp_path / "source.db"
    Store(source)
    manager = BackupManager(source, tmp_path / "snapshots", keep=1)
    assert manager.snapshot(100).name == "c2-100.db"
    assert manager.snapshot(200).name == "c2-200.db"
    assert [path.name for path in (tmp_path / "snapshots").glob("*.db")] == ["c2-200.db"]
