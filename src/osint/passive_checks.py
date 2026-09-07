"""Passive parsers for already-collected public metadata."""

from __future__ import annotations

import re
from dataclasses import dataclass


@dataclass(frozen=True)
class PassiveSignal:
    kind: str
    value: str
    source: str


def parse_asn_netblock(payload: dict[str, object]) -> list[PassiveSignal]:
    signals: list[PassiveSignal] = []
    for key in ("asn", "netblock", "prefix", "organization"):
        value = str(payload.get(key, "")).strip()
        if value:
            signals.append(PassiveSignal(key, value, "asn-lookup"))
    return signals


def parse_robots_txt(content: str) -> list[str]:
    paths: set[str] = set()
    for line in content.splitlines():
        match = re.match(r"^\s*(?:disallow|allow)\s*:\s*(\S+)", line, re.IGNORECASE)
        if match and match.group(1) != "/":
            paths.add(match.group(1))
    return sorted(paths)


def detect_cms_framework(headers: dict[str, str], body: str) -> list[PassiveSignal]:
    haystack = (" ".join(f"{key}:{value}" for key, value in headers.items()) + " " + body).lower()
    markers = {
        "wordpress": "wp-content",
        "drupal": "drupal-settings-json",
        "django": "csrfmiddlewaretoken",
        "nextjs": "__next_f.push",
        "rails": "x-runtime",
    }
    return [
        PassiveSignal("framework", name, "captured-response") for name, marker in markers.items() if marker in haystack
    ]


def detect_cloud_exposure(records: list[dict[str, object]]) -> list[PassiveSignal]:
    signals: set[PassiveSignal] = set()
    for record in records:
        host = str(record.get("hostname", "")).lower()
        provider = str(record.get("provider", "")).lower()
        if any(
            marker in host or marker in provider for marker in ("s3", "blob.core", "storage.googleapis", "cloudfront")
        ):
            signals.add(PassiveSignal("cloud", host or provider, "captured-metadata"))
    return sorted(signals, key=lambda signal: signal.value)


__all__ = ["PassiveSignal", "detect_cloud_exposure", "detect_cms_framework", "parse_asn_netblock", "parse_robots_txt"]
