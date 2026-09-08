package pivoting

import (
"context"
"net"
"strconv"
"time"
)

type Proxy struct {
TargetAddr string
Host       string
Port       int
}

type ProxyConfig struct {
Port int
}

func (p *Proxy) Start(ctx context.Context) error {
listener, err := net.Listen("tcp", ":"+strconv.Itoa(p.Port))
if err != nil {
return err
}
defer listener.Close()

go func() {
for {
conn, err := listener.Accept()
if err != nil {
continue
}
go p.HandleClient(ctx, conn)
}
}()

select {
case <-ctx.Done():
return listener.Close()
}
}

func (p *Proxy) HandleClient(ctx context.Context, client net.Conn) {
defer client.Close()
target, err := net.DialTimeout("tcp", p.TargetAddr, 10*time.Second)
if err != nil {
return
}
defer target.Close()

go func() {
buf := make([]byte, 4096)
for {
n, err := client.Read(buf)
if err != nil {
return
}
target.Write(buf[:n])
}
}()

buf := make([]byte, 4096)
for {
n, err := target.Read(buf)
if err != nil {
return
}
client.Write(buf[:n])
}
}
