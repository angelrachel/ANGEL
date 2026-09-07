import platform
import socket


def run_windows_implant():
    hostname = socket.gethostname()
    os_info = platform.system() + " " + platform.release()
    print(f"[+] Windows implant started on {hostname}")
    print(f"[+] OS: {os_info}")
    print("[+] Implant is ready for tasking")
