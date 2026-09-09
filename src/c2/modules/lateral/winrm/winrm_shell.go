//go:build windows

package winrm

import (
	"os/exec"
	"strings"
)

func WinRMShell(targetHost, username, password string) string {
	cmd := exec.Command("cmd", "/c", "powershell -Command \"$s = New-PSSession -ComputerName "+targetHost+" -Credential (New-Object PSCredential('"+username+"', (ConvertTo-SecureString '"+password+"' -AsPlainText -Force))); Invoke-Command -Session $s -ScriptBlock {whoami; hostname; systeminfo} | Out-String\"")
	output, err := cmd.Output()
	if err != nil {
		return "Error: " + err.Error()
	}
	return strings.TrimSpace(string(output))
}

func WinRMGetUsers(targetHost, username, password string) string {
	cmd := exec.Command("cmd", "/c", "powershell -Command \"$s = New-PSSession -ComputerName "+targetHost+" -Credential (New-Object PSCredential('"+username+"', (ConvertTo-SecureString '"+password+"' -AsPlainText -Force))); Invoke-Command -Session $s -ScriptBlock {Get-ChildItem C:\\Users\\ -Force | Select-Object -Property Name} | Out-String\"")
	output, err := cmd.Output()
	if err != nil {
		return "Error: " + err.Error()
	}
	return strings.TrimSpace(string(output))
}
