#!/usr/bin/env python3
from __future__ import annotations

import json
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]


def load_json(path: Path) -> dict:
    return json.loads(path.read_text(encoding="utf-8"))


def require(condition: bool, message: str) -> None:
    if not condition:
        raise ValueError(message)


def validate_task(task: dict) -> None:
    required = {"id", "agent_type", "target_ref", "technique", "mode", "requested_by"}
    require(required <= task.keys(), "simulation task is missing required fields")
    require(task["target_ref"].startswith("fixture://"), "simulation task target must use fixture://")
    require(task["mode"] in {"observe", "simulate"}, "simulation task mode is not allowed")
    require(all(isinstance(task[field], str) and task[field].strip() for field in required - {"mode"}), "simulation task has empty string fields")


def validate_result(task: dict, result: dict) -> None:
    required = {"task_id", "status", "simulation", "authorized", "evidence_refs"}
    require(required <= result.keys(), "simulation result is missing required fields")
    require(result["task_id"] == task["id"], "simulation result task ID does not match task")
    require(result["status"] in {"accepted", "rejected", "completed", "failed"}, "simulation result status is not allowed")
    require(result["simulation"] is True, "simulation result must be marked simulation")
    require(result["authorized"] is True, "repository fixture result must be authorized")
    require(all(isinstance(ref, str) and ref.startswith("fixture://") for ref in result["evidence_refs"]), "evidence reference is not fixture-backed")


def validate_lifecycle_schema(schema: dict) -> None:
    required = {"task_id", "state", "version", "updated_at"}
    properties = schema.get("properties", {})
    require(set(schema.get("required", [])) == required, "lifecycle schema required fields are incorrect")
    require(properties.get("state", {}).get("enum") == [
        "created", "approved", "queued", "running", "completed",
        "failed", "rejected", "cancelled", "expired", "compensated",
    ], "lifecycle state set is incorrect")


def main() -> int:
    openapi = (ROOT / "contracts/api/openapi.yaml").read_text(encoding="utf-8")
    require("/healthz:" in openapi, "OpenAPI health path is missing")
    require("/v1/simulation/tasks:" in openapi, "OpenAPI simulation path is missing")
    task = load_json(ROOT / "simulation/fixtures/tasks/recon-http.json")
    result = load_json(ROOT / "simulation/fixtures/reports/recon-http-result.json")
    lifecycle = load_json(ROOT / "contracts/task/lifecycle.schema.json")
    validate_task(task)
    validate_result(task, result)
    validate_lifecycle_schema(lifecycle)
    print("safe contracts validated")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
