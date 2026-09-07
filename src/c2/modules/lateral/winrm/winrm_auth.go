//go:build windows

package winrm

import (
"os/exec"
"strings"
)

func WinRMLogin(targetHost, username, password string) bool {
cmd := exec.Command("cmd", "/c", "powershell -Command \"$s = New-PSSession -ComputerName "+targetHost+" -Credential (New-Object PSCredential('"+username+"', (ConvertTo-SecureString '"+password+"' -AsPlainText -Force))); Get-PSSession | Select-Object -Property Name\"")
err := cmd.Run()
if err != nil {
return false
}
return true
}

func WinRMQuery(targetHost, username, password, query string) string {
cmd := exec.Command("cmd", "/c", "powershell -Command \"$s = New-PSSession -ComputerName "+targetHost+" -Credential (New-Object PSCredential('"+username+"', (ConvertTo-SecureString '"+password+"' -AsPlainText -Force))); Invoke-Command -Session $s -ScriptBlock {"+query+"} | Out-String\"")
output, err := cmd.Output()
if err != nil {
return "Error: " + err.Error()
}
return strings.TrimSpace(string(output))
}

func WinRMExec(targetHost, username, password, command string) bool {
cmd := exec.Command("cmd", "/c", "powershell -Command \"$s = New-PSSession -ComputerName "+targetHost+" -Credential (New-Object PSCredential('"+username+"', (ConvertTo-SecureString '"+password+"' -AsPlainText -Force))); Invoke-Command -Session $s -ScriptBlock {"+command+"}\"")
err := cmd.Run()
if err != nil {
return false
}
return true
}
