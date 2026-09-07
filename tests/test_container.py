from __future__ import annotations

import pytest

from src.container import Container, DependencyError


class Service:
    pass


def test_container_resolves_singleton_factory() -> None:
    container = Container()
    container.register_factory(Service, Service)
    assert container.resolve(Service) is container.resolve(Service)
    with pytest.raises(DependencyError, match="already"):
        container.register_instance(Service, Service())


def test_container_rejects_missing_dependency() -> None:
    with pytest.raises(DependencyError, match="not registered"):
        Container().resolve(Service)
