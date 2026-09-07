"""Small structured JSON logger used by service boundaries."""

from __future__ import annotations

import json
import logging
from typing import Any

SENSITIVE_FIELDS = frozenset({"password", "secret", "token", "api_key", "authorization"})


def configure_logging(level: int = logging.INFO) -> None:
    logging.basicConfig(level=level, format="%(message)s")


def log_event(logger: logging.Logger, event: str, *, correlation_id: str, **fields: Any) -> None:
    safe = {
        key: ("[REDACTED]" if key.lower() in SENSITIVE_FIELDS else value)
        for key, value in fields.items()
    }
    logger.info(json.dumps({"event": event, "correlation_id": correlation_id, **safe}, sort_keys=True, default=str))


__all__ = ["configure_logging", "log_event"]
