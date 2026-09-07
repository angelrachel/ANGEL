package windows

import (
"syscall"
"unsafe"
)

type Clipboard struct{}

func NewClipboard() *Clipboard {
return &Clipboard{}
}

func (c *Clipboard) Capture() (string, error) {
user32 := syscall.NewLazyDLL("user32.dll")
procOpenClipboard := user32.NewProc("OpenClipboard")
procGetClipboardData := user32.NewProc("GetClipboardData")
procCloseClipboard := user32.NewProc("CloseClipboard")

procOpenClipboard.Call(0)
handle, _, _ := procGetClipboardData.Call(1)
procCloseClipboard.Call()

if handle == 0 {
return "", nil
}
data := (*uint16)(unsafe.Pointer(handle))
return syscall.UTF16ToString((*[1 << 20]uint16)(unsafe.Pointer(data))[:]), nil
}
