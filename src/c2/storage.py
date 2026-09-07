"""SQLite persistence for ANGEL control-plane state."""

from __future__ import annotations

import json
import sqlite3
import time
from dataclasses import dataclass
from pathlib import Path
from typing import Any

from .policy import validate_task
from .records import EvidenceRecordRow, ReportRecord, ScopeRecord


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
                CREATE TABLE IF NOT EXISTS scopes (
                    id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT NOT NULL,
                    hosts TEXT NOT NULL, paths TEXT NOT NULL, expires_at INTEGER,
                    status TEXT NOT NULL
                );
                CREATE TABLE IF NOT EXISTS evidence (
                    id INTEGER PRIMARY KEY AUTOINCREMENT, scope_id INTEGER NOT NULL,
                    evidence_type TEXT NOT NULL, actor TEXT NOT NULL, payload TEXT NOT NULL,
                    record_hash TEXT NOT NULL, created_at INTEGER NOT NULL,
                    FOREIGN KEY(scope_id) REFERENCES scopes(id)
                );
                CREATE TABLE IF NOT EXISTS reports (
                    id INTEGER PRIMARY KEY AUTOINCREMENT, scope_id INTEGER NOT NULL,
                    title TEXT NOT NULL, payload TEXT NOT NULL, created_at INTEGER NOT NULL,
                    FOREIGN KEY(scope_id) REFERENCES scopes(id)
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

    def create_scope(self, name: str, hosts: list[str], paths: list[str], expires_at: int | None = None) -> ScopeRecord:
        if not name.strip() or not hosts or not paths:
            raise ValueError("scope name, hosts, and paths are required")
        with self._connect() as db:
            cursor = db.execute(
                "INSERT INTO scopes(name, hosts, paths, expires_at, status) VALUES (?, ?, ?, ?, 'active')",
                (name, json.dumps(hosts), json.dumps(paths), expires_at),
            )
            if cursor.lastrowid is None:
                raise RuntimeError("scope insert did not return an id")
            scope_id = cursor.lastrowid
        self.audit("scope.created", "operator", str(scope_id), {"name": name})
        return ScopeRecord(scope_id, name, hosts, paths, expires_at, "active")

    def get_scope(self, scope_id: int) -> ScopeRecord | None:
        with self._connect() as db:
            row = db.execute("SELECT * FROM scopes WHERE id = ?", (scope_id,)).fetchone()
        if row is None:
            return None
        return ScopeRecord(
            row["id"], row["name"], json.loads(row["hosts"]), json.loads(row["paths"]), row["expires_at"], row["status"]
        )

    def add_evidence(
        self, scope_id: int, evidence_type: str, actor: str, payload: dict[str, Any], record_hash: str
    ) -> EvidenceRecordRow:
        now = int(time.time())
        with self._connect() as db:
            cursor = db.execute(
                "INSERT INTO evidence(scope_id, evidence_type, actor, payload, record_hash, created_at) "
                "VALUES (?, ?, ?, ?, ?, ?)",
                (scope_id, evidence_type, actor, json.dumps(payload, sort_keys=True), record_hash, now),
            )
            if cursor.lastrowid is None:
                raise RuntimeError("evidence insert did not return an id")
            evidence_id = cursor.lastrowid
        self.audit("evidence.added", actor, str(evidence_id), {"scope_id": scope_id, "type": evidence_type})
        return EvidenceRecordRow(evidence_id, scope_id, evidence_type, actor, payload, record_hash, now)

    def add_report(self, scope_id: int, title: str, payload: dict[str, Any]) -> ReportRecord:
        now = int(time.time())
        with self._connect() as db:
            cursor = db.execute(
                "INSERT INTO reports(scope_id, title, payload, created_at) VALUES (?, ?, ?, ?)",
                (scope_id, title, json.dumps(payload, sort_keys=True), now),
            )
            if cursor.lastrowid is None:
                raise RuntimeError("report insert did not return an id")
            report_id = cursor.lastrowid
        self.audit("report.created", "operator", str(report_id), {"scope_id": scope_id, "title": title})
        return ReportRecord(report_id, scope_id, title, payload, now)


__all__ = ["Agent", "Task", "Store"]
