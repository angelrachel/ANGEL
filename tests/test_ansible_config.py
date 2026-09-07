from __future__ import annotations

from pathlib import Path


def test_ansible_deployment_is_non_root_and_health_checked() -> None:
    config = Path("deploy/ansible/angel.yml").read_text(encoding="utf-8")
    assert "user: \"10001:10001\"" in config
    assert "healthcheck:" in config
    assert "/var/lib/angel:/data" in config
