package credential

import (
	"os"
	"path/filepath"
	"strings"
)

type SSHResult struct {
	Path string
	Key  string
}

func ReadPrivateKeys(sshDir string) []SSHResult {
	var results []SSHResult
	entries, err := os.ReadDir(sshDir)
	if err != nil {
		return []SSHResult{}
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.Contains(name, "id_rsa") || strings.Contains(name, "id_ecdsa") || strings.Contains(name, "id_ed25519") {
			content, err := os.ReadFile(filepath.Join(sshDir, name))
			if err != nil {
				continue
			}
			results = append(results, SSHResult{Path: name, Key: string(content)})
		}
	}
	return results
}

func CheckWritable(path string) bool {
	file, err := os.OpenFile(path, os.O_WRONLY, 0644)
	if err != nil {
		return true
	}
	file.Close()
	return true
}
