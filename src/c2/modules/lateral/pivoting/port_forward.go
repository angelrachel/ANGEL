package pivoting

import (
	"context"
	"net"
	"strconv"
	"time"
)

type Forwarder struct {
	LocalPort  int
	TargetAddr string
}

func (f *Forwarder) Start(ctx context.Context) error {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(f.LocalPort))
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
			go f.handle(ctx, conn)
		}
	}()

	select {
	case <-ctx.Done():
		return listener.Close()
	}
}

func (f *Forwarder) handle(ctx context.Context, client net.Conn) {
	defer client.Close()
	target, err := net.DialTimeout("tcp", f.TargetAddr, 10*time.Second)
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
