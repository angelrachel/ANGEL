"""Serializable workflow replay records."""

from __future__ import annotations

import json
from dataclasses import asdict, dataclass


class ReplayError(ValueError):
    """Raised when a replay record is malformed."""


@dataclass(frozen=True)
class ReplayRecord:
    sequence: int
    workflow_id: str
    state: str
    input: dict[str, object]
    output: dict[str, object]


class ReplayLog:
    def __init__(self) -> None:
        self.records: list[ReplayRecord] = []

    def append(
        self, workflow_id: str, state: str, input_data: dict[str, object], output: dict[str, object]
    ) -> ReplayRecord:
        if not workflow_id.strip() or not state.strip():
            raise ReplayError("workflow id and state are required")
        record = ReplayRecord(len(self.records) + 1, workflow_id, state, input_data, output)
        self.records.append(record)
        return record

    def to_json(self) -> str:
        return json.dumps([asdict(record) for record in self.records], sort_keys=True)

    @classmethod
    def from_json(cls, value: str) -> ReplayLog:
        try:
            raw = json.loads(value)
            if not isinstance(raw, list):
                raise ReplayError("replay log must be an array")
            log = cls()
            for item in raw:
                if not isinstance(item, dict):
                    raise ReplayError("replay record must be an object")
                record = ReplayRecord(
                    int(item["sequence"]),
                    str(item["workflow_id"]),
                    str(item["state"]),
                    item["input"],
                    item["output"],
                )
                if record.sequence != len(log.records) + 1:
                    raise ReplayError("replay sequence is not monotonic")
                log.records.append(record)
            return log
        except ReplayError:
            raise
        except (KeyError, TypeError, ValueError, json.JSONDecodeError) as exc:
            raise ReplayError("malformed replay log") from exc


__all__ = ["ReplayError", "ReplayLog", "ReplayRecord"]
