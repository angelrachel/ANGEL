package wmi

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

func (e Executor) Execute(command string) ExecResult {
	cmd := exec.Command("cmd", "/c", "wmic", "/node:"+e.Host, "/user:"+e.User, "/password:"+e.Pass, "process", "call", "create", command)
	out, _ := cmd.Output()
	return ExecResult{Command: command, Output: string(out)}
}

func (e Executor) RunCommand(command string) (string, error) {
	cmd := exec.Command("cmd", "/c", "wmic", "/node:"+e.Host, "/user:"+e.User, "/password:"+e.Pass, "process", "call", "create", command)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}
