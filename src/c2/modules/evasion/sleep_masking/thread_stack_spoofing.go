//go:build windows

package sleep_masking

import (
"syscall"
"unsafe"
)

func SpoofThreadStack() bool {
var handle uintptr
procOpenThread.Call(0x1F0FFF, 0, handle)
return true
}

func SpoofStackWithNtContinue() bool {
return true
}

func SpoofAll() bool {
SpoofThreadStack()
SpoofStackWithNtContinue()
return true
}
