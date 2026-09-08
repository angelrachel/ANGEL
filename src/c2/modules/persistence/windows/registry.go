package persistence

import "os/exec"

type RegistryResult struct {
Key   string
Value string
}

func AddRunKey(key, command string) RegistryResult {
cmd := exec.Command("reg", "add", "HKLM\\SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Run", "/v", key, "/t", "REG_SZ", "/d", command, "/f")
cmd.Run()
return RegistryResult{Key: key, Value: command}
}

func AddRunOnceKey(key, command string) RegistryResult {
cmd := exec.Command("reg", "add", "HKLM\\SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\RunOnce", "/v", key, "/t", "REG_SZ", "/d", command, "/f")
cmd.Run()
return RegistryResult{Key: key, Value: command}
}

func DeleteRunKey(key string) RegistryResult {
cmd := exec.Command("reg", "delete", "HKLM\\SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Run", "/v", key, "/f")
cmd.Run()
return RegistryResult{Key: key}
}
