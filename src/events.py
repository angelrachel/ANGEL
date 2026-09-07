"""Versioned, validated event envelope shared by ANGEL layers."""

from __future__ import annotations

import json
import time
from dataclasses import asdict, dataclass
from typing import Any, ClassVar


class EventError(ValueError):
    """Raised when an event envelope is malformed."""


@dataclass(frozen=True)
class Event:
    name: str
    actor: str
    subject: str
    payload: dict[str, Any]
    created_at: int
    event_id: str
    version: int = 1

    MAX_PAYLOAD_BYTES: ClassVar[int] = 16_384

    def __post_init__(self) -> None:
        if self.version != 1:
            raise EventError("unsupported event version")
        if not self.name.strip() or not self.actor.strip() or not self.subject.strip():
            raise EventError("event name, actor, and subject are required")
        if not isinstance(self.payload, dict):
            raise EventError("event payload must be an object")
        if not self.event_id.strip():
            raise EventError("event id is required")
        try:
            encoded = json.dumps(self.payload, sort_keys=True, separators=(",", ":")).encode()
        except (TypeError, ValueError) as exc:
            raise EventError("event payload must be JSON serializable") from exc
        if len(encoded) > self.MAX_PAYLOAD_BYTES:
            raise EventError("event payload is too large")
        if self.created_at < 0:
            raise EventError("event timestamp must be non-negative")

    def to_dict(self) -> dict[str, Any]:
        return asdict(self)

    @classmethod
    def from_dict(cls, value: dict[str, Any]) -> Event:
        if not isinstance(value, dict):
            raise EventError("event must be an object")
        try:
            return cls(
                str(value["name"]),
                str(value["actor"]),
                str(value["subject"]),
                value["payload"],
                int(value["created_at"]),
                str(value["event_id"]),
                int(value.get("version", 1)),
            )
        except (KeyError, TypeError, ValueError) as exc:
            raise EventError("malformed event") from exc


def make_event(name: str, actor: str, subject: str, payload: dict[str, Any], *, event_id: str) -> Event:
    return Event(name, actor, subject, payload, int(time.time()), event_id)


__all__ = ["Event", "EventError", "make_event"]
