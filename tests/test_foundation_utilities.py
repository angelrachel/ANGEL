from __future__ import annotations

import json
import logging

import pytest

from src.cancellation import CancellationContext, CancellationError
from src.errors import AuthorizationError, ValidationError
from src.structured_logging import log_event


def test_error_model_serializes_stable_contract() -> None:
    error = ValidationError("bad input")
    assert error.status == 400
    assert error.to_dict()["error"] == "validation_error"
    assert AuthorizationError().status == 403


def test_structured_logger_redacts_secrets(caplog: pytest.LogCaptureFixture) -> None:
    logger = logging.getLogger("angel-test")
    with caplog.at_level(logging.INFO):
        log_event(logger, "request", correlation_id="corr-1", token="secret-value", path="/healthz")
    payload = json.loads(caplog.records[-1].message)
    assert payload["token"] == "[REDACTED]"
    assert payload["correlation_id"] == "corr-1"


def test_cancellation_context() -> None:
    context = CancellationContext()
    assert not context.cancelled
    context.cancel()
    with pytest.raises(CancellationError):
        context.raise_if_cancelled()
    assert context.wait(0) is True
