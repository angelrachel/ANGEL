package gcp

import (
"os"
)

type ADC struct {
Path string
}

func NewADC(path string) *ADC {
return &ADC{Path: path}
}

func (a *ADC) ReadADC() (string, error) {
file, err := os.Open(a.Path)
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

func (a *ADC) GetProjectID() (string, error) {
content, err := a.ReadADC()
if err != nil {
return "", err
}
return content, nil
}

func (a *ADC) GetCredentials() (string, error) {
content, err := a.ReadADC()
if err != nil {
return "", err
}
return content, nil
}
