package persistence

import (
	"os"
	"os/exec"
)

type StartupResult struct {
	Path   string
	Status string
}

func CreateStartupShortcut(targetPath, command string) StartupResult {
	cmd := exec.Command("powershell", "-Command", "New-Item -Path '"+targetPath+"' -ItemType SymbolicLink -Value '"+command+"' -Force")
	cmd.Run()
	return StartupResult{Path: targetPath, Status: "success"}
}

func DeleteStartupShortcut(targetPath string) StartupResult {
	os.Remove(targetPath)
	return StartupResult{Path: targetPath, Status: "success"}
}
