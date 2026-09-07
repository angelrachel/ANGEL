from __future__ import annotations

import json
from pathlib import Path
from urllib.request import Request, urlopen

import pytest

from src.auth import AuthError, authenticate, issue_token
from src.c2.server import build_server
from src.c2.storage import Store


def test_rbac_token_and_expiry() -> None:
    token = issue_token("auditor-1", "auditor", ttl=10, now=100)
    principal = authenticate(token, now=101)
    assert principal.subject == "auditor-1"
    assert principal.can("evidence")
    assert not principal.can("write")
    with pytest.raises(AuthError):
        authenticate(token[:-1] + "0", now=101)
    with pytest.raises(AuthError):
        authenticate(token, now=111)


def test_scope_evidence_report_storage(tmp_path: Path) -> None:
    store = Store(tmp_path / "records.db")
    assert store.schema_version() == store.SCHEMA_VERSION
    scope = store.create_scope("engagement", ["example.test"], ["/api/"])
    evidence = store.add_evidence(scope.id, "request", "tester", {"status": 200}, "a" * 64)
    report = store.add_report(scope.id, "Report", {"findings": []})
    assert store.get_scope(scope.id) == scope
    assert evidence.scope_id == scope.id
    assert report.scope_id == scope.id
    assert evidence.payload["status"] == 200
    with pytest.raises(ValueError, match="SHA-256"):
        store.add_evidence(scope.id, "response", "tester", {}, "invalid")
    with pytest.raises(ValueError, match="paths"):
        store.create_scope("bad", ["example.test"], ["api"])
    assert store.list_audit_events(event_type="scope.created", actor="operator")


def test_store_rejects_newer_schema(tmp_path: Path) -> None:
    import sqlite3

    database = tmp_path / "future.db"
    with sqlite3.connect(database) as connection:
        connection.execute("PRAGMA user_version = 999")
    with pytest.raises(RuntimeError, match="newer"):
        Store(database)


def test_rbac_api_endpoints(tmp_path: Path) -> None:
    server = build_server(host="127.0.0.1", port=0, database=str(tmp_path / "api.db"))
    server_thread = __import__("threading").Thread(target=server.serve_forever, daemon=True)
    server_thread.start()
    base = f"http://127.0.0.1:{server.server_port}"
    token = issue_token("operator-1", "operator")
    headers = {"Content-Type": "application/json", "Authorization": f"Bearer {token}"}

    def post(path: str, body: dict) -> dict:
        request = Request(base + path, data=json.dumps(body).encode(), headers=headers, method="POST")
        with urlopen(request, timeout=2) as response:
            return json.loads(response.read())

    scope = post("/scope", {"name": "test", "hosts": ["example.test"], "paths": ["/"]})
    assert scope["status"] == "active"
    evidence = post(
        "/evidence",
        {"scope_id": scope["id"], "evidence_type": "response", "payload": {"status": 200}, "record_hash": "b" * 64},
    )
    assert evidence["id"] > 0
    assert post("/report", {"scope_id": scope["id"], "title": "Test"})["title"] == "Test"
    server.shutdown()
    server_thread.join(timeout=2)
