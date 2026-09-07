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
return "Error: " + err.Error()
}
return strings.TrimSpace(string(output))
}

func CollectUsers() string {
cmd := exec.Command("cmd", "/c", "net user")
output, err := cmd.Output()
if err != nil {
return "Error: " + err.Error()
}
return strings.TrimSpace(string(output))
}

func CollectGroups() string {
cmd := exec.Command("cmd", "/c", "net localgroup")
output, err := cmd.Output()
if err != nil {
return "Error: " + err.Error()
}
return strings.TrimSpace(string(output))
}

func CollectNetworkInfo() string {
cmd := exec.Command("cmd", "/c", "ipconfig /all")
output, err := cmd.Output()
if err != nil {
return "Error: " + err.Error()
}
return strings.TrimSpace(string(output))
}

func CollectRunningProcesses() string {
cmd := exec.Command("cmd", "/c", "tasklist /v")
output, err := cmd.Output()
if err != nil {
return "Error: " + err.Error()
}
return strings.TrimSpace(string(output))
}

func CollectServices() string {
cmd := exec.Command("cmd", "/c", "wmic service list brief")
output, err := cmd.Output()
if err != nil {
return "Error: " + err.Error()
}
return strings.TrimSpace(string(output))
}

func CollectAllEvidence() string {
var result strings.Builder
result.WriteString(CollectSystemInfo())
result.WriteString(CollectUsers())
result.WriteString(CollectGroups())
result.WriteString(CollectNetworkInfo())
result.WriteString(CollectRunningProcesses())
result.WriteString(CollectServices())
return result.String()
}
