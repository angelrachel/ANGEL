#!/usr/bin/env python3
from __future__ import annotations

import json
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
SCHEMA = json.loads((ROOT / "contracts/task/lab-agent.schema.json").read_text())

properties = SCHEMA["properties"]
assert properties["fixture_ref"]["pattern"] == "^fixture://"
assert properties["action"]["enum"] == ["observe", "simulate", "collect-evidence", "validate-control"]
assert "arbitrary-command" not in properties["action"]["enum"]
assert SCHEMA["additionalProperties"] is False
print("lab agent contract validated")
