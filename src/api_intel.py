"""OpenAPI inventory and schema drift analysis."""

from __future__ import annotations

from dataclasses import dataclass
from typing import Any


@dataclass(frozen=True)
class Endpoint:
    path: str
    method: str
    operation_id: str
    parameters: tuple[str, ...]
    security: tuple[str, ...]


@dataclass(frozen=True)
class SchemaDrift:
    path: str
    method: str
    issue: str
    details: str


def inventory_openapi(document: dict[str, Any]) -> list[Endpoint]:
    paths = document.get("paths", {})
    if not isinstance(paths, dict):
        raise ValueError("OpenAPI paths must be an object")
    endpoints: list[Endpoint] = []
    for path, item in paths.items():
        if not isinstance(path, str) or not isinstance(item, dict):
            continue
        for method, operation in item.items():
            if method.lower() not in {"get", "post", "put", "patch", "delete", "head", "options"}:
                continue
            if not isinstance(operation, dict):
                continue
            parameters = operation.get("parameters", [])
            names = tuple(
                sorted(
                    parameter["name"]
                    for parameter in parameters
                    if isinstance(parameter, dict) and isinstance(parameter.get("name"), str)
                )
            )
            security = operation.get("security", document.get("security", []))
            security_names = tuple(
                sorted(
                    name
                    for entry in security
                    if isinstance(entry, dict)
                    for name in entry
                    if isinstance(name, str)
                )
            )
            endpoints.append(
                Endpoint(
                    path,
                    method.upper(),
                    str(operation.get("operationId", "")),
                    names,
                    security_names,
                )
            )
    return sorted(endpoints, key=lambda endpoint: (endpoint.path, endpoint.method))


def compare_openapi(expected: dict[str, Any], observed: dict[str, Any]) -> list[SchemaDrift]:
    expected_endpoints = {(item.path, item.method): item for item in inventory_openapi(expected)}
    observed_endpoints = {(item.path, item.method): item for item in inventory_openapi(observed)}
    drift: list[SchemaDrift] = []
    for key, endpoint in observed_endpoints.items():
        if key not in expected_endpoints:
            drift.append(
                SchemaDrift(
                    endpoint.path,
                    endpoint.method,
                    "undocumented_endpoint",
                    "Observed endpoint is absent from the expected schema",
                )
            )
    for key, endpoint in expected_endpoints.items():
        if key not in observed_endpoints:
            drift.append(
                SchemaDrift(
                    endpoint.path, endpoint.method, "unavailable_endpoint", "Expected endpoint was not observed"
                )
            )
            continue
        observed_endpoint = observed_endpoints[key]
        missing = sorted(set(endpoint.parameters) - set(observed_endpoint.parameters))
        if missing:
            drift.append(
                SchemaDrift(
                    endpoint.path,
                    endpoint.method,
                    "parameter_drift",
                    f"Observed schema is missing: {', '.join(missing)}",
                )
            )
        if endpoint.security != observed_endpoint.security:
            drift.append(
                SchemaDrift(
                    endpoint.path,
                    endpoint.method,
                    "security_drift",
                    "Observed security requirements differ from expected schema",
                )
            )
    return drift


__all__ = ["Endpoint", "SchemaDrift", "compare_openapi", "inventory_openapi"]
