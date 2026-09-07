//go:build windows

package process_injection

import (
"golang.org/x/sys/windows"
"unsafe"
)

// procThreadContext returns a pointer to thread context structure
func procThreadContext() unsafe.Pointer {
var ctx [1024]byte
return unsafe.Pointer(&ctx)
}

// GetCurrentProcess returns handle to current process
func GetCurrentProcess() uintptr {
kernel32 := windows.NewLazyDLL("kernel32.dll")
procGetCurrentProcess := kernel32.NewProc("GetCurrentProcess")
ret, _, _ := procGetCurrentProcess.Call()
return ret
}
