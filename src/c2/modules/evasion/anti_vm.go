package evasion

import (
"os"
"strings"
)

type VMResult struct {
Status string
}

func CheckCPUInfo() VMResult {
data, err := os.ReadFile("/proc/cpuinfo")
if err != nil {
return VMResult{Status: "error"}
}
if strings.Contains(string(data), "hypervisor") || strings.Contains(string(data), "QEMU") || strings.Contains(string(data), "VMware") || strings.Contains(string(data), "VirtualBox") {
return VMResult{Status: "vm_detected"}
}
return VMResult{Status: "not_vm"}
}

func CheckSystemVendor() VMResult {
data, err := os.ReadFile("/sys/class/dmi/id/sys_vendor")
if err != nil {
return VMResult{Status: "error"}
}
if strings.Contains(string(data), "QEMU") || strings.Contains(string(data), "VMware") || strings.Contains(string(data), "Microsoft Corporation") {
return VMResult{Status: "vm_detected"}
}
return VMResult{Status: "not_vm"}
}
