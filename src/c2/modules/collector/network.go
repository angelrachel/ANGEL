//go:build windows

package collector

import (
"os/exec"
)

func CollectNetworkInfo() string {
cmd := exec.Command("cmd", "/c", "ipconfig /all")
output, err := cmd.Output()
if err != nil {
return ""
}
return string(output)
}

func CollectARPTable() string {
cmd := exec.Command("cmd", "/c", "arp -a")
output, err := cmd.Output()
if err != nil {
return ""
}
return string(output)
}

func CollectRoutingTable() string {
cmd := exec.Command("cmd", "/c", "route print")
output, err := cmd.Output()
if err != nil {
return ""
}
return string(output)
}

func CollectDNSInfo() string {
cmd := exec.Command("cmd", "/c", "ipconfig /displaydns")
output, err := cmd.Output()
if err != nil {
return ""
}
return string(output)
}

func CollectAllNetwork() string {
var result string
result += CollectNetworkInfo()
result += CollectARPTable()
result += CollectRoutingTable()
result += CollectDNSInfo()
return result
}
