package smb

import (
	"context"
	"os/exec"
	"time"
)

type ExecResult struct {
	Command string
	Output  string
	Status  string
}

type Execute struct {
	Host string
	User string
	Pass string
}

func (e Execute) PsExec(command string) ExecResult {
	cmd := exec.Command("psexec.exe", "\\\\"+e.Host, "-u", e.User, "-p", e.Pass, command)
	out, _ := cmd.Output()
	return ExecResult{Command: command, Output: string(out), Status: "success"}
}

func (e Execute) SMBExec() ExecResult {
	cmd := exec.Command("cmd", "/c", "net", "use", "\\\\"+e.Host+"\\C$", "/user:"+e.User, e.Pass)
	out, _ := cmd.Output()
	return ExecResult{Command: "net_use", Output: string(out), Status: "success"}
}

func (e Execute) RunPsExecWithTimeout(command string, timeout time.Duration) ExecResult {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, "psexec.exe", "\\\\"+e.Host, "-u", e.User, "-p", e.Pass, command)
	out, _ := cmd.Output()
	return ExecResult{Command: command, Output: string(out), Status: "success"}
}
