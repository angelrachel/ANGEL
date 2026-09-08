package collector

import "os/exec"

type SystemInfoResult struct {
Command string
Output  string
}

func GetSystemInfo() SystemInfoResult {
cmd := exec.Command("systeminfo")
out, _ := cmd.Output()
return SystemInfoResult{Command: "systeminfo", Output: string(out)}
}

func GetProcessList() SystemInfoResult {
cmd := exec.Command("tasklist")
out, _ := cmd.Output()
return SystemInfoResult{Command: "tasklist", Output: string(out)}
}

func GetNetworkInfo() SystemInfoResult {
cmd := exec.Command("ipconfig", "/all")
out, _ := cmd.Output()
return SystemInfoResult{Command: "ipconfig", Output: string(out)}
}
