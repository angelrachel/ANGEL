//go:build windows

package ad_recon

import (
"os/exec"
)

func EnumerateGroup(groupName string) string {
cmd := exec.Command("cmd", "/c", "net group "+groupName+" /domain")
output, err := cmd.Output()
if err != nil {
return ""
}
return string(output)
}

func EnumerateLocalGroup(groupName string) string {
cmd := exec.Command("cmd", "/c", "net localgroup "+groupName)
output, err := cmd.Output()
if err != nil {
return ""
}
return string(output)
}

func EnumerateDomainAdmins() string {
return EnumerateGroup("Domain Admins")
}

func EnumerateEnterpriseAdmins() string {
return EnumerateGroup("Enterprise Admins")
}

func EnumerateSchemaAdmins() string {
return EnumerateGroup("Schema Admins")
}

func EnumerateAllGroups() string {
cmd := exec.Command("cmd", "/c", "net group /domain")
output, err := cmd.Output()
if err != nil {
return ""
}
return string(output)
}
