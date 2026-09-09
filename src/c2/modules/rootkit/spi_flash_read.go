package rootkit

import "os/exec"

type SPIReadResult struct {
	Status string
}

func SPIRead() SPIReadResult {
	cmd := exec.Command("flashrom.exe", "-p", "spi", "-r", "firmware.bin")
	cmd.Run()
	return SPIReadResult{Status: "success"}
}
