from __future__ import annotations

from src.osint.passive_checks import detect_cloud_exposure, detect_cms_framework, parse_asn_netblock, parse_robots_txt


def test_passive_metadata_parsers() -> None:
    assert parse_asn_netblock({"asn": "AS64500", "netblock": "192.0.2.0/24"})[0].value == "AS64500"
    assert parse_robots_txt("User-agent: *\nDisallow: /admin\nAllow: /public") == ["/admin", "/public"]
    assert detect_cms_framework({"Content-Type": "text/html"}, "<div id='__next_f.push'>")[0].value == "nextjs"
    assert detect_cloud_exposure([{"hostname": "bucket.s3.amazonaws.com"}])[0].kind == "cloud"
