//go:build windows

package evasion

import (
"syscall"
"unsafe"
)

var (
kernel32Memory = syscall.NewLazyDLL("kernel32.dll")
)

func CheckMemorySize() uint32 {
return 0
}

func CheckMemoryProtection() uint32 {
return 0
}

func IsMemoryExecutable() bool {
return true
}

func CheckProcessMemory() bool {
return true
}

func DumpProcessMemory() bool {
return true
}
