package cleanup

import (
"os"
)

type CredentialCleanup struct{}

func NewCredentialCleanup() *CredentialCleanup {
return &CredentialCleanup{}
}

func (c *CredentialCleanup) RevokeTempCreds() error {
// Revoke temporary credentials
// In real implementation, would call AWS/Azure/GCP APIs
return os.RemoveAll("/tmp/angel_creds")
}

func (c *CredentialCleanup) RotateTokens() error {
// Rotate tokens and keys
return os.RemoveAll("/tmp/angel_tokens")
}

func (c *CredentialCleanup) DeleteSSHKeys() error {
// Delete SSH keys used during assessment
return os.RemoveAll("/tmp/angel_ssh_keys")
}
