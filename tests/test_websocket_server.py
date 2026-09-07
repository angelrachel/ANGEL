from __future__ import annotations

import asyncio
import json

from websockets.asyncio.client import connect
from websockets.asyncio.server import serve

from src.c2.websocket_server import EventWebSocketServer
from src.event_stream import EventStream
from src.events import make_event


def test_websocket_event_transport_e2e() -> None:
    async def scenario() -> None:
        stream = EventStream()
        stream.publish(make_event("audit", "alice", "subject", {"ok": True}, event_id="evt-1"))
        transport = EventWebSocketServer(stream, "token")
        async with serve(transport.handler, "127.0.0.1", 0) as server:
            port = server.sockets[0].getsockname()[1]
            async with connect(
                f"ws://127.0.0.1:{port}", additional_headers={"Authorization": "Bearer token"}
            ) as client:
                await client.send(json.dumps({"cursor": 0}))
                response = json.loads(await client.recv())
                assert response["events"][0]["event_id"] == "evt-1"

    asyncio.run(scenario())
