package darwin

import (
	"os/exec"
)

type Keychain struct{}

func NewKeychain() *Keychain {
	return &Keychain{}
}

func (k *Keychain) DumpKeychain() error {
	cmd := exec.Command("security", "dump-keychain")
	return cmd.Run()
}

func (k *Keychain) DumpKeychainWithPass(password string) error {
	cmd := exec.Command("security", "dump-keychain", "-p", password)
	return cmd.Run()
}
