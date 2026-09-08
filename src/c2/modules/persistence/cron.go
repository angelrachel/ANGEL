package persistence

import "os/exec"

type CronResult struct {
Schedule string
Status   string
}

func AddCronJob(schedule, command string) CronResult {
cmd := exec.Command("bash", "-c", "(crontab -l; echo '"+schedule+" "+command+"') | crontab -")
cmd.Run()
return CronResult{Schedule: schedule, Status: "success"}
}

func DeleteCronJob() CronResult {
cmd := exec.Command("bash", "-c", "crontab -r")
cmd.Run()
return CronResult{Status: "success"}
}
