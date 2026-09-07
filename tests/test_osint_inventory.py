from __future__ import annotations

from datetime import UTC, datetime, timedelta

import pytest
from cryptography import x509
from cryptography.hazmat.primitives import hashes, serialization
from cryptography.hazmat.primitives.asymmetric import ec
from cryptography.x509.oid import NameOID

from src.osint.inventory import InventoryError, inventory_urls, parse_certificate


def test_inventory_is_scoped_and_deduplicated() -> None:
    assets = inventory_urls(
        ["https://api.example.test/v1", "https://api.example.test/other", "http://outside.test"],
        {"example.test"},
    )
    assert len(assets) == 1
    assert assets[0].host == "api.example.test"
    with pytest.raises(InventoryError):
        inventory_urls(["ftp://example.test"], {"example.test"})


def test_certificate_metadata_parser() -> None:
    key = ec.generate_private_key(ec.SECP256R1())
    name = x509.Name([x509.NameAttribute(NameOID.COMMON_NAME, "example.test")])
    cert = (
        x509.CertificateBuilder()
        .subject_name(name)
        .issuer_name(name)
        .public_key(key.public_key())
        .serial_number(x509.random_serial_number())
        .not_valid_before(datetime.now(UTC) - timedelta(minutes=1))
        .not_valid_after(datetime.now(UTC) + timedelta(days=1))
        .add_extension(x509.SubjectAlternativeName([x509.DNSName("example.test")]), critical=False)
        .sign(key, hashes.SHA256())
    )
    parsed = parse_certificate(cert.public_bytes(serialization.Encoding.PEM))
    assert parsed["dns_names"] == ["example.test"]
    assert len(str(parsed["sha256"])) == 64
