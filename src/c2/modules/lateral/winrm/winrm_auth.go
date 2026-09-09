package winrm

import (
	"os/exec"
	"strings"
)

type Result struct {
	Command string
	Output  string
	Error   string
}

type Auth struct {
	Host string
	User string
	Pass string
}

func (a Auth) ExecuteScript(script string) Result {
	cmd := exec.Command("cmd", "/c", "powershell", "-Command", "New-PSSession -ComputerName "+a.Host+" -Credential (New-Object PSCredential('"+a.User+"',(ConvertTo-SecureString '"+a.Pass+"' -AsPlainText -Force))) | Invoke-Command -ScriptBlock {"+script+"}")
	out, err := cmd.Output()
	result := Result{Command: script}
	if err != nil {
		result.Error = err.Error()
		return result
	}
	result.Output = string(out)
	return result
}

func (a Auth) GetUsers() Result {
	cmd := exec.Command("cmd", "/c", "powershell", "-Command", "Get-WmiObject -Class Win32_UserAccount | Select-Object Name")
	out, err := cmd.Output()
	result := Result{Command: "users"}
	if err != nil {
		result.Error = err.Error()
		return result
	}
	result.Output = string(out)
	return result
}

func TrimOutput(input string) string {
	return strings.TrimSpace(input)
}
