from __future__ import annotations

from src.event_stream import EventStream
from src.events import make_event


def test_event_stream_cursor_and_retention() -> None:
    stream = EventStream(max_items=2)
    stream.publish(make_event("one", "a", "s", {}, event_id="1"))
    stream.publish(make_event("two", "a", "s", {}, event_id="2"))
    stream.publish(make_event("three", "a", "s", {}, event_id="3"))
    assert stream.latest_cursor == 3
    assert [item.cursor for item in stream.since(0)] == [2, 3]
    assert [item.event.name for item in stream.since(2)] == ["three"]
    assert [item.event.name for item in stream.since(0, name="two")] == ["two"]
    assert stream.since(0, actor="missing") == []
