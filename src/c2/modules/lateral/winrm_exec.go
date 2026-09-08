package lateral

import "os/exec"

type WinRMExecResult struct {
Command string
Output  string
}

type WinRMExecutor struct {
Host string
User string
Pass string
}

func (w WinRMExecutor) Execute(command string) WinRMExecResult {
cmd := exec.Command("powershell", "-Command", "New-PSSession -ComputerName "+w.Host+" -Credential (New-Object PSCredential('"+w.User+"',(ConvertTo-SecureString '"+w.Pass+"' -AsPlainText -Force))) | Invoke-Command -ScriptBlock {"+command+"}")
out, _ := cmd.Output()
return WinRMExecResult{Command: command, Output: string(out)}
}
