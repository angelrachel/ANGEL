package linux

import (
"os"
)

type UdevPersistence struct{}

func NewUdevPersistence() *UdevPersistence {
return &UdevPersistence{}
}

func (u *UdevPersistence) AddUdevRule(name, command string) error {
ruleContent := `ACTION=="add", KERNEL=="sda", RUN+="` + command + `"`
return os.WriteFile("/etc/udev/rules.d/99-"+name+".rules", []byte(ruleContent), 0644)
}

func (u *UdevPersistence) TriggerUdev() error {
return os.WriteFile("/sys/kernel/uevent_helper", []byte(""), 0644)
}
