"""SQLite persistence for ANGEL control-plane state."""

from __future__ import annotations

import json
import sqlite3
import time
from dataclasses import dataclass
from pathlib import Path
from typing import Any

from .policy import validate_task


@dataclass(frozen=True)
class Agent:
    id: str
    hostname: str
    os: str
    arch: str
    last_seen: int
    status: str


@dataclass(frozen=True)
class Task:
    id: int
    agent_id: str
    task_type: str
    payload: dict[str, Any]
    status: str
    created_at: int
    expires_at: int


class Store:
    def __init__(self, path: str | Path = "c2.db") -> None:
        self.path = str(path)
        self._initialize()

    def _connect(self) -> sqlite3.Connection:
        connection = sqlite3.connect(self.path)
        connection.row_factory = sqlite3.Row
        connection.execute("PRAGMA foreign_keys = ON")
        return connection

    def _initialize(self) -> None:
        with self._connect() as db:
            db.executescript(
                """
                CREATE TABLE IF NOT EXISTS agents (
                    id TEXT PRIMARY KEY, hostname TEXT NOT NULL, os TEXT NOT NULL,
                    arch TEXT NOT NULL, last_seen INTEGER NOT NULL, status TEXT NOT NULL
                );
                CREATE TABLE IF NOT EXISTS tasks (
                    id INTEGER PRIMARY KEY AUTOINCREMENT, agent_id TEXT NOT NULL,
                    task_type TEXT NOT NULL, payload TEXT NOT NULL, status TEXT NOT NULL,
                    created_at INTEGER NOT NULL, expires_at INTEGER NOT NULL,
                    FOREIGN KEY(agent_id) REFERENCES agents(id)
                );
                CREATE TABLE IF NOT EXISTS results (
                    id INTEGER PRIMARY KEY AUTOINCREMENT, task_id INTEGER NOT NULL,
                    agent_id TEXT NOT NULL, payload TEXT NOT NULL, created_at INTEGER NOT NULL,
                    FOREIGN KEY(task_id) REFERENCES tasks(id), FOREIGN KEY(agent_id) REFERENCES agents(id)
                );
                CREATE TABLE IF NOT EXISTS audit_events (
                    id INTEGER PRIMARY KEY AUTOINCREMENT, event_type TEXT NOT NULL,
                    actor TEXT NOT NULL, subject TEXT NOT NULL, details TEXT NOT NULL,
                    created_at INTEGER NOT NULL
                );
                CREATE INDEX IF NOT EXISTS idx_tasks_agent_status ON tasks(agent_id, status, expires_at);
                CREATE INDEX IF NOT EXISTS idx_audit_created ON audit_events(created_at);
                """
            )

    def upsert_agent(self, agent_id: str, hostname: str, os_name: str, arch: str) -> Agent:
        now = int(time.time())
        with self._connect() as db:
            db.execute(
                """INSERT INTO agents(id, hostname, os, arch, last_seen, status) VALUES (?, ?, ?, ?, ?, 'online')
                   ON CONFLICT(id) DO UPDATE SET hostname=excluded.hostname, os=excluded.os,
                   arch=excluded.arch, last_seen=excluded.last_seen, status='online'""",
                (agent_id, hostname, os_name, arch, now),
            )
        self.audit("agent.registered", agent_id, agent_id, {"hostname": hostname, "os": os_name, "arch": arch})
        return Agent(agent_id, hostname, os_name, arch, now, "online")

    def get_agent(self, agent_id: str) -> Agent | None:
        with self._connect() as db:
            row = db.execute("SELECT * FROM agents WHERE id = ?", (agent_id,)).fetchone()
        return None if row is None else Agent(**dict(row))

    def enqueue_task(self, agent_id: str, task_type: str, payload: dict[str, Any], ttl: int = 300) -> Task:
        validate_task(task_type, payload)
        now = int(time.time())
        with self._connect() as db:
            cursor = db.execute(
                "INSERT INTO tasks(agent_id, task_type, payload, status, created_at, expires_at) "
                "VALUES (?, ?, ?, 'pending', ?, ?)",
                (agent_id, task_type, json.dumps(payload, sort_keys=True), now, now + ttl),
            )
            if cursor.lastrowid is None:
                raise RuntimeError("task insert did not return an id")
            task_id = cursor.lastrowid
        self.audit("task.queued", "operator", str(task_id), {"agent_id": agent_id, "task_type": task_type})
        return Task(task_id, agent_id, task_type, payload, "pending", now, now + ttl)

    def claim_task(self, agent_id: str) -> Task | None:
        now = int(time.time())
        with self._connect() as db:
            db.execute(
                "UPDATE tasks SET status='expired' WHERE agent_id=? AND status='pending' AND expires_at < ?",
                (agent_id, now),
            )
            row = db.execute(
                "SELECT * FROM tasks WHERE agent_id=? AND status='pending' ORDER BY id LIMIT 1",
                (agent_id,),
            ).fetchone()
            if row is None:
                return None
            db.execute("UPDATE tasks SET status='assigned' WHERE id=?", (row["id"],))
        return Task(
            row["id"],
            row["agent_id"],
            row["task_type"],
            json.loads(row["payload"]),
            "assigned",
            row["created_at"],
            row["expires_at"],
        )

    def record_result(self, task_id: int, agent_id: str, payload: dict[str, Any]) -> None:
        now = int(time.time())
        with self._connect() as db:
            db.execute(
                "INSERT INTO results(task_id, agent_id, payload, created_at) VALUES (?, ?, ?, ?)",
                (task_id, agent_id, json.dumps(payload, sort_keys=True), now),
            )
            db.execute("UPDATE tasks SET status='completed' WHERE id=? AND agent_id=?", (task_id, agent_id))
        self.audit("task.completed", agent_id, str(task_id), {"keys": sorted(payload)})

    def audit(self, event_type: str, actor: str, subject: str, details: dict[str, Any]) -> None:
        with self._connect() as db:
            db.execute(
                "INSERT INTO audit_events(event_type, actor, subject, details, created_at) VALUES (?, ?, ?, ?, ?)",
                (event_type, actor, subject, json.dumps(details, sort_keys=True), int(time.time())),
            )


__all__ = ["Agent", "Task", "Store"]
