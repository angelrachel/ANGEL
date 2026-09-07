package linux

import (
"os"
"strconv"
)

type CronPersistence struct{}

func NewCronPersistence() *CronPersistence {
return &CronPersistence{}
}

func (c *CronPersistence) AddCronJob(command, schedule, user string) error {
cronLine := schedule + " " + command + "\n"
cronFile := "/etc/cron.d/angel"
if user != "" {
cronLine = schedule + " " + user + " " + command + "\n"
}
return os.WriteFile(cronFile, []byte(cronLine), 0644)
}

func (c *CronPersistence) AddUserCron(command, schedule, user string) error {
cronFile := "/var/spool/cron/crontabs/" + user
cronLine := schedule + " " + command + "\n"
f, err := os.OpenFile(cronFile, os.O_APPEND|os.O_WRONLY, 0600)
if err != nil {
return err
}
defer f.Close()
_, err = f.WriteString(cronLine)
return err
}
