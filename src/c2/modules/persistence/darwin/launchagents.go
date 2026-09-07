//go:build darwin

package persistence

import (
"os/exec"
)

func InstallLaunchAgent(plistPath string) bool {
cmd := exec.Command("cp", plistPath, "/Library/LaunchAgents/")
err := cmd.Run()
if err != nil {
return false
}
cmd = exec.Command("launchctl", "load", "/Library/LaunchAgents/"+plistPath)
err = cmd.Run()
if err != nil {
return false
}
return true
}

func RemoveLaunchAgent(plistPath string) bool {
cmd := exec.Command("launchctl", "unload", "/Library/LaunchAgents/"+plistPath)
err := cmd.Run()
if err != nil {
return false
}
return true
}
