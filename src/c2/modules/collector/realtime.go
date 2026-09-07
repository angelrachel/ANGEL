//go:build windows

package collector

import (
"os/exec"
"strings"
)

func GetSystemTime() string {
cmd := exec.Command("cmd", "/c", "time /t")
output, err := cmd.Output()
if err != nil {
return "Error: " + err.Error()
}
return strings.TrimSpace(string(output))
}

func GetCurrentDirectory() string {
cmd := exec.Command("cmd", "/c", "cd")
output, err := cmd.Output()
if err != nil {
return "Error: " + err.Error()
}
return strings.TrimSpace(string(output))
}

func GetEnvironmentVariables() string {
cmd := exec.Command("cmd", "/c", "set")
output, err := cmd.Output()
if err != nil {
return "Error: " + err.Error()
}
return strings.TrimSpace(string(output))
}

func GetRunningProcesses() string {
cmd := exec.Command("cmd", "/c", "tasklist /v")
output, err := cmd.Output()
if err != nil {
return "Error: " + err.Error()
}
return strings.TrimSpace(string(output))
}

func GetNetworkShares() string {
cmd := exec.Command("cmd", "/c", "net share")
output, err := cmd.Output()
if err != nil {
return "Error: " + err.Error()
}
return strings.TrimSpace(string(output))
}

func GetLoggedInUsers() string {
cmd := exec.Command("cmd", "/c", "query user")
output, err := cmd.Output()
if err != nil {
return "Error: " + err.Error()
}
return strings.TrimSpace(string(output))
}
