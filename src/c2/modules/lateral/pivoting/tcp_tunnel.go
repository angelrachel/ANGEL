package pivoting

import (
	"context"
	"net"
	"strconv"
	"time"
)

type Tunnel struct {
	LocalPort  int
	RemoteAddr string
}

func (t *Tunnel) Start(ctx context.Context) error {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(t.LocalPort))
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
			go t.handle(ctx, conn)
		}
	}()

	select {
	case <-ctx.Done():
		return listener.Close()
	}
}

func (t *Tunnel) handle(ctx context.Context, client net.Conn) {
	defer client.Close()
	target, err := net.DialTimeout("tcp", t.RemoteAddr, 10*time.Second)
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
