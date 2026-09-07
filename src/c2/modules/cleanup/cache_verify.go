package cleanup

import (
"crypto/sha256"
"encoding/hex"
"os"
"path/filepath"
)

type CacheVerify struct{}

func NewCacheVerify() *CacheVerify {
return &CacheVerify{}
}

func (c *CacheVerify) ScanCache() ([]string, error) {
var files []string
cacheDirs := []string{
"/tmp/angel_cache",
"/var/cache/angel",
"/home/*/.cache/angel",
}

for _, dir := range cacheDirs {
matches, _ := filepath.Glob(dir)
for _, match := range matches {
filepath.Walk(match, func(path string, info os.FileInfo, err error) error {
if err == nil && !info.IsDir() {
files = append(files, path)
}
return nil
})
}
}
return files, nil
}

func (c *CacheVerify) VerifyClean() (bool, error) {
files, err := c.ScanCache()
if err != nil {
return false, err
}
return len(files) == 0, nil
}

func (c *CacheVerify) HashCache() (string, error) {
files, err := c.ScanCache()
if err != nil {
return "", err
}
hash := sha256.New()
for _, file := range files {
data, err := os.ReadFile(file)
if err == nil {
hash.Write(data)
}
}
return hex.EncodeToString(hash.Sum(nil)), nil
}
