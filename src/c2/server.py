"""Authenticated ANGEL control-plane HTTP server."""

from __future__ import annotations

import hashlib
import json
import os
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer
from typing import Any, ClassVar, cast

from ..auth import AuthError, authenticate
from ..evidence.chain import redact
from ..reporting import Report
from .crypto import CryptoError, ReplayGuard, SessionCipher
from .policy import validate_task
from .storage import Store


class C2Handler(BaseHTTPRequestHandler):
    store: ClassVar[Store]
    cipher: ClassVar[SessionCipher]
    replay: ClassVar[ReplayGuard]
    operator_key: ClassVar[str]

    def _json(self, status: int, payload: dict[str, Any]) -> None:
        body = json.dumps(payload, sort_keys=True).encode()
        self.send_response(status)
        self.send_header("Content-Type", "application/json")
        self.send_header("Content-Length", str(len(body)))
        self.end_headers()
        self.wfile.write(body)

    def _authorized(self) -> bool:
        supplied = self.headers.get("X-ANGEL-Key", "")
        return bool(supplied) and supplied == self.operator_key

    def _principal(self) -> Any:
        header = self.headers.get("Authorization", "")
        if not header.startswith("Bearer "):
            raise AuthError("bearer token required")
        return authenticate(header[7:])

    def _body(self) -> dict[str, Any]:
        length = int(self.headers.get("Content-Length", "0"))
        if length <= 0 or length > 1_048_576:
            raise ValueError("invalid request size")
        parsed = json.loads(self.rfile.read(length))
        if not isinstance(parsed, dict):
            raise ValueError("request must be a JSON object")
        return parsed

    def do_GET(self) -> None:  # noqa: N802
        if self.path == "/healthz":
            self._json(200, {"status": "ok"})
            return
        if not self._authorized():
            self._json(401, {"error": "unauthorized"})
            return
        if self.path == "/capabilities":
            self._json(200, {"service": "angel-c2", "protocol": 1})
            return
        self._json(404, {"error": "not found"})

    def do_POST(self) -> None:  # noqa: N802
        bearer_paths = {"/scope", "/evidence", "/report"}
        has_bearer = self.headers.get("Authorization", "").startswith("Bearer ")
        if (
            not self._authorized()
            and self.path not in {"/register", "/healthz"}
            and not (self.path in bearer_paths and has_bearer)
        ):
            self._json(401, {"error": "unauthorized"})
            return
        try:
            data = self._body()
            if self.path == "/register":
                required = ("agent_id", "hostname", "os", "arch")
                if any(not isinstance(data.get(key), str) or not data[key] for key in required):
                    raise ValueError("missing agent registration field")
                agent = self.store.upsert_agent(data["agent_id"], data["hostname"], data["os"], data["arch"])
                self._json(200, {"status": "registered", "agent_id": agent.id})
                return
            if self.path == "/task/queue":
                agent_id = data.get("agent_id")
                task_type = data.get("task_type")
                payload = data.get("payload", {})
                if not isinstance(agent_id, str) or not isinstance(task_type, str):
                    raise ValueError("invalid task fields")
                if self.store.get_agent(agent_id) is None:
                    self._json(404, {"error": "agent not found"})
                    return
                validate_task(task_type, payload)
                queued_task = self.store.enqueue_task(agent_id, task_type, payload)
                self._json(202, {"task_id": queued_task.id, "status": queued_task.status})
                return
            if self.path == "/task/claim":
                agent_id = data.get("agent_id")
                if not isinstance(agent_id, str) or self.store.get_agent(agent_id) is None:
                    self._json(404, {"error": "agent not found"})
                    return
                claimed_task = self.store.claim_task(agent_id)
                self._json(
                    200,
                    {
                        "task": None
                        if claimed_task is None
                        else {
                            "id": claimed_task.id,
                            "type": claimed_task.task_type,
                            "payload": claimed_task.payload,
                        }
                    },
                )
                return
            if self.path == "/result":
                task_id, agent_id, payload = data.get("task_id"), data.get("agent_id"), data.get("payload")
                if not isinstance(task_id, int) or not isinstance(agent_id, str) or not isinstance(payload, dict):
                    raise ValueError("invalid result fields")
                self.store.record_result(task_id, agent_id, payload)
                self._json(201, {"status": "recorded"})
                return
            if self.path == "/scope":
                principal = self._principal()
                if not principal.can("write"):
                    raise AuthError("write permission required")
                scope = self.store.create_scope(
                    data["name"], data["hosts"], data.get("paths", ["/"]), data.get("expires_at")
                )
                self._json(201, {"id": scope.id, "name": scope.name, "status": scope.status})
                return
            if self.path == "/evidence":
                principal = self._principal()
                if not principal.can("evidence"):
                    raise AuthError("evidence permission required")
                payload = json.loads(redact(json.dumps(data.get("payload", {}))))
                evidence = self.store.add_evidence(
                    data["scope_id"], data["evidence_type"], principal.subject, payload, data["record_hash"]
                )
                self._json(201, {"id": evidence.id, "record_hash": evidence.record_hash})
                return
            if self.path == "/report":
                principal = self._principal()
                if not principal.can("report"):
                    raise AuthError("report permission required")
                report = Report(data["title"], data.get("engagement", ""))
                record = self.store.add_report(data["scope_id"], report.title, report.to_dict())
                self._json(201, {"id": record.id, "title": record.title})
                return
            self._json(404, {"error": "not found"})
        except (ValueError, json.JSONDecodeError, AuthError, KeyError) as exc:
            self._json(400, {"error": str(exc)})
        except Exception:
            self._json(500, {"error": "internal error"})

    def log_message(self, _format: str, *_args: object) -> None:
        return


def build_server(host: str = "127.0.0.1", port: int = 8000, database: str = "c2.db") -> ThreadingHTTPServer:
    Store(database)
    key = os.environ.get("ANGEL_OPERATOR_KEY", "development-operator-key")
    material = hashlib.sha512(os.environ.get("ANGEL_SHARED_KEY", "development-only-key").encode()).digest()
    handler = cast(type[C2Handler], type("ConfiguredC2Handler", (C2Handler,), {}))
    handler.store = Store(database)
    handler.cipher = SessionCipher(material)
    handler.replay = ReplayGuard()
    handler.operator_key = key
    return ThreadingHTTPServer((host, port), handler)


def run_server(port: int = 8000) -> None:
    server = build_server(port=port)
    print(f"ANGEL control plane listening on 127.0.0.1:{port}")
    try:
        server.serve_forever()
    except KeyboardInterrupt:
        pass
    finally:
        server.server_close()


__all__ = ["C2Handler", "build_server", "run_server", "CryptoError"]
