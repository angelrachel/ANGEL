//go:build android

package persistence

import (
"os/exec"
)

func ForegroundService() bool {
cmd := exec.Command("am", "start-foreground-service", "-n", "com.angel/.ForegroundService")
err := cmd.Run()
if err != nil {
return false
}
return true
}

func ForegroundServiceStop() bool {
cmd := exec.Command("am", "stop-service", "-n", "com.angel/.ForegroundService")
err := cmd.Run()
if err != nil {
return false
}
return true
}
