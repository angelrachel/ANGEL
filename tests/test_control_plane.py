from __future__ import annotations

import json
import threading
from pathlib import Path
from urllib.error import HTTPError
from urllib.request import Request, urlopen

import pytest

from src.auth import issue_token
from src.c2.crypto import (
    CryptoError,
    KeyPair,
    KeyRegistry,
    ReplayGuard,
    SessionCipher,
    derive_session_key,
    sign_message,
    verify_signature,
)
from src.c2.implant_windows import execute_allowlisted_task
from src.c2.policy import validate_task
from src.c2.server import build_server
from src.c2.storage import Store


def test_ecdh_and_authenticated_envelope() -> None:
    first, second = KeyPair.generate(), KeyPair.generate()
    left = derive_session_key(first.private_key, second.public_bytes)
    right = derive_session_key(second.private_key, first.public_bytes)
    assert left == right
    cipher = SessionCipher(left)
    envelope = cipher.encrypt({"type": "self_test"}, aad=b"agent-1")
    assert cipher.decrypt(envelope, aad=b"agent-1", replay=ReplayGuard()) == {"type": "self_test"}


def test_replay_and_tampering_are_rejected() -> None:
    cipher = SessionCipher(b"x" * 64)
    envelope = cipher.encrypt({"value": 1})
    guard = ReplayGuard()
    cipher.decrypt(envelope, replay=guard)
    with pytest.raises(CryptoError):
        cipher.decrypt(envelope, replay=guard)
    envelope["data"] = envelope["data"][:-2] + "AA"
    with pytest.raises(CryptoError):
        cipher.decrypt(envelope)


def test_key_registry_rotation_and_revocation() -> None:
    registry = KeyRegistry()
    registry.register("k1", b"x" * 64)
    assert registry.get("k1") == b"x" * 64
    registry.revoke("k1")
    with pytest.raises(CryptoError, match="revoked"):
        registry.get("k1")
    with pytest.raises(CryptoError, match="revoked"):
        registry.register("k1", b"y" * 64)


def test_signature_verification() -> None:
    key = KeyPair.generate()
    message = b"ANGEL test vector"
    signature = sign_message(key.private_key, message)
    assert verify_signature(key.private_key.public_key(), message, signature)
    assert not verify_signature(key.private_key.public_key(), b"tampered", signature)


def test_store_lifecycle(tmp_path: Path) -> None:
    store = Store(tmp_path / "test.db")
    store.upsert_agent("a1", "host", "linux", "amd64")
    task = store.enqueue_task("a1", "self_test", {})
    claimed = store.claim_task("a1")
    assert claimed is not None and claimed.id == task.id
    store.record_result(task.id, "a1", {"status": "passed"})
    assert store.get_agent("a1") is not None


def test_result_requires_assigned_task_owner(tmp_path: Path) -> None:
    store = Store(tmp_path / "test.db")
    store.upsert_agent("a1", "host", "linux", "amd64")
    store.upsert_agent("a2", "other", "linux", "amd64")
    task = store.enqueue_task("a1", "self_test", {})
    with pytest.raises(ValueError, match="assigned"):
        store.record_result(task.id, "a2", {"status": "passed"})
    with pytest.raises(ValueError, match="awaiting"):
        store.record_result(task.id, "a1", {"status": "passed"})
    store.claim_task("a1")
    store.record_result(task.id, "a1", {"status": "passed"})
    with pytest.raises(ValueError, match="awaiting"):
        store.record_result(task.id, "a1", {"status": "passed"})


def test_storage_maintenance_marks_expired_and_stale(tmp_path: Path) -> None:
    store = Store(tmp_path / "test.db")
    store.upsert_agent("a1", "host", "linux", "amd64")
    task = store.enqueue_task("a1", "self_test", {}, ttl=1)
    assert store.expire_tasks(now=task.expires_at + 1) == 1
    saved_task = store.get_task(task.id)
    assert saved_task is not None and saved_task.status == "expired"
    assert store.mark_stale_agents(stale_after=1, now=task.created_at + 2) == 1
    saved_agent = store.get_agent("a1")
    assert saved_agent is not None and saved_agent.status == "offline"
    with pytest.raises(ValueError, match="positive"):
        store.mark_stale_agents(stale_after=0)


