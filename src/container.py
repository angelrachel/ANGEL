"""Minimal dependency injection container for service composition."""

from __future__ import annotations

from collections.abc import Callable
from typing import Any, TypeVar, cast

T = TypeVar("T")


class DependencyError(LookupError):
    """Raised when a dependency is missing or duplicated."""


class Container:
    def __init__(self) -> None:
        self._factories: dict[type[Any], Callable[[], Any]] = {}
        self._instances: dict[type[Any], Any] = {}

    def register_instance(self, interface: type[T], instance: T) -> None:
        if interface in self._factories or interface in self._instances:
            raise DependencyError("dependency already registered")
        self._instances[interface] = instance

    def register_factory(self, interface: type[T], factory: Callable[[], T]) -> None:
        if interface in self._factories or interface in self._instances:
            raise DependencyError("dependency already registered")
        self._factories[interface] = factory

    def resolve(self, interface: type[T]) -> T:
        if interface in self._instances:
            return cast(T, self._instances[interface])
        try:
            factory = self._factories[interface]
        except KeyError as exc:
            raise DependencyError("dependency is not registered") from exc
        instance = factory()
        self._instances[interface] = instance
        return cast(T, instance)


__all__ = ["Container", "DependencyError"]
