package wmi

import "os/exec"

type PersistResult struct {
Command string
Status  string
}

type Persistence struct {
Host string
User string
Pass string
}

func (p Persistence) CreateService(serviceName, binaryPath string) PersistResult {
cmd := exec.Command("cmd", "/c", "wmic", "/node:"+p.Host, "/user:"+p.User, "/password:"+p.Pass, "service", "call", "create", serviceName, "binpath="+binaryPath)
out, _ := cmd.Output()
status := "success"
if len(out) == 0 {
status = "error"
}
return PersistResult{Command: "create_service", Status: status}
}

func (p Persistence) DeleteService(serviceName string) PersistResult {
cmd := exec.Command("cmd", "/c", "wmic", "/node:"+p.Host, "/user:"+p.User, "/password:"+p.Pass, "service", serviceName, "delete")
out, _ := cmd.Output()
status := "success"
if len(out) == 0 {
status = "error"
}
return PersistResult{Command: "delete_service", Status: status}
}
