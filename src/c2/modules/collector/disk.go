//go:build windows

package collector

import (
"os/exec"
)

func CollectDiskInfo() string {
cmd := exec.Command("cmd", "/c", "wmic logicaldisk get caption,size,freespace")
output, err := cmd.Output()
if err != nil {
return ""
}
return string(output)
}

func CollectDiskUsage() string {
cmd := exec.Command("cmd", "/c", "fsutil volume diskfree c:")
output, err := cmd.Output()
if err != nil {
return ""
}
return string(output)
}

func CollectDiskDrives() string {
cmd := exec.Command("cmd", "/c", "wmic logicaldisk get name")
output, err := cmd.Output()
if err != nil {
return ""
}
return string(output)
}

func CollectAllDisk() string {
var result string
result += CollectDiskInfo()
result += CollectDiskUsage()
result += CollectDiskDrives()
return result
}
