//go:build windows

package anti_analysis

import (
"golang.org/x/sys/windows"
"syscall"
)

// fs is a placeholder for filesystem operations
type fs struct{}

func (f *fs) Stat(path string) error {
_, err := windows.Stat(path)
return err
}
