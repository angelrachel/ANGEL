"""Command-line entrypoint for the ANGEL control plane."""

from __future__ import annotations

import os

from .server import run_server

if __name__ == "__main__":
    run_server(int(os.environ.get("ANGEL_PORT", "8000")))
