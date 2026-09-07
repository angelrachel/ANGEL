from __future__ import annotations

from pathlib import Path


def test_terraform_lab_module_is_bounded() -> None:
    config = Path("deploy/terraform/main.tf").read_text(encoding="utf-8")
    assert 'name = "angel-lab"' in config
    assert "internal = 8000" in config
    assert 'restart = "unless-stopped"' in config
    assert "docker_network" in config
    assert 'ip       = "127.0.0.1"' in config
    assert "read_only     = true" in config
    assert "no-new-privileges:true" in config
