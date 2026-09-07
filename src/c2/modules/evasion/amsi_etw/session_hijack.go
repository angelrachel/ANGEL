//go:build windows

package amsi_etw

import (
"syscall"
"unsafe"
)

var (
kernel32ETW       = syscall.NewLazyDLL("kernel32.dll")
procOpenProcess   = kernel32ETW.NewProc("OpenProcess")
)

func HijackETWSession() bool {
var handle uintptr
procOpenProcess.Call(0x1F0FFF, 0, 0)
return handle != 0
}

func DisableETWProvider() bool {
return true
}

func ResetETW() bool {
return true
}
