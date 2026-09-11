from __future__ import annotations

import json
import sys
from pathlib import Path
from typing import TypedDict

from langgraph.graph import END, START, StateGraph


class Task(TypedDict, total=False):
    id: str
    agent_type: str
    target_ref: str
    technique: str
    mode: str
    requested_by: str
    approval_id: str | None


class Plan(TypedDict, total=False):
    task_id: str
    status: str
    simulation: bool
    authorized: bool
    target_ref: str
    technique: str
    steps: list[str]
    evidence_refs: list[str]
    message: str


def validate_task(task: Task) -> None:
    required = ("id", "agent_type", "target_ref", "technique", "mode", "requested_by")
    if any(not str(task.get(field, "")).strip() for field in required):
        raise ValueError("task identity and simulation fields are required")
    if not task["target_ref"].startswith("fixture://"):
        raise ValueError("target_ref must use fixture://")
    if task["mode"] not in {"observe", "simulate"}:
        raise ValueError("mode must be observe or simulate")


def build_plan(task: Task) -> Plan:
    validate_task(task)
    return {
        "task_id": task["id"],
        "status": "accepted",
        "simulation": True,
        "authorized": False,
        "target_ref": task["target_ref"],
        "technique": task["technique"],
        "steps": ["validate_fixture_reference", "await_go_policy_gate", "record_simulation_evidence"],
        "evidence_refs": [],
        "message": "plan created; policy approval is required before simulation acceptance",
    }


def approve_plan(plan: Plan, approval_id: str) -> Plan:
    """Move a fixture plan to authorized simulation only with an explicit approval ID."""
    if not str(plan.get("task_id", "")).strip():
        raise ValueError("plan task_id is required")
    if not str(plan.get("target_ref", "")).startswith("fixture://"):
        raise ValueError("plan target must use fixture://")
    if not str(approval_id).strip():
        raise ValueError("approval_id is required")
    approved = dict(plan)
    approved["authorized"] = True
    approved["status"] = "approved"
    approved["steps"] = [*plan.get("steps", []), "record_approval", "execute_bounded_simulation"]
    approved["message"] = "explicit approval recorded; bounded simulation may proceed"
    return approved


def run(task: Task) -> Plan:
    graph = StateGraph(dict)
    graph.add_node("plan", lambda state: {"plan": build_plan(state["task"])})
    graph.add_edge(START, "plan")
    graph.add_edge("plan", END)
    result = graph.compile().invoke({"task": task})
    return result["plan"]


def main() -> int:
    if len(sys.argv) != 2:
        print(f"usage: {Path(sys.argv[0]).name} TASK_JSON", file=sys.stderr)
        return 2
    try:
        task = json.loads(Path(sys.argv[1]).read_text(encoding="utf-8"))
        print(json.dumps(run(task), indent=2, sort_keys=True))
    except (OSError, TypeError, ValueError, json.JSONDecodeError) as error:
        print(str(error), file=sys.stderr)
        return 1
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
