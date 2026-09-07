package sleep_masking

import (
"syscall"
"unsafe"
)

type VirtualProtectSleep struct{}

func NewVirtualProtectSleep() *VirtualProtectSleep {
return &VirtualProtectSleep{}
}

func (v *VirtualProtectSleep) Sleep(ms int) error {
kernel32 := syscall.NewLazyDLL("kernel32.dll")
procVirtualProtect := kernel32.NewProc("VirtualProtect")
procSleep := kernel32.NewProc("Sleep")

var oldProtect uint32
addr := uintptr(unsafe.Pointer(&oldProtect))
procVirtualProtect.Call(addr, 1, 0x80, uintptr(unsafe.Pointer(&oldProtect)))
procSleep.Call(uintptr(ms))
procVirtualProtect.Call(addr, 1, oldProtect, uintptr(unsafe.Pointer(&oldProtect)))

return nil
}
