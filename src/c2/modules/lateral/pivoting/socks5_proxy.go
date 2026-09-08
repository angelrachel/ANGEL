//go:build windows

package pivoting

import (
"net"
"sync"
"time"
)

type SOCKS5Proxy struct {
ListenAddr string
TargetAddr string
Client     *net.Dialer
}

func NewSOCKS5Proxy(listenAddr, targetAddr string) *SOCKS5Proxy {
return &SOCKS5Proxy{
ListenAddr: listenAddr,
TargetAddr: targetAddr,
Client:     &net.Dialer{Timeout: 10 * time.Second},
}
}

func (p *SOCKS5Proxy) HandleClient(conn net.Conn) {
defer conn.Close()
target, err := p.Client.Dial("tcp", p.TargetAddr)
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
n, err := conn.Read(buffer)
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
if _, err := conn.Write(buffer[:n]); err != nil {
return
}
}
}()

wg.Wait()
}

func (p *SOCKS5Proxy) Start() error {
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
go p.HandleClient(conn)
}
}
