//go:build darwin

package persistence

import (
"os/exec"
)

func InstallLaunchDaemon(plistPath string) bool {
cmd := exec.Command("cp", plistPath, "/Library/LaunchDaemons/")
err := cmd.Run()
if err != nil {
return false
}
cmd = exec.Command("launchctl", "load", "/Library/LaunchDaemons/"+plistPath)
err = cmd.Run()
if err != nil {
return false
}
return true
}

func RemoveLaunchDaemon(plistPath string) bool {
cmd := exec.Command("launchctl", "unload", "/Library/LaunchDaemons/"+plistPath)
err := cmd.Run()
if err != nil {
return false
}
return true
}
