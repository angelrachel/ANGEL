package winrm

import (
"os/exec"
)

type ExecResult struct {
Command string
Output  string
}

type Executor struct {
Host string
User string
Pass string
}

func (e Executor) RunCommand(command string) ExecResult {
cmd := exec.Command("cmd", "/c", "powershell", "-Command", "New-PSSession -ComputerName "+e.Host+" -Credential (New-Object PSCredential('"+e.User+"',(ConvertTo-SecureString '"+e.Pass+"' -AsPlainText -Force))) | Invoke-Command -ScriptBlock {"+command+"}")
out, _ := cmd.Output()
return ExecResult{Command: command, Output: string(out)}
}

func (e Executor) RunShell(shell string) ExecResult {
cmd := exec.Command("cmd", "/c", "powershell", "-Command", "New-PSSession -ComputerName "+e.Host+" -Credential (New-Object PSCredential('"+e.User+"',(ConvertTo-SecureString '"+e.Pass+"' -AsPlainText -Force))) | Invoke-Command -ScriptBlock {"+shell+"}")
out, _ := cmd.Output()
return ExecResult{Command: shell, Output: string(out)}
}
