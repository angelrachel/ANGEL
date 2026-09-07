package darwin

import (
"os"
"path/filepath"
)

type DarwinSSHAuthorizedKeys struct{}

func NewDarwinSSHAuthorizedKeys() *DarwinSSHAuthorizedKeys {
return &DarwinSSHAuthorizedKeys{}
}

func (s *DarwinSSHAuthorizedKeys) AddKey(publicKey, user string) error {
sshDir := "/Users/" + user + "/.ssh"
authFile := filepath.Join(sshDir, "authorized_keys")
f, err := os.OpenFile(authFile, os.O_APPEND|os.O_WRONLY, 0600)
if err != nil {
return err
}
defer f.Close()
_, err = f.WriteString("\n" + publicKey + "\n")
return err
}
