from __future__ import annotations

import pytest

from src.config import ConfigError, Settings
from src.events import Event, EventError, make_event


def test_settings_defaults_and_explicit_values() -> None:
    settings = Settings.from_env(host="127.0.0.1", port=9000, database="test.db")
    assert settings.environment == "development"
    assert settings.port == 9000
    assert settings.database == "test.db"


def test_settings_reject_invalid_port(monkeypatch: pytest.MonkeyPatch) -> None:
    monkeypatch.setenv("ANGEL_PORT", "not-a-number")
    with pytest.raises(ConfigError, match="integer"):
        Settings.from_env()
    monkeypatch.setenv("ANGEL_PORT", "70000")
    with pytest.raises(ConfigError, match="between"):
        Settings.from_env()


def test_deployment_requires_real_secrets(monkeypatch: pytest.MonkeyPatch) -> None:
    monkeypatch.setenv("ANGEL_ENV", "production")
    with pytest.raises(ConfigError, match="explicitly"):
        Settings.from_env()
    monkeypatch.setenv("ANGEL_OPERATOR_KEY", "operator-key-that-is-long-enough")
    monkeypatch.setenv("ANGEL_SHARED_KEY", "shared-key-that-is-long-enough")
    monkeypatch.setenv("ANGEL_AUTH_SECRET", "auth-secret-that-is-long-enough")
    assert Settings.from_env().environment == "production"


def test_event_round_trip_and_validation() -> None:
    event = make_event("task.completed", "agent-1", "task-7", {"status": "passed"}, event_id="evt-1")
    assert Event.from_dict(event.to_dict()) == event
    with pytest.raises(EventError, match="unsupported"):
        Event("x", "a", "s", {}, 1, "id", version=2)
    with pytest.raises(EventError, match="too large"):
        Event("x", "a", "s", {"data": "x" * 20_000}, 1, "id")
