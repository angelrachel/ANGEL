//go:build windows

package ad_recon

import (
"os/exec"
)

func EnumerateUserSessions() string {
cmd := exec.Command("cmd", "/c", "qwinsta")
output, err := cmd.Output()
if err != nil {
return ""
}
return string(output)
}

func EnumerateLoggedInUsers() string {
cmd := exec.Command("cmd", "/c", "query user")
output, err := cmd.Output()
if err != nil {
return ""
}
return string(output)
}

func EnumerateUserPrivileges() string {
cmd := exec.Command("cmd", "/c", "whoami /priv")
output, err := cmd.Output()
if err != nil {
return ""
}
return string(output)
}

func EnumerateUserGroups() string {
cmd := exec.Command("cmd", "/c", "whoami /groups")
output, err := cmd.Output()
if err != nil {
return ""
}
return string(output)
}

func EnumerateUserToken() string {
cmd := exec.Command("cmd", "/c", "whoami /all")
output, err := cmd.Output()
if err != nil {
return ""
}
return string(output)
}

func DumpUserSessions() string {
var result string
result += EnumerateUserSessions()
result += EnumerateLoggedInUsers()
result += EnumerateUserPrivileges()
result += EnumerateUserGroups()
result += EnumerateUserToken()
return result
}
