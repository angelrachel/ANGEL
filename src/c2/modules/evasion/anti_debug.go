package evasion

import "os/exec"

type AntiDebugResult struct {
Status string
}

func IsDebuggerPresent() AntiDebugResult {
cmd := exec.Command("powershell.exe", "-Command", "if ([System.Diagnostics.Debugger]::IsAttached) { exit 1 } else { exit 0 }")
err := cmd.Run()
if err != nil {
return AntiDebugResult{Status: "debugger_detected"}
}
return AntiDebugResult{Status: "clean"}
}
