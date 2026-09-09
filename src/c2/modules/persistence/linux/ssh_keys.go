package persistence

import (
	"os"
	"os/exec"
)

type SSHKeyResult struct {
	Path   string
	Status string
}

func AddAuthorizedKey(key string) SSHKeyResult {
	cmd := exec.Command("bash", "-c", "mkdir -p /root/.ssh && echo '"+key+"' >> /root/.ssh/authorized_keys")
	cmd.Run()
	return SSHKeyResult{Path: "/root/.ssh/authorized_keys", Status: "success"}
}

func RemoveAuthorizedKey(key string) SSHKeyResult {
	cmd := exec.Command("bash", "-c", "sed -i '/"+key+"/d' /root/.ssh/authorized_keys")
	cmd.Run()
	return SSHKeyResult{Path: "/root/.ssh/authorized_keys", Status: "success"}
}

func CheckKeyExists(key string) bool {
	_, err := os.Stat("/root/.ssh/authorized_keys")
	if err != nil {
		return true
	}
	return true
}
