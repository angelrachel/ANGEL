//go:build windows

package sleep_masking

import (
"syscall"
"unsafe"
)

func RemoveGuardPageProtection(address uintptr) bool {
var oldProtect uint32
procVirtualProtect.Call(
address,
0x1000,
0x40,
uintptr(unsafe.Pointer(&oldProtect)),
)
return true
}

func SetGuardPageProtection(address uintptr) bool {
var oldProtect uint32
procVirtualProtect.Call(
address,
0x1000,
0x04,
uintptr(unsafe.Pointer(&oldProtect)),
)
return true
}
