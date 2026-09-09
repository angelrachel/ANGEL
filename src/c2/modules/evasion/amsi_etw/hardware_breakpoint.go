//go:build windows

package amsi_etw

import (
	"syscall"
	"unsafe"
)

var (
	ntdllHard                  = syscall.NewLazyDLL("ntdll.dll")
	procNtSetInformationThread = ntdllHard.NewProc("NtSetInformationThread")
)

func PatchHardwareBreakpoint() bool {
	var handle uintptr
	procNtSetInformationThread.Call(handle, 0x11, 0, 0)
	return true
}

func ClearHardwareBreakpoint() bool {
	var handle uintptr
	procNtSetInformationThread.Call(handle, 0x11, 0, 0)
	return true
}

func DisableHardwareBreakpoints() bool {
	var ctx [512]byte
	procNtSetInformationThread.Call(0, 0x11, uintptr(unsafe.Pointer(&ctx[0])), 0)
	return true
}
