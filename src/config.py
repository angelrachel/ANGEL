"""Validated runtime configuration for the ANGEL control plane."""

from __future__ import annotations

import os
from dataclasses import dataclass


class ConfigError(ValueError):
    """Raised when runtime configuration is missing or invalid."""


_PLACEHOLDERS = {
    "development-auth-secret",
    "development-operator-key",
    "development-only-key",
    "change-me",
    "your-secret-key-here",
}


@dataclass(frozen=True)
class Settings:
    environment: str
    host: str
    port: int
    database: str
    operator_key: str
    shared_key: str
    auth_secret: str

    @classmethod
    def from_env(
        cls,
        *,
        host: str | None = None,
        port: int | None = None,
        database: str | None = None,
    ) -> Settings:
        environment = os.environ.get("ANGEL_ENV", "development").strip().lower()
        if environment not in {"development", "test", "staging", "production"}:
            raise ConfigError("ANGEL_ENV must be development, test, staging, or production")
        selected_port = port if port is not None else _int_env("ANGEL_PORT", 8000)
        if selected_port != 0 and not 1 <= selected_port <= 65535:
            raise ConfigError("ANGEL_PORT must be between 1 and 65535")
        selected_host = host if host is not None else os.environ.get("ANGEL_HOST", "127.0.0.1")
        selected_database = database if database is not None else os.environ.get("ANGEL_DATABASE", "c2.db")
        if not selected_host.strip() or not selected_database.strip():
            raise ConfigError("ANGEL_HOST and ANGEL_DATABASE must not be empty")
        settings = cls(
            environment,
            selected_host,
            selected_port,
            selected_database,
            os.environ.get("ANGEL_OPERATOR_KEY", "development-operator-key"),
            os.environ.get("ANGEL_SHARED_KEY", "development-only-key"),
            os.environ.get("ANGEL_AUTH_SECRET", "development-auth-secret"),
        )
        if environment in {"staging", "production"}:
            settings._validate_deployment_secrets()
        return settings

    def _validate_deployment_secrets(self) -> None:
        values = (self.operator_key, self.shared_key, self.auth_secret)
        if any(not value.strip() or value in _PLACEHOLDERS for value in values):
            raise ConfigError("deployment secrets must be explicitly configured")
        if len(self.operator_key) < 16 or len(self.shared_key) < 16 or len(self.auth_secret) < 16:
            raise ConfigError("deployment secrets must be at least 16 characters")


def _int_env(name: str, default: int) -> int:
    raw = os.environ.get(name)
    if raw is None:
        return default
    try:
        return int(raw)
    except ValueError as exc:
        raise ConfigError(f"{name} must be an integer") from exc


__all__ = ["ConfigError", "Settings"]
