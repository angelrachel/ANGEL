//go:build windows

package pivoting

import (
"net"
"sync"
"time"
)

type TCPTunnel struct {
ListenAddr string
TargetAddr string
Dialer     *net.Dialer
}

func NewTCPTunnel(listenAddr, targetAddr string) *TCPTunnel {
return &TCPTunnel{
ListenAddr: listenAddr,
TargetAddr: targetAddr,
Dialer:     &net.Dialer{Timeout: 10 * time.Second},
}
}

func (t *TCPTunnel) Tunnel(client net.Conn) {
defer client.Close()
target, err := t.Dialer.Dial("tcp", t.TargetAddr)
if err != nil {
return
}
defer target.Close()
var wg sync.WaitGroup
wg.Add(2)
go func() {
defer wg.Done()
buffer := make([]byte, 4096)
for {
n, err := client.Read(buffer)
if err != nil {
return
}
if _, err := target.Write(buffer[:n]); err != nil {
return
}
}
}()
go func() {
defer wg.Done()
buffer := make([]byte, 4096)
for {
n, err := target.Read(buffer)
if err != nil {
return
}
if _, err := client.Write(buffer[:n]); err != nil {
return
}
}
}()
wg.Wait()
}

func (t *TCPTunnel) Start() error {
listener, err := net.Listen("tcp", t.ListenAddr)
if err != nil {
return err
}
defer listener.Close()
for {
conn, err := listener.Accept()
if err != nil {
continue
}
go t.Tunnel(conn)
}
}

func (t *TCPTunnel) GetListenAddr() string {
return t.ListenAddr
}

func (t *TCPTunnel) GetTargetAddr() string {
return t.TargetAddr
}
