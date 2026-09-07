package cloud

import (
"os"
"path/filepath"
)

type AWS struct{}

func NewAWS() *AWS {
return &AWS{}
}

func (a *AWS) ExtractCredentials() (string, error) {
path := filepath.Join(os.Getenv("HOME"), ".aws", "credentials")
data, err := os.ReadFile(path)
if err != nil {
return "", err
}
return string(data), nil
}

func (a *AWS) ExtractConfig() (string, error) {
path := filepath.Join(os.Getenv("HOME"), ".aws", "config")
data, err := os.ReadFile(path)
if err != nil {
return "", err
}
return string(data), nil
}
