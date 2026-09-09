package persistence

import "os/exec"

type ScheduledTaskResult struct {
	Name   string
	Status string
}

type ScheduledTask struct {
	Name    string
	Command string
}

func (s ScheduledTask) Create() ScheduledTaskResult {
	cmd := exec.Command("schtasks", "/create", "/tn", s.Name, "/tr", s.Command, "/sc", "daily", "/f")
	cmd.Run()
	return ScheduledTaskResult{Name: s.Name, Status: "success"}
}

func (s ScheduledTask) Delete() ScheduledTaskResult {
	cmd := exec.Command("schtasks", "/delete", "/tn", s.Name, "/f")
	cmd.Run()
	return ScheduledTaskResult{Name: s.Name, Status: "success"}
}
