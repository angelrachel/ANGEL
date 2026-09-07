from __future__ import annotations

from pathlib import Path

from src.c2.simulator import AgentSimulator
from src.c2.storage import Store


def test_agent_simulator_registers_claims_and_submits(tmp_path: Path) -> None:
    store = Store(tmp_path / "sim.db")
    simulator = AgentSimulator(store)
    simulator.register()
    task = store.enqueue_task(simulator.agent_id, "self_test", {})
    result = simulator.run_once()
    assert result is not None
    assert result["task_id"] == task.id
    assert result["result"]["status"] == "passed"
    assert store.get_task(task.id).status == "completed"
    assert simulator.run_once() is None
