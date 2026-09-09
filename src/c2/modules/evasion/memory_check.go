package evasion

import (
	"os"
	"runtime"
	"strconv"
)

type MemoryResult struct {
	Memory int
	CPU    int
	Status string
}

func CheckSystemMemory() MemoryResult {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	result := MemoryResult{Memory: int(m.Sys / 1024 / 1024), CPU: runtime.NumCPU()}
	if m.Sys < 1024*1024*1024 {
		result.Status = "low_memory"
	} else {
		result.Status = "sufficient"
	}
	return result
}

func CheckDiskSize() MemoryResult {
	info, err := os.Stat("/")
	if err != nil {
		return MemoryResult{Status: "error"}
	}
	return MemoryResult{Memory: int(info.Size() / 1024 / 1024 / 1024), Status: "disk_size"}
}

func CheckCoreCount() MemoryResult {
	return MemoryResult{CPU: runtime.NumCPU(), Status: "core_count"}
}

func FormatMemory(size int) string {
	return strconv.Itoa(size) + " MB"
}
