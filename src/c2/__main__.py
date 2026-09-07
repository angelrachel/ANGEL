import asyncio
import os

from .server import run_server

if __name__ == "__main__":
    try:
        asyncio.run(run_server(int(os.environ.get("ANGEL_PORT", "8000"))))
    except KeyboardInterrupt:
        print("\n[!] Server stopped.")
