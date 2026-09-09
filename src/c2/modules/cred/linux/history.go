//go:build linux

package cred

import (
	"os"
)

func ReadBashHistory() string {
	data, err := os.ReadFile("/root/.bash_history")
	if err != nil {
		return ""
	}
	return string(data)
}

func ReadZshHistory() string {
	data, err := os.ReadFile("/root/.zsh_history")
	if err != nil {
		return ""
	}
	return string(data)
}
