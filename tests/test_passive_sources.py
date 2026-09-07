from __future__ import annotations

from src.osint.passive_sources import parse_certificate_transparency, parse_passive_dns


def test_passive_dns_is_scoped_and_deduplicated() -> None:
    records = [
        {"hostname": "api.example.test.", "value": "1.2.3.4"},
        {"hostname": "api.example.test", "value": "1.2.3.4"},
        {"hostname": "outside.test", "value": "5.6.7.8"},
    ]
    result = parse_passive_dns(records, {"example.test"})
    assert len(result) == 1 and result[0].hostname == "api.example.test"


def test_ct_parser_keeps_issuer_and_source() -> None:
    result = parse_certificate_transparency(
        [{"hostname": "www.example.test", "issuer": "Test CA", "first_seen": "2026-01-01"}],
        {"example.test"},
    )
    assert result[0].source == "certificate-transparency"
    assert result[0].value == "Test CA"
