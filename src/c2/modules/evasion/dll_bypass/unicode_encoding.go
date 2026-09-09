//go:build windows

package dll_bypass

import (
	"syscall"
	"unsafe"
)

var (
	kernel32U        = syscall.NewLazyDLL("kernel32.dll")
	procLoadLibraryW = kernel32U.NewProc("LoadLibraryW")
)

func LoadUnicodeDLL(dllPath string) bool {
	pathPtr, _ := syscall.UTF16PtrFromString(dllPath)
	procLoadLibraryW.Call(uintptr(unsafe.Pointer(pathPtr)))
	return true
}

func EncodeToUnicode(s string) string {
	runes := []rune(s)
	result := make([]uint16, len(runes))
	for i, r := range runes {
		result[i] = uint16(r)
	}
	return string(result)
}
