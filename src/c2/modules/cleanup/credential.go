package cleanup

import "os"

type CredentialResult struct {
	Path   string
	Status string
}

type CredentialCleanup struct {
	KeyPath string
}

func (c CredentialCleanup) RevokeSSHKeys() CredentialResult {
	os.RemoveAll(c.KeyPath)
	return CredentialResult{Path: c.KeyPath, Status: "success"}
}

func (c CredentialCleanup) RotateTokens() CredentialResult {
	os.RemoveAll("/tmp/angel_tokens")
	return CredentialResult{Path: "/tmp/angel_tokens", Status: "success"}
}

func (c CredentialCleanup) CleanTempCredentials() CredentialResult {
	os.RemoveAll("/tmp/angel_creds")
	return CredentialResult{Path: "/tmp/angel_creds", Status: "success"}
}
