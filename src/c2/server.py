"""Authenticated ANGEL control-plane HTTP server."""

from __future__ import annotations

import hashlib
import json
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer
from typing import Any, ClassVar, cast
from urllib.parse import parse_qs, urlparse

from ..auth import AuthError, authenticate
from ..config import Settings
from ..evidence.chain import redact
from ..reporting import Report
from .crypto import CryptoError, ReplayGuard, SessionCipher
from .policy import validate_task
from .storage import Store


def openapi_document() -> dict[str, Any]:
    return {
        "openapi": "3.0.3",
        "info": {"title": "ANGEL Control Plane", "version": "1.0.0"},
        "paths": {
            "/healthz": {"get": {"responses": {"200": {"description": "Service is alive"}}}},
            "/readyz": {"get": {"responses": {"200": {"description": "Service is ready"}}}},
            "/openapi.json": {"get": {"responses": {"200": {"description": "OpenAPI document"}}}},
            "/capabilities": {
                "get": {"security": [{"operatorKey": []}], "responses": {"200": {"description": "Capabilities"}}}
            },
            "/agents": {
                "get": {"security": [{"operatorKey": []}], "responses": {"200": {"description": "Registered agents"}}}
            },
            "/tasks": {
                "get": {"security": [{"operatorKey": []}], "responses": {"200": {"description": "Tasks"}}}
            },
            "/audit": {
                "get": {"security": [{"bearerAuth": []}], "responses": {"200": {"description": "Audit events"}}}
            },
            "/register": {
                "post": {"security": [{"operatorKey": []}], "responses": {"200": {"description": "Registered"}}}
            },
            "/heartbeat": {
                "post": {"security": [{"operatorKey": []}], "responses": {"200": {"description": "Acknowledged"}}}
            },
            "/task/queue": {
                "post": {"security": [{"operatorKey": []}], "responses": {"202": {"description": "Queued"}}}
            },
            "/task/claim": {
                "post": {"security": [{"operatorKey": []}], "responses": {"200": {"description": "Claimed task"}}}
            },
            "/result": {
                "post": {"security": [{"operatorKey": []}], "responses": {"201": {"description": "Recorded"}}}
            },
            "/scope": {
                "post": {"security": [{"bearerAuth": []}], "responses": {"201": {"description": "Scope created"}}}
            },
            "/evidence": {
                "post": {"security": [{"bearerAuth": []}], "responses": {"201": {"description": "Evidence stored"}}}
            },
            "/report": {
                "post": {"security": [{"bearerAuth": []}], "responses": {"201": {"description": "Report stored"}}}
            },
        },
        "components": {
            "securitySchemes": {
                "operatorKey": {"type": "apiKey", "in": "header", "name": "X-ANGEL-Key"},
                "bearerAuth": {"type": "http", "scheme": "bearer"},
            }
        },
    }


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
        parsed = urlparse(self.path)
        if parsed.path == "/healthz":
            self._json(200, {"status": "ok"})
            return
        if parsed.path == "/readyz":
            try:
                self.store.get_agent("")
            except Exception:
                self._json(503, {"status": "not_ready"})
                return
            self._json(200, {"status": "ready"})
            return
        if parsed.path == "/openapi.json":
            self._json(200, openapi_document())
            return
        if not self._authorized():
            if parsed.path != "/audit" or not self.headers.get("Authorization", "").startswith("Bearer "):
                self._json(401, {"error": "unauthorized"})
                return
            try:
                if not self._principal().can("read"):
                    raise AuthError("read permission required")
            except AuthError as exc:
                self._json(401, {"error": str(exc)})
                return
        if parsed.path == "/capabilities":
            self._json(200, {"service": "angel-c2", "protocol": 1})
            return
        if parsed.path == "/agents":
            self._json(200, {"agents": [agent.__dict__ for agent in self.store.list_agents()]})
            return
        if parsed.path == "/tasks":
            agent_values = parse_qs(parsed.query).get("agent_id", [])
            agent_id = agent_values[0] if agent_values else None
            self._json(200, {"tasks": [task.__dict__ for task in self.store.list_tasks(agent_id)]})
            return
        if parsed.path == "/audit":
            values = parse_qs(parsed.query).get("limit", ["100"])
            try:
                limit = int(values[0])
                events = self.store.list_audit_events(limit)
            except ValueError as exc:
                self._json(400, {"error": str(exc)})
                return
            self._json(200, {"events": events})
            return
        self._json(404, {"error": "not found"})

    def do_POST(self) -> None:  # noqa: N802
        parsed_path = urlparse(self.path).path
        bearer_paths = {"/scope", "/evidence", "/report"}
        has_bearer = self.headers.get("Authorization", "").startswith("Bearer ")
        if (
            not self._authorized()
            and parsed_path not in {"/register", "/heartbeat", "/healthz"}
            and not (parsed_path in bearer_paths and has_bearer)
        ):
            self._json(401, {"error": "unauthorized"})
            return
        if parsed_path in {"/register", "/heartbeat"} and not self._authorized():
            self._json(401, {"error": "operator key required"})
            return
        try:
            data = self._body()
            if parsed_path in {"/register", "/heartbeat"}:
                required = ("agent_id", "hostname", "os", "arch")
                if any(not isinstance(data.get(key), str) or not data[key] for key in required):
                    raise ValueError("missing agent registration field")
                agent = self.store.upsert_agent(data["agent_id"], data["hostname"], data["os"], data["arch"])
                status = "registered" if parsed_path == "/register" else "heartbeat_ack"
                self._json(200, {"status": status, "agent_id": agent.id, "last_seen": agent.last_seen})
                return
            if parsed_path == "/task/queue":
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
            if parsed_path == "/task/claim":
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
            if parsed_path == "/result":
                task_id, agent_id, payload = data.get("task_id"), data.get("agent_id"), data.get("payload")
                if not isinstance(task_id, int) or not isinstance(agent_id, str) or not isinstance(payload, dict):
                    raise ValueError("invalid result fields")
                self.store.record_result(task_id, agent_id, payload)
                self._json(201, {"status": "recorded"})
                return
            if parsed_path == "/scope":
                principal = self._principal()
                if not principal.can("write"):
                    raise AuthError("write permission required")
                scope = self.store.create_scope(
                    data["name"], data["hosts"], data.get("paths", ["/"]), data.get("expires_at")
                )
                self._json(201, {"id": scope.id, "name": scope.name, "status": scope.status})
                return
            if parsed_path == "/evidence":
                principal = self._principal()
                if not principal.can("evidence"):
                    raise AuthError("evidence permission required")
                payload = json.loads(redact(json.dumps(data.get("payload", {}))))
                evidence = self.store.add_evidence(
                    data["scope_id"], data["evidence_type"], principal.subject, payload, data["record_hash"]
                )
                self._json(201, {"id": evidence.id, "record_hash": evidence.record_hash})
                return
            if parsed_path == "/report":
                principal = self._principal()
                if not principal.can("report"):
                    raise AuthError("report permission required")
                report = Report(data["title"], data.get("engagement", ""))
                record = self.store.add_report(data["scope_id"], report.title, report.to_dict())
                self._json(201, {"id": record.id, "title": record.title})
                return
            self._json(404, {"error": "not found"})
        except AuthError as exc:
            error_status = 403 if "permission" in str(exc) else 401
            self._json(error_status, {"error": str(exc)})
        except (ValueError, json.JSONDecodeError, KeyError) as exc:
            self._json(400, {"error": str(exc)})
        except Exception:
            self._json(500, {"error": "internal error"})

    def log_message(self, _format: str, *_args: object) -> None:
        return


def build_server(host: str = "127.0.0.1", port: int = 8000, database: str = "c2.db") -> ThreadingHTTPServer:
    settings = Settings.from_env(host=host, port=port, database=database)
    key = settings.operator_key
    material = hashlib.sha512(settings.shared_key.encode()).digest()
    handler = cast(type[C2Handler], type("ConfiguredC2Handler", (C2Handler,), {}))
    handler.store = Store(settings.database)
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
