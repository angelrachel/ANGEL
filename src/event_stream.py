"""Bounded in-process event stream for dashboard and audit consumers."""

from __future__ import annotations

import threading
from collections import deque
from dataclasses import dataclass

from .events import Event


@dataclass(frozen=True)
class StreamItem:
    cursor: int
    event: Event


class EventStream:
    def __init__(self, max_items: int = 1000) -> None:
        if max_items < 1:
            raise ValueError("max_items must be positive")
        self._items: deque[StreamItem] = deque(maxlen=max_items)
        self._next_cursor = 1
        self._lock = threading.Lock()

    def publish(self, event: Event) -> StreamItem:
        with self._lock:
            item = StreamItem(self._next_cursor, event)
            self._next_cursor += 1
            self._items.append(item)
            return item

    def since(self, cursor: int = 0, limit: int = 100) -> list[StreamItem]:
        if cursor < 0 or not 1 <= limit <= 500:
            raise ValueError("invalid stream cursor or limit")
        with self._lock:
            return [item for item in self._items if item.cursor > cursor][:limit]

    @property
    def latest_cursor(self) -> int:
        with self._lock:
            return self._next_cursor - 1


__all__ = ["EventStream", "StreamItem"]
