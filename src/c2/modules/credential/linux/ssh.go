//go:build linux

package credential

import (
"os"
)

type SSHDump struct {
Path string
}

func NewSSHDump(path string) *SSHDump {
return &SSHDump{Path: path}
}

func (s *SSHDump) ReadPrivateKey() (string, error) {
file, err := os.Open(s.Path)
if err != nil {
return "", err
}
defer file.Close()
buf := make([]byte, 4096)
n, err := file.Read(buf)
if err != nil {
return "", err
}
return string(buf[:n]), nil
}

func (s *SSHDump) GetAllPrivateKeys() ([]string, error) {
dir := "/root/.ssh"
entries, err := os.ReadDir(dir)
if err != nil {
return nil, err
}
var keys []string
for _, entry := range entries {
if entry.IsDir() {
continue
}
keys = append(keys, dir+"/"+entry.Name())
}
return keys, nil
}

func (s *SSHDump) IsWritable() bool {
info, err := os.Stat(s.Path)
if err != nil {
return false
}
return info.Mode().Perm()&0200 != 0
}
