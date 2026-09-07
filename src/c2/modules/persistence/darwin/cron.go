package darwin

import (
"os"
)

type DarwinCronPersistence struct{}

func NewDarwinCronPersistence() *DarwinCronPersistence {
return &DarwinCronPersistence{}
}

func (c *DarwinCronPersistence) AddCron(command, schedule, user string) error {
cronLine := schedule + " " + command + "\n"
cronFile := "/etc/cron.d/angel"
if user != "" {
cronLine = schedule + " " + user + " " + command + "\n"
}
return os.WriteFile(cronFile, []byte(cronLine), 0644)
}
