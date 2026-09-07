from __future__ import annotations

from concurrent.futures import ThreadPoolExecutor
from pathlib import Path

from src.c2.storage import Store


def test_task_claim_is_single_winner(tmp_path: Path) -> None:
    store = Store(tmp_path / "concurrency.db")
    store.upsert_agent("a1", "host-1", "linux", "amd64")
    store.upsert_agent("a2", "host-2", "linux", "amd64")
    task = store.enqueue_task("a1", "self_test", {})

    def claim(agent_id: str):
        return store.claim_task(agent_id)

    with ThreadPoolExecutor(max_workers=2) as executor:
        results = list(executor.map(claim, ["a1", "a1"]))
    winners = [result for result in results if result is not None]
    assert len(winners) == 1
    assert winners[0].id == task.id
