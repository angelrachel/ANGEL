package persistence

import "os/exec"

type BootResult struct {
Status string
}

func EnableBootReceiver() BootResult {
cmd := exec.Command("pm", "enable", "com.angel/.BootReceiver")
cmd.Run()
return BootResult{Status: "success"}
}

func StartForegroundService() BootResult {
cmd := exec.Command("am", "start-foreground-service", "-n", "com.angel/.ForegroundService")
cmd.Run()
return BootResult{Status: "success"}
}
