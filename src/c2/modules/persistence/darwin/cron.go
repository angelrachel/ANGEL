//go:build darwin

package persistence

import (
"os/exec"
"strings"
)

func InstallCron(command string) bool {
cmd := exec.Command("crontab", "-l")
output, err := cmd.Output()
if err != nil {
cmd = exec.Command("crontab", "-")
cmd.Stdin = strings.NewReader("* * * * * " + command)
cmd.Run()
return true
}
newCron := string(output) + "* * * * * " + command
cmd = exec.Command("crontab", "-")
cmd.Stdin = strings.NewReader(newCron)
err = cmd.Run()
if err != nil {
return false
}
return true
}

func UninstallCron() bool {
cmd := exec.Command("crontab", "-r")
err := cmd.Run()
if err != nil {
return false
}
return true
}
