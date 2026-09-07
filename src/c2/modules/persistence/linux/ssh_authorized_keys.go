package linux

import (
"os"
"path/filepath"
)

type SSHAuthorizedKeys struct{}

func NewSSHAuthorizedKeys() *SSHAuthorizedKeys {
return &SSHAuthorizedKeys{}
}

func (s *SSHAuthorizedKeys) AddKey(publicKey, user string) error {
sshDir := "/home/" + user + "/.ssh"
if user == "root" {
sshDir = "/root/.ssh"
}
authFile := filepath.Join(sshDir, "authorized_keys")
f, err := os.OpenFile(authFile, os.O_APPEND|os.O_WRONLY, 0600)
if err != nil {
return err
}
defer f.Close()
_, err = f.WriteString("\n" + publicKey + "\n")
return err
}
