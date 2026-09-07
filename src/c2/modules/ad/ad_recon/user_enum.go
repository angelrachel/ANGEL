//go:build windows

package ad_recon

import (
"os/exec"
"strings"
)

func EnumerateUser(username string) string {
cmd := exec.Command("cmd", "/c", "net user "+username+" /domain")
output, err := cmd.Output()
if err != nil {
return "Error: " + err.Error()
}
return strings.TrimSpace(string(output))
}

func EnumerateUserSPN() string {
cmd := exec.Command("cmd", "/c", "setspn -Q */*")
output, err := cmd.Output()
if err != nil {
return "Error: " + err.Error()
}
return strings.TrimSpace(string(output))
}

func EnumerateUserSessions() string {
cmd := exec.Command("cmd", "/c", "qwinsta")
output, err := cmd.Output()
if err != nil {
return "Error: " + err.Error()
}
return strings.TrimSpace(string(output))
}

func EnumerateUserPrivileges() string {
cmd := exec.Command("cmd", "/c", "whoami /priv")
output, err := cmd.Output()
if err != nil {
return "Error: " + err.Error()
}
return strings.TrimSpace(string(output))
}

func EnumerateUserGroups() string {
cmd := exec.Command("cmd", "/c", "whoami /groups")
output, err := cmd.Output()
if err != nil {
return "Error: " + err.Error()
}
return strings.TrimSpace(string(output))
}
