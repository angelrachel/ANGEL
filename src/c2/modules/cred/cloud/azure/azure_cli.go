package azure

import (
"os"
)

type AzureCLI struct {
Path string
}

func NewAzureCLI(path string) *AzureCLI {
return &AzureCLI{Path: path}
}

func (a *AzureCLI) ReadAzureProfile() (string, error) {
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

func (a *AzureCLI) GetSubscriptions() (string, error) {
content, err := a.ReadAzureProfile()
if err != nil {
return "", err
}
return content, nil
}

func (a *AzureCLI) GetTokens() (string, error) {
content, err := a.ReadAzureProfile()
if err != nil {
return "", err
}
return content, nil
}
