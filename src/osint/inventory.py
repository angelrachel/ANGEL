"""Passive, scope-aware asset inventory helpers."""

from __future__ import annotations

import hashlib
from dataclasses import dataclass
from datetime import UTC, datetime
from urllib.parse import urlparse

from cryptography import x509
from cryptography.hazmat.primitives import hashes, serialization


class InventoryError(ValueError):
    """Raised when an inventory input is malformed."""


@dataclass(frozen=True)
class Asset:
    host: str
    scheme: str
    port: int
    source: str
    fingerprint: str


def inventory_urls(urls: list[str], allowed_hosts: set[str]) -> list[Asset]:
    normalized_hosts = {host.lower().rstrip(".") for host in allowed_hosts if host.strip()}
    if not normalized_hosts:
        raise InventoryError("at least one allowed host is required")
    assets: dict[tuple[str, int], Asset] = {}
    for raw_url in urls:
        parsed = urlparse(raw_url)
        host = (parsed.hostname or "").lower().rstrip(".")
        if parsed.scheme not in {"http", "https"} or not host:
            raise InventoryError("only HTTP(S) URLs are supported")
        if not any(host == allowed or host.endswith("." + allowed) for allowed in normalized_hosts):
            continue
        port = parsed.port or (443 if parsed.scheme == "https" else 80)
        key = (host, port)
        fingerprint = hashlib.sha256(f"{parsed.scheme}://{host}:{port}".encode()).hexdigest()
        assets[key] = Asset(host, parsed.scheme, port, "passive-url", fingerprint)
    return sorted(assets.values(), key=lambda asset: (asset.host, asset.port))


def parse_certificate(certificate: bytes) -> dict[str, object]:
    try:
        cert = x509.load_pem_x509_certificate(certificate)
    except ValueError:
        try:
            cert = x509.load_der_x509_certificate(certificate)
        except ValueError as exc:
            raise InventoryError("invalid X.509 certificate") from exc
    names: list[str] = []
    try:
        extension = cert.extensions.get_extension_for_class(x509.SubjectAlternativeName)
        names = sorted(extension.value.get_values_for_type(x509.DNSName))
    except x509.ExtensionNotFound:
        pass
    return {
        "subject": cert.subject.rfc4514_string(),
        "issuer": cert.issuer.rfc4514_string(),
        "serial": str(cert.serial_number),
        "not_before": _iso(cert.not_valid_before_utc),
        "not_after": _iso(cert.not_valid_after_utc),
        "dns_names": names,
        "sha256": cert.fingerprint(hashes.SHA256()).hex(),
        "public_key_type": type(cert.public_key()).__name__,
        "der_size": len(cert.public_bytes(serialization.Encoding.DER)),
    }


def _iso(value: datetime) -> str:
    return value.astimezone(UTC).isoformat()


__all__ = ["Asset", "InventoryError", "inventory_urls", "parse_certificate"]


def fingerprint_response(headers: dict[str, str], body: str = "") -> list[str]:
    """Infer coarse technology hints from already-captured response data."""
    normalized = {key.lower(): value.lower() for key, value in headers.items()}
    haystack = body.lower()
    findings: set[str] = set()
    server = normalized.get("server", "")
    powered = normalized.get("x-powered-by", "")
    if "nginx" in server:
        findings.add("nginx")
    if "apache" in server:
        findings.add("apache")
    if "express" in powered or "express" in haystack:
        findings.add("express")
    if "wordpress" in haystack or "wp-content" in haystack:
        findings.add("wordpress")
    if "graphql" in haystack or "graphql" in normalized.get("content-type", ""):
        findings.add("graphql")
    return sorted(findings)


__all__ = ["Asset", "InventoryError", "fingerprint_response", "inventory_urls", "parse_certificate"]
