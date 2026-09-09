package smbbeacon

import (
	"fmt"
	"os/exec"
)

type BeaconResult struct {
	Status string
	Pipe   string
}

type SMBBeacon struct {
	PipeName string
}

func (b SMBBeacon) Connect() BeaconResult {
	cmd := exec.Command("cmd", "/c", "echo", "beacon", ">", "\\\\.\\pipe\\"+b.PipeName)
	cmd.Run()
	return BeaconResult{Status: "success", Pipe: b.PipeName}
}

func (b SMBBeacon) SendMessage(message string) BeaconResult {
	cmd := exec.Command("cmd", "/c", "echo", message, ">", "\\\\.\\pipe\\"+b.PipeName)
	cmd.Run()
	return BeaconResult{Status: "success", Pipe: b.PipeName}
}

func (b SMBBeacon) ReceiveMessage() BeaconResult {
	cmd := exec.Command("cmd", "/c", "type", "\\\\.\\pipe\\"+b.PipeName)
	out, _ := cmd.Output()
	return BeaconResult{Status: string(out), Pipe: b.PipeName}
}

func FormatPipeName(name string) string {
	return fmt.Sprintf("\\\\.\\pipe\\%s", name)
}
