from __future__ import annotations

import json

import pytest

from src.orchestrator.replay import ReplayError, ReplayLog


def test_replay_log_round_trip() -> None:
    log = ReplayLog()
    log.append("wf-1", "running", {"task": "inventory"}, {})
    log.append("wf-1", "succeeded", {}, {"count": 2})
    restored = ReplayLog.from_json(log.to_json())
    assert restored.records == log.records


def test_replay_log_rejects_sequence_tampering() -> None:
    value = json.dumps([{"sequence": 2, "workflow_id": "wf", "state": "running", "input": {}, "output": {}}])
    with pytest.raises(ReplayError, match="monotonic"):
        ReplayLog.from_json(value)
