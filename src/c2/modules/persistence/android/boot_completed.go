//go:build android

package persistence

import (
"os/exec"
)

func BootReceiver() bool {
cmd := exec.Command("pm", "enable", "com.angel/.BootReceiver")
err := cmd.Run()
if err != nil {
return false
}
return true
}

func BootService() bool {
cmd := exec.Command("am", "start-foreground-service", "-n", "com.angel/.Service")
err := cmd.Run()
if err != nil {
return false
}
return true
}
