package rootkit

import "os/exec"

type JTAGResult struct {
	Status string
}

func JTAGDebug() JTAGResult {
	cmd := exec.Command("jtag.exe", "-debug", "chip.bin")
	cmd.Run()
	return JTAGResult{Status: "success"}
}
