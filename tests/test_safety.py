from __future__ import annotations

import pytest

from src.safety import Capability, SafetyError, authorize_capability, reject_destructive_request


def test_safe_capabilities_are_explicitly_allowed() -> None:
    decision = authorize_capability(" inventory ")
    assert decision.capability is Capability.INVENTORY
    assert decision.allowed is True


def test_unknown_capability_fails_closed() -> None:
    with pytest.raises(SafetyError, match="not allowed"):
        authorize_capability("remote_exec")


@pytest.mark.parametrize("text", ["credential theft", "disable antivirus", "delete logs", "ransomware wiper"])
def test_destructive_requests_are_rejected(text: str) -> None:
    with pytest.raises(SafetyError, match="rejected"):
        reject_destructive_request(text)


def test_normal_defensive_request_is_allowed() -> None:
    reject_destructive_request("collect passive inventory and produce an evidence report")
