//go:build windows

package collector

import (
"os/exec"
"strings"
)

func CollectSystemInfo() string {
cmd := exec.Command("cmd", "/c", "systeminfo")
output, err := cmd.Output()
if err != nil {
return ""
}
return strings.TrimSpace(string(output))
}

func CollectUsers() string {
cmd := exec.Command("cmd", "/c", "net user")
output, err := cmd.Output()
if err != nil {
return ""
}
return strings.TrimSpace(string(output))
}

func CollectRunningProcesses() string {
cmd := exec.Command("cmd", "/c", "tasklist /v")
output, err := cmd.Output()
if err != nil {
return ""
}
return strings.TrimSpace(string(output))
}

func CollectServices() string {
cmd := exec.Command("cmd", "/c", "wmic service list brief")
output, err := cmd.Output()
if err != nil {
return ""
}
return strings.TrimSpace(string(output))
}

func CollectAll() string {
var result strings.Builder
result.WriteString(CollectSystemInfo())
result.WriteString(CollectUsers())
result.WriteString(CollectRunningProcesses())
result.WriteString(CollectServices())
return result.String()
}