def test_disallowed_task_is_rejected(tmp_path: Path) -> None:
    store = Store(tmp_path / "test.db")
    store.upsert_agent("a1", "host", "linux", "amd64")
    with pytest.raises(ValueError):
        store.enqueue_task("a1", "exec_shell", {})


def test_agent_executor_is_allowlisted() -> None:
    assert execute_allowlisted_task({"type": "self_test"})["status"] == "passed"
    assert execute_allowlisted_task({"type": "get_config"})["agent_mode"] == "synthetic"
    assert execute_allowlisted_task({"type": "submit_synthetic_result", "result": {"ok": True}})["status"] == "accepted"
    assert execute_allowlisted_task({"type": "rotate_key", "key_id": "k2"})["key_id"] == "k2"
    with pytest.raises(ValueError):
        execute_allowlisted_task({"type": "exec_shell"})


def test_policy_rejects_invalid_payloads() -> None:
    with pytest.raises(ValueError):
        validate_task("self_test", "not-an-object")  # type: ignore[arg-type]
    with pytest.raises(ValueError):
        validate_task("self_test", {"data": "x" * 9000})


def test_http_lifecycle(tmp_path: Path) -> None:
    server = build_server(host="127.0.0.1", port=0, database=str(tmp_path / "api.db"))
    thread = threading.Thread(target=server.serve_forever, daemon=True)
    thread.start()
    base = f"http://127.0.0.1:{server.server_port}"
    headers = {"Content-Type": "application/json", "X-ANGEL-Key": "development-operator-key"}

    def post(path: str, body: dict) -> dict:
        request = Request(base + path, data=json.dumps(body).encode(), headers=headers, method="POST")
        with urlopen(request, timeout=2) as response:
            return json.loads(response.read())

    with urlopen(base + "/healthz", timeout=2) as response:
        assert json.loads(response.read())["status"] == "ok"
    with urlopen(base + "/readyz", timeout=2) as response:
        assert json.loads(response.read())["status"] == "ready"
    with urlopen(base + "/openapi.json", timeout=2) as response:
        openapi = json.loads(response.read())
        assert openapi["openapi"] == "3.0.3"
        assert "/audit" in openapi["paths"]

    unauthenticated = Request(
        base + "/register",
        data=json.dumps({"agent_id": "unauth", "hostname": "h", "os": "linux", "arch": "amd64"}).encode(),
        headers={"Content-Type": "application/json"},
        method="POST",
    )
    with pytest.raises(HTTPError) as error:
        urlopen(unauthenticated, timeout=2)
    assert error.value.code == 401

    registration = {"agent_id": "a1", "hostname": "h", "os": "linux", "arch": "amd64"}
    assert post("/register", registration)["status"] == "registered"
    assert post("/heartbeat", registration)["status"] == "heartbeat_ack"
    with urlopen(Request(base + "/agents", headers=headers), timeout=2) as response:
        assert json.loads(response.read())["agents"][0]["id"] == "a1"
    queued = post("/task/queue", {"agent_id": "a1", "task_type": "self_test", "payload": {}})
    task = post("/task/claim", {"agent_id": "a1"})["task"]
    assert task["id"] == queued["task_id"]
    result = post("/result", {"agent_id": "a1", "task_id": task["id"], "payload": {"status": "passed"}})
    assert result["status"] == "recorded"
    with urlopen(Request(base + "/tasks?agent_id=a1", headers=headers), timeout=2) as response:
        assert json.loads(response.read())["tasks"][0]["status"] == "completed"
    with urlopen(Request(base + "/events?cursor=0", headers=headers), timeout=2) as response:
        assert json.loads(response.read())["events"]
    token_headers = {"Authorization": f"Bearer {issue_token('auditor', 'auditor')}"}
    with urlopen(Request(base + "/audit?limit=5", headers=token_headers), timeout=2) as response:
        assert json.loads(response.read())["events"]
    server.shutdown()
    thread.join(timeout=2)
