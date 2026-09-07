"""ANGEL agent client for authorized lab control-plane tests."""

from __future__ import annotations

import json
import platform
import socket
import time
import urllib.request
import uuid
from typing import Any

from .policy import ALLOWED_TASKS


def get_agent_id() -> str:
    return str(uuid.uuid4())


def _post(c2_url: str, path: str, payload: dict[str, Any], operator_key: str) -> dict[str, Any]:
    request = urllib.request.Request(
        f"{c2_url.rstrip('/')}{path}",
        data=json.dumps(payload).encode(),
        headers={"Content-Type": "application/json", "X-ANGEL-Key": operator_key},
        method="POST",
    )
    with urllib.request.urlopen(request, timeout=10) as response:  # nosec B310 - URL is operator-configured
        value = json.loads(response.read().decode())
        if not isinstance(value, dict):
            raise ValueError("invalid server response")
        return value


def register(agent_id: str, c2_url: str, operator_key: str = "development-operator-key") -> dict[str, Any]:
    return _post(
        c2_url,
        "/register",
        {"agent_id": agent_id, "hostname": socket.gethostname(), "os": platform.system(), "arch": platform.machine()},
        operator_key,
    )


def get_task(agent_id: str, c2_url: str, operator_key: str = "development-operator-key") -> dict[str, Any] | None:
    response = _post(c2_url, "/task/claim", {"agent_id": agent_id}, operator_key)
    task = response.get("task")
    return task if isinstance(task, dict) else None


def send_result(
    agent_id: str,
    c2_url: str,
    task_id: int,
    payload: dict[str, Any],
    operator_key: str = "development-operator-key",
) -> dict[str, Any]:
    return _post(c2_url, "/result", {"agent_id": agent_id, "task_id": task_id, "payload": payload}, operator_key)


def execute_allowlisted_task(task: dict[str, Any]) -> dict[str, Any]:
    task_type = task.get("type")
    if task_type not in ALLOWED_TASKS:
        raise ValueError("task type is not allowed")
    if task_type == "heartbeat":
        return {"status": "alive", "timestamp": int(time.time())}
    if task_type == "self_test":
        return {"status": "passed", "agent": "angel-agent"}
    if task_type == "get_capabilities":
        return {"tasks": sorted(ALLOWED_TASKS)}
    if task_type == "get_config":
        return {"agent_mode": "synthetic", "heartbeat_seconds": 30, "max_task_runtime_seconds": 60}
    if task_type == "collect_synthetic_inventory":
        return {"hostname": "synthetic-host", "os": "synthetic-os", "source": "test-fixture"}
    if task_type == "submit_synthetic_result":
        result = task.get("result")
        if not isinstance(result, dict) or len(str(result)) > 4096:
            raise ValueError("synthetic result must be a bounded object")
        return {"status": "accepted", "fields": sorted(result)}
    if task_type == "rotate_key":
        return {"status": "rotation-requested", "key_id": str(task.get("key_id", "next"))}
    if task_type == "shutdown_agent":
        return {"status": "shutdown-requested"}
    return {"status": "accepted", "task": task_type}


def run_windows_implant() -> None:
    """Run a bounded lab agent loop; no arbitrary command execution is supported."""
    agent_id = get_agent_id()
    c2_url = "http://127.0.0.1:8000"
    register(agent_id, c2_url)
    for _ in range(1):
        task = get_task(agent_id, c2_url)
        if task and isinstance(task.get("id"), int):
            result = execute_allowlisted_task({"type": task.get("type")})
            send_result(agent_id, c2_url, task["id"], result)


__all__ = ["get_agent_id", "register", "get_task", "send_result", "execute_allowlisted_task", "run_windows_implant"]
