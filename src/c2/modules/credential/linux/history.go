//go:build linux

package credential

import (
"os"
)

type HistoryDump struct {
Path string
}

func NewHistoryDump(path string) *HistoryDump {
return &HistoryDump{Path: path}
}

func (h *HistoryDump) ReadHistory() (string, error) {
file, err := os.Open(h.Path)
if err != nil {
return "", err
}
defer file.Close()
buf := make([]byte, 8192)
n, err := file.Read(buf)
if err != nil {
return "", err
}
return string(buf[:n]), nil
}

func (h *HistoryDump) GetBashHistory() (string, error) {
return h.ReadHistory()
}

func (h *HistoryDump) GetZshHistory() (string, error) {
return h.ReadHistory()
}

func (h *HistoryDump) GetSSHKeys() ([]string, error) {
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
