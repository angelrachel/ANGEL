package lateral

import "os/exec"

type WMIExecResult struct {
	Command string
	Output  string
}

type WMIExecutor struct {
	Host string
	User string
	Pass string
}

func (w WMIExecutor) Execute(command string) WMIExecResult {
	cmd := exec.Command("wmic", "/node:"+w.Host, "/user:"+w.User, "/password:"+w.Pass, "process", "call", "create", command)
	out, _ := cmd.Output()
	return WMIExecResult{Command: command, Output: string(out)}
}
