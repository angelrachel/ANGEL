//go:build windows

package ad_recon

import (
"os/exec"
)

func EnumerateOUs() string {
cmd := exec.Command("cmd", "/c", "dsquery ou -name *")
output, err := cmd.Output()
if err != nil {
return ""
}
return string(output)
}

func EnumerateOUComputers(ou string) string {
cmd := exec.Command("cmd", "/c", "dsquery OU=\""+ou+"\",DC=angel,DC=local -filter \"(objectCategory=computer)\"")
output, err := cmd.Output()
if err != nil {
return ""
}
return string(output)
}

func EnumerateOUUsers(ou string) string {
cmd := exec.Command("cmd", "/c", "dsquery OU=\""+ou+"\",DC=angel,DC=local -filter \"(objectCategory=person)\"")
output, err := cmd.Output()
if err != nil {
return ""
}
return string(output)
}

func EnumerateAllOUs() string {
return EnumerateOUs()
}
