//go:build linux

package credential

import (
"os"
"strings"
"syscall"
)

type ShadowDump struct {
Path string
}

func NewShadowDump(path string) *ShadowDump {
return &ShadowDump{Path: path}
}

func (s *ShadowDump) ReadShadow() ([]byte, error) {
file, err := os.Open(s.Path)
if err != nil {
return nil, err
}
defer file.Close()
buf := make([]byte, 4096)
n, err := file.Read(buf)
if err != nil {
return nil, err
}
return buf[:n], nil
}

func (s *ShadowDump) IsReadable() bool {
var stat syscall.Stat_t
if err := syscall.Stat(s.Path, &stat); err != nil {
return false
}
return stat.Mode&syscall.S_IRUSR != 0
}

func (s *ShadowDump) DumpAllHashes() ([]string, error) {
content, err := s.ReadShadow()
if err != nil {
return nil, err
}
var hashes []string
lines := strings.Split(string(content), "\n")
for _, line := range lines {
if strings.Contains(line, ":") {
hashes = append(hashes, line)
}
}
return hashes, nil
}
