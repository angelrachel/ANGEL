//go:build windows

package pivoting

import (
"net"
"sync"
"time"
)

type PortForward struct {
ListenAddr string
TargetAddr string
Dialer     *net.Dialer
}

func NewPortForward(listenAddr, targetAddr string) *PortForward {
return &PortForward{
ListenAddr: listenAddr,
TargetAddr: targetAddr,
Dialer:     &net.Dialer{Timeout: 10 * time.Second},
}
}

func (p *PortForward) Forward(client net.Conn) {
defer client.Close()
target, err := p.Dialer.Dial("tcp", p.TargetAddr)
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

func (p *PortForward) Start() error {
listener, err := net.Listen("tcp", p.ListenAddr)
if err != nil {
return err
}
defer listener.Close()
for {
conn, err := listener.Accept()
if err != nil {
continue
}
go p.Forward(conn)
}
}

func (p *PortForward) GetListenAddr() string {
return p.ListenAddr
}

func (p *PortForward) GetTargetAddr() string {
return p.TargetAddr
}
