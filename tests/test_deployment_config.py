from __future__ import annotations

from pathlib import Path


def test_nginx_lab_config_has_proxy_hardening() -> None:
    config = Path("deploy/nginx/nginx.conf").read_text(encoding="utf-8")
    assert "proxy_pass http://angel_control_plane" in config
    assert "limit_req_zone" in config
    assert "X-Content-Type-Options nosniff" in config
    assert "server_tokens off" in config
