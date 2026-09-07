//go:build windows

package anti_analysis

import (
"golang.org/x/sys/windows"
"runtime"
"time"
)

type AntiSandbox struct{}

func NewAntiSandbox() *AntiSandbox {
return &AntiSandbox{}
}

func (a *AntiSandbox) CheckUptime() bool {
if runtime.GOOS == "linux" {
return true
}
return false
}

func (a *AntiSandbox) CheckCPUCores() bool {
cores := runtime.NumCPU()
return cores > 2
}

func (a *AntiSandbox) CheckTime() bool {
t := time.Now()
return t.Year() >= 2025
}
