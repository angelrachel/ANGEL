package windows

import (
"os/exec"
)

type ScheduledTaskPersistence struct{}

func NewScheduledTaskPersistence() *ScheduledTaskPersistence {
return &ScheduledTaskPersistence{}
}

func (s *ScheduledTaskPersistence) CreateTask(name, command, trigger string) error {
cmd := exec.Command("schtasks", "/create", "/tn", name, "/tr", command, "/sc", trigger, "/f")
return cmd.Run()
}

func (s *ScheduledTaskPersistence) CreateOnLogon(name, command string) error {
return s.CreateTask(name, command, "ONLOGON")
}

func (s *ScheduledTaskPersistence) CreateOnBoot(name, command string) error {
return s.CreateTask(name, command, "ONBOOT")
}

func (s *ScheduledTaskPersistence) CreateDaily(name, command string, time string) error {
return s.CreateTask(name, command, "DAILY")
}
