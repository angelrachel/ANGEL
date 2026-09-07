package windows

import (
"io"
"os"
"path/filepath"
)

type StartupFolderPersistence struct{}

func NewStartupFolderPersistence() *StartupFolderPersistence {
return &StartupFolderPersistence{}
}

func (s *StartupFolderPersistence) AddToStartup(source, filename string) error {
startupFolder := os.Getenv("APPDATA") + `\Microsoft\Windows\Start Menu\Programs\Startup`
dest := filepath.Join(startupFolder, filename)
srcFile, err := os.Open(source)
if err != nil {
return err
}
defer srcFile.Close()
destFile, err := os.Create(dest)
if err != nil {
return err
}
defer destFile.Close()
_, err = io.Copy(destFile, srcFile)
return err
}
