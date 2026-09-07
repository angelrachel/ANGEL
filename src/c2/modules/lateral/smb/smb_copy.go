//go:build windows

package smb

import (
"os/exec"
)

func SMBFileCopy(targetHost, username, password, sourceFile, destPath string) bool {
cmd := exec.Command("cmd", "/c", "net use \\\\"+targetHost+" /user:"+username+" "+password)
err := cmd.Run()
if err != nil {
return false
}
cmd = exec.Command("cmd", "/c", "copy \\\\"+targetHost+"\\C$\\"+sourceFile+" "+destPath)
err = cmd.Run()
if err != nil {
return false
}
return true
}

func SMBFileUpload(targetHost, username, password, localFile, remotePath string) bool {
cmd := exec.Command("cmd", "/c", "net use \\\\"+targetHost+" /user:"+username+" "+password)
err := cmd.Run()
if err != nil {
return false
}
cmd = exec.Command("cmd", "/c", "copy "+localFile+" \\\\"+targetHost+"\\C$\\"+remotePath)
err = cmd.Run()
if err != nil {
return false
}
return true
}

func SMBFileDelete(targetHost, username, password, remotePath string) bool {
cmd := exec.Command("cmd", "/c", "net use \\\\"+targetHost+" /user:"+username+" "+password)
err := cmd.Run()
if err != nil {
return false
}
cmd = exec.Command("cmd", "/c", "del \\\\"+targetHost+"\\C$\\"+remotePath)
err = cmd.Run()
if err != nil {
return false
}
return true
}
