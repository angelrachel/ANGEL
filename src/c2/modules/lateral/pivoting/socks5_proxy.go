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
Dialer     *net.Dialer
}

func NewSOCKS5Proxy(listenAddr, targetAddr string) *SOCKS5Proxy {
return &SOCKS5Proxy{
ListenAddr: listenAddr,
TargetAddr: targetAddr,
Dialer:     &net.Dialer{Timeout: 10 * time.Second},
}
}

func (p *SOCKS5Proxy) HandleClient(client net.Conn) {
defer client.Close()
remote, err := p.Dialer.Dial("tcp", p.TargetAddr)
if err != nil {
return
}
defer remote.Close()
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
if _, err := remote.Write(buffer[:n]); err != nil {
return
}
}
}()
go func() {
defer wg.Done()
buffer := make([]byte, 4096)
for {
n, err := remote.Read(buffer)
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

func (p *SOCKS5Proxy) GetListenAddr() string {
return p.ListenAddr
}

func (p *SOCKS5Proxy) GetTargetAddr() string {
return p.TargetAddr
}
