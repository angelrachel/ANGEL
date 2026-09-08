package smbbeacon

import (
"fmt"
"net"
"time"
)

type PeerToPeerResult struct {
Status string
Target string
}

type PeerToPeer struct {
Host string
Port int
}

func (p PeerToPeer) Connect() PeerToPeerResult {
conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", p.Host, p.Port), 5*time.Second)
if err != nil {
return PeerToPeerResult{Status: "error", Target: p.Host}
}
conn.Close()
return PeerToPeerResult{Status: "success", Target: p.Host}
}

func (p PeerToPeer) SendData(data string) PeerToPeerResult {
conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", p.Host, p.Port), 5*time.Second)
if err != nil {
return PeerToPeerResult{Status: "error", Target: p.Host}
}
defer conn.Close()
conn.Write([]byte(data))
return PeerToPeerResult{Status: "success", Target: p.Host}
}

func (p PeerToPeer) ReceiveData() PeerToPeerResult {
conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", p.Host, p.Port), 5*time.Second)
if err != nil {
return PeerToPeerResult{Status: "error", Target: p.Host}
}
defer conn.Close()
buf := make([]byte, 4096)
conn.Read(buf)
return PeerToPeerResult{Status: string(buf), Target: p.Host}
}
