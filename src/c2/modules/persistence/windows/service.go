package persistence

import "os/exec"

type ServiceResult struct {
Name   string
Status string
}

func CreateService(name, binaryPath string) ServiceResult {
cmd := exec.Command("sc", "create", name, "binPath=", binaryPath, "start=", "auto")
cmd.Run()
return ServiceResult{Name: name, Status: "success"}
}

func StartService(name string) ServiceResult {
cmd := exec.Command("sc", "start", name)
cmd.Run()
return ServiceResult{Name: name, Status: "success"}
}

func DeleteService(name string) ServiceResult {
cmd := exec.Command("sc", "delete", name)
cmd.Run()
return ServiceResult{Name: name, Status: "success"}
}
