package persistence

import "os/exec"

type TaskResult struct {
	Name   string
	Status string
}

func CreateScheduledTask(name, command, trigger string) TaskResult {
	cmd := exec.Command("schtasks", "/create", "/tn", name, "/tr", command, "/sc", trigger, "/f")
	cmd.Run()
	return TaskResult{Name: name, Status: "success"}
}

func DeleteScheduledTask(name string) TaskResult {
	cmd := exec.Command("schtasks", "/delete", "/tn", name, "/f")
	cmd.Run()
	return TaskResult{Name: name, Status: "success"}
}
