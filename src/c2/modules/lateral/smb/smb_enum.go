//go:build windows

package smb

import (
"os/exec"
"strings"
)

func SMBShareEnum(targetHost, username, password string) string {
cmd := exec.Command("cmd", "/c", "net view \\\\"+targetHost)
output, err := cmd.Output()
if err != nil {
return "Error: " + err.Error()
}
return strings.TrimSpace(string(output))
}

func SMBUserEnum(targetHost string) string {
cmd := exec.Command("cmd", "/c", "net user /domain")
output, err := cmd.Output()
if err != nil {
return "Error: " + err.Error()
}
return strings.TrimSpace(string(output))
}

func SMBGroupEnum(targetHost string) string {
cmd := exec.Command("cmd", "/c", "net group /domain")
output, err := cmd.Output()
if err != nil {
return "Error: " + err.Error()
}
return strings.TrimSpace(string(output))
}

func SMBSystemInfo(targetHost string) string {
cmd := exec.Command("cmd", "/c", "systeminfo /s "+targetHost)
output, err := cmd.Output()
if err != nil {
return "Error: " + err.Error()
}
return strings.TrimSpace(string(output))
}
