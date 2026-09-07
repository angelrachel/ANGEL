"""Fail-closed policy for agent tasks."""

from __future__ import annotations

from typing import Any

ALLOWED_TASKS = {
    "heartbeat",
    "self_test",
    "get_config",
    "get_capabilities",
    "collect_synthetic_inventory",
    "submit_synthetic_result",
    "rotate_key",
    "shutdown_agent",
}


def validate_task(task_type: str, payload: dict[str, Any]) -> None:
    if task_type not in ALLOWED_TASKS:
        raise ValueError("task type is not allowed")
    if not isinstance(payload, dict):
        raise ValueError("task payload must be an object")
    if len(str(payload)) > 8192:
        raise ValueError("task payload is too large")


__all__ = ["ALLOWED_TASKS", "validate_task"]
