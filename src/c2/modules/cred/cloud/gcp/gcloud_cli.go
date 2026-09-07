package gcp

import (
"os"
)

type GCloudCLI struct {
Path string
}

func NewGCloudCLI(path string) *GCloudCLI {
return &GCloudCLI{Path: path}
}

func (g *GCloudCLI) ReadConfig() (string, error) {
file, err := os.Open(g.Path)
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

func (g *GCloudCLI) GetAccounts() (string, error) {
content, err := g.ReadConfig()
if err != nil {
return "", err
}
return content, nil
}

func (g *GCloudCLI) GetCredentials() (string, error) {
content, err := g.ReadConfig()
if err != nil {
return "", err
}
return content, nil
}
