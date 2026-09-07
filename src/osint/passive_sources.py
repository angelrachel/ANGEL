"""Parse passive DNS and CT records supplied by an external source."""

from __future__ import annotations

from dataclasses import dataclass


@dataclass(frozen=True)
class PassiveRecord:
    hostname: str
    source: str
    value: str
    first_seen: str | None = None


def parse_passive_dns(records: list[dict[str, object]], allowed_hosts: set[str]) -> list[PassiveRecord]:
    return _parse(records, allowed_hosts, "passive-dns", "value")


def parse_certificate_transparency(records: list[dict[str, object]], allowed_hosts: set[str]) -> list[PassiveRecord]:
    return _parse(records, allowed_hosts, "certificate-transparency", "issuer")


def _parse(
    records: list[dict[str, object]],
    allowed_hosts: set[str],
    source: str,
    value_key: str,
) -> list[PassiveRecord]:
    normalized = {host.lower().rstrip(".") for host in allowed_hosts if host.strip()}
    output: set[PassiveRecord] = set()
    for record in records:
        hostname = str(record.get("hostname", "")).lower().rstrip(".")
        value = str(record.get(value_key, ""))
        if not hostname or not value:
            continue
        if not any(hostname == allowed or hostname.endswith("." + allowed) for allowed in normalized):
            continue
        seen = record.get("first_seen")
        output.add(PassiveRecord(hostname, source, value, str(seen) if seen is not None else None))
    return sorted(output, key=lambda item: (item.hostname, item.value))


__all__ = ["PassiveRecord", "parse_certificate_transparency", "parse_passive_dns"]
