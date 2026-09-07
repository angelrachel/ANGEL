"""Stable application error model for API and orchestrator boundaries."""

from __future__ import annotations

from dataclasses import dataclass


@dataclass(frozen=True)
class AppError(Exception):
    code: str
    message: str
    status: int = 400
    retryable: bool = False

    def __str__(self) -> str:
        return self.message

    def to_dict(self) -> dict[str, object]:
        return {"error": self.code, "message": self.message, "retryable": self.retryable}


class ValidationError(AppError):
    def __init__(self, message: str) -> None:
        super().__init__("validation_error", message, 400, False)


class AuthorizationError(AppError):
    def __init__(self, message: str = "not authorized") -> None:
        super().__init__("authorization_error", message, 403, False)


class ConflictError(AppError):
    def __init__(self, message: str) -> None:
        super().__init__("conflict", message, 409, False)


__all__ = ["AppError", "AuthorizationError", "ConflictError", "ValidationError"]
