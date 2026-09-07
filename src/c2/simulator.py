"""Local agent simulator for protocol compatibility tests."""

from __future__ import annotations

from dataclasses import dataclass
from typing import Any

from .implant_windows import execute_allowlisted_task
from .storage import Store


@dataclass
class AgentSimulator:
    store: Store
    agent_id: str = "sim-agent"
    hostname: str = "simulator"
    os_name: str = "linux"
    arch: str = "amd64"

    def register(self) -> None:
        self.store.upsert_agent(self.agent_id, self.hostname, self.os_name, self.arch)

    def run_once(self) -> dict[str, Any] | None:
        task = self.store.claim_task(self.agent_id)
        if task is None:
            return None
        result = execute_allowlisted_task({"type": task.task_type, **task.payload})
        self.store.record_result(task.id, self.agent_id, result)
        return {"task_id": task.id, "result": result}


__all__ = ["AgentSimulator"]
