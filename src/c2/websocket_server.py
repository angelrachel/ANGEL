"""Authenticated WebSocket transport for the bounded event stream."""

from __future__ import annotations

import asyncio
import json
from typing import Any

from websockets.asyncio.server import ServerConnection, serve

from ..event_stream import EventStream


class WebSocketAuthError(Exception):
    """Raised when a WebSocket client is not authorized."""


class EventWebSocketServer:
    def __init__(self, stream: EventStream, bearer_token: str) -> None:
        if not bearer_token.strip():
            raise ValueError("bearer token is required")
        self.stream = stream
        self.bearer_token = bearer_token

    async def handler(self, websocket: ServerConnection) -> None:
        request = websocket.request
        if request is None or request.headers.get("Authorization") != f"Bearer {self.bearer_token}":
            await websocket.close(code=1008, reason="not authorized")
            return
        async for message in websocket:
            try:
                request = json.loads(message)
                if not isinstance(request, dict):
                    raise ValueError("request must be an object")
                cursor = int(request.get("cursor", 0))
                limit = int(request.get("limit", 100))
                items = self.stream.since(
                    cursor,
                    limit,
                    name=str(request["name"]) if request.get("name") else None,
                    actor=str(request["actor"]) if request.get("actor") else None,
                )
                payload: dict[str, Any] = {
                    "events": [{"cursor": item.cursor, **item.event.to_dict()} for item in items],
                    "latest_cursor": self.stream.latest_cursor,
                }
                await websocket.send(json.dumps(payload, sort_keys=True))
            except (TypeError, ValueError, json.JSONDecodeError) as exc:
                await websocket.send(json.dumps({"error": str(exc)}))

    async def run(self, host: str = "127.0.0.1", port: int = 8765) -> None:
        async with serve(self.handler, host, port):
            await asyncio.Future()


__all__ = ["EventWebSocketServer", "WebSocketAuthError"]
