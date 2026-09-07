"""GraphQL schema and operation inventory without active network access."""

from __future__ import annotations

import re
from dataclasses import dataclass
from typing import Any, cast


@dataclass(frozen=True)
class GraphQLField:
    parent_type: str
    name: str
    arguments: tuple[str, ...]
    return_type: str
    requires_auth: bool


@dataclass(frozen=True)
class GraphQLOperation:
    operation: str
    name: str
    fields: tuple[str, ...]


@dataclass(frozen=True)
class GraphQLDrift:
    category: str
    subject: str
    details: str


def inventory_schema(document: dict[str, Any]) -> list[GraphQLField]:
    data = document.get("data", document)
    schema = data.get("__schema", {}) if isinstance(data, dict) else {}
    types = schema.get("types", []) if isinstance(schema, dict) else []
    fields: list[GraphQLField] = []
    for gql_type in types:
        if not isinstance(gql_type, dict) or gql_type.get("kind") not in {"OBJECT", "INTERFACE"}:
            continue
        parent = gql_type.get("name")
        if not isinstance(parent, str):
            continue
        for field in gql_type.get("fields", []) or []:
            if not isinstance(field, dict) or not isinstance(field.get("name"), str):
                continue
            arguments = tuple(
                sorted(
                    argument["name"]
                    for argument in field.get("args", [])
                    if isinstance(argument, dict) and isinstance(argument.get("name"), str)
                )
            )
            return_ref = field.get("type", {})
            return_type = _type_name(return_ref)
            fields.append(
                GraphQLField(parent, field["name"], arguments, return_type, bool(field.get("requiresAuth", False)))
            )
    return sorted(fields, key=lambda field: (field.parent_type, field.name))


def _type_name(type_ref: Any) -> str:
    if not isinstance(type_ref, dict):
        return ""
    if isinstance(type_ref.get("name"), str):
        return cast(str, type_ref["name"])
    return _type_name(type_ref.get("ofType"))


def inventory_operation(document: str) -> list[GraphQLOperation]:
    pattern = re.compile(r"\b(query|mutation|subscription)\s*([A-Za-z_][A-Za-z0-9_]*)?[^\{]*\{", re.DOTALL)
    operations: list[GraphQLOperation] = []
    for match in pattern.finditer(document):
        operation, name = match.groups()
        body = document[match.end() :]
        fields = tuple(sorted(set(re.findall(r"\b([A-Za-z_][A-Za-z0-9_]*)\s*(?=\(|\{)", body))))
        operations.append(GraphQLOperation(operation, name or "anonymous", fields))
    return operations


def compare_schema(expected: dict[str, Any], observed: dict[str, Any]) -> list[GraphQLDrift]:
    expected_fields = {(field.parent_type, field.name): field for field in inventory_schema(expected)}
    observed_fields = {(field.parent_type, field.name): field for field in inventory_schema(observed)}
    drift: list[GraphQLDrift] = []
    for key, field in observed_fields.items():
        if key not in expected_fields:
            drift.append(
                GraphQLDrift(
                    "undocumented_field",
                    f"{field.parent_type}.{field.name}",
                    "Observed field is absent from expected schema",
                )
            )
    for key, expected_field in expected_fields.items():
        observed_field = observed_fields.get(key)
        if observed_field is None:
            drift.append(
                GraphQLDrift(
                    "missing_field",
                    f"{expected_field.parent_type}.{expected_field.name}",
                    "Expected field is absent from observed schema",
                )
            )
            continue
        if expected_field.requires_auth != observed_field.requires_auth:
            drift.append(
                GraphQLDrift(
                    "authorization_drift",
                    f"{expected_field.parent_type}.{expected_field.name}",
                    "Authentication requirement differs between schemas",
                )
            )
        if expected_field.arguments != observed_field.arguments:
            drift.append(
                GraphQLDrift(
                    "argument_drift",
                    f"{expected_field.parent_type}.{expected_field.name}",
                    "Field argument set differs between schemas",
                )
            )
    return drift


__all__ = [
    "GraphQLDrift",
    "GraphQLField",
    "GraphQLOperation",
    "compare_schema",
    "inventory_operation",
    "inventory_schema",
]
