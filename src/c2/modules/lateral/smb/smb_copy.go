package smb

import (
	"fmt"
	"os/exec"
)

type CopyResult struct {
	Command string
	Status  string
}

type Copy struct {
	Host string
	User string
	Pass string
}

func (c Copy) Upload(localPath, remotePath string) CopyResult {
	cmd := exec.Command("cmd", "/c", "net", "use", "\\\\"+c.Host+"\\C$", "/user:"+c.User, c.Pass)
	cmd.Run()
	cmd = exec.Command("cmd", "/c", "copy", localPath, "\\\\"+c.Host+"\\C$\\"+remotePath)
	cmd.Run()
	return CopyResult{Command: "upload", Status: "success"}
}

func (c Copy) Delete(remotePath string) CopyResult {
	cmd := exec.Command("cmd", "/c", "del", "\\\\"+c.Host+"\\C$\\"+remotePath)
	cmd.Run()
	return CopyResult{Command: "delete", Status: "success"}
}

func FormatUNC(host, share, path string) string {
	return fmt.Sprintf("\\\\%s\\%s\\%s", host, share, path)
}
