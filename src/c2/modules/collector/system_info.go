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

func CollectServices() string {
cmd := exec.Command("cmd", "/c", "wmic service list brief")
output, err := cmd.Output()
if err != nil {
return "Error: " + err.Error()
}
return strings.TrimSpace(string(output))
}

func CollectInstalledSoftware() string {
cmd := exec.Command("cmd", "/c", "wmic product get name,version")
output, err := cmd.Output()
if err != nil {
return "Error: " + err.Error()
}
return strings.TrimSpace(string(output))
}

func CollectDiskInfo() string {
cmd := exec.Command("cmd", "/c", "wmic logicaldisk get caption,size,freespace")
output, err := cmd.Output()
if err != nil {
return "Error: " + err.Error()
}
return strings.TrimSpace(string(output))
}

func CollectAll() string {
var result strings.Builder
result.WriteString(CollectSystemInfo())
result.WriteString(CollectUsers())
result.WriteString(CollectGroups())
result.WriteString(CollectServices())
result.WriteString(CollectInstalledSoftware())
result.WriteString(CollectDiskInfo())
return result.String()
}
