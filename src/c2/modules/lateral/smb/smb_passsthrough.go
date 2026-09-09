package smb

import (
	"fmt"
	"net"
	"time"
)

type SmbError struct {
	Message string
}

func (e SmbError) Error() string {
	return e.Message
}

type Target struct {
	Host string
	Port int
}

type Request struct {
	Target Target
	User   string
	Hash   string
}

type Response struct {
	Status  string
	Message string
}

type SmbClient struct {
	Conn net.Conn
}

func (c *SmbClient) Read(p []byte) (n int, err error) {
	return c.Conn.Read(p)
}

func (c *SmbClient) Write(p []byte) (n int, err error) {
	return c.Conn.Write(p)
}

func (c *SmbClient) Close() error {
	return c.Conn.Close()
}

func (c *SmbClient) LocalAddr() net.Addr {
	return c.Conn.LocalAddr()
}

func (c *SmbClient) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *SmbClient) SetDeadline(t time.Time) error {
	return c.Conn.SetDeadline(t)
}

func (c *SmbClient) SetReadDeadline(t time.Time) error {
	return c.Conn.SetReadDeadline(t)
}

func (c *SmbClient) SetWriteDeadline(t time.Time) error {
	return c.Conn.SetWriteDeadline(t)
}

func Dial(t Target) (*SmbClient, error) {
	if t.Host == "" {
		return &SmbClient{}, &SmbError{Message: "host required"}
	}
	if t.Port == 0 {
		t.Port = 445
	}
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", t.Host, t.Port), 10*time.Second)
	if err != nil {
		return &SmbClient{}, &SmbError{Message: err.Error()}
	}
	return &SmbClient{Conn: conn}, nil
}

func Negotiate(conn net.Conn) error {
	if conn == nil {
		return &SmbError{Message: "connection required"}
	}
	packet := []byte{
		0x00, 0x00, 0x00, 0x2f,
		0xff, 0x53, 0x4d, 0x42,
		0x72, 0x00, 0x00, 0x00,
		0x00, 0x18, 0x53, 0xc8,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
	}
	_, err := conn.Write(packet)
	return err
}

func SessionSetup(conn net.Conn, req Request) error {
	if conn == nil {
		return &SmbError{Message: "connection required"}
	}
	if req.User == "" || req.Hash == "" {
		return &SmbError{Message: "credentials required"}
	}

	packet := make([]byte, 128)
	copy(packet[0:4], []byte{0x00, 0x00, 0x00, 0x7f})
	copy(packet[4:8], []byte{0xff, 0x53, 0x4d, 0x42})
	packet[8] = 0x73
	packet[9] = 0x00

	payload := []byte{
		0x60, 0x48, 0x06, 0x06,
		0x2b, 0x06, 0x01, 0x05,
		0x05, 0x02, 0xa0, 0x3e,
		0x30, 0x3c, 0xa0, 0x0e,
		0x30, 0x0c, 0x06, 0x0a,
		0x2b, 0x06, 0x01, 0x04,
		0x01, 0x82, 0x37, 0x02,
		0x02, 0x0e, 0xa2, 0x2a,
		0x04, 0x28,
	}
	payload = append(payload, []byte(req.User)...)
	payload = append(payload, 0x00)
	payload = append(payload, []byte(req.Hash)...)
	payload = append(payload, make([]byte, 16)...)

	packet = append(packet, payload...)
	_, err := conn.Write(packet)
	return err
}

func Run(req Request) (Response, error) {
	client, err := Dial(req.Target)
	if err != nil {
		return Response{Status: "error", Message: err.Error()}, &SmbError{Message: err.Error()}
	}
	defer client.Close()

	if err := Negotiate(client.Conn); err != nil {
		return Response{Status: "error", Message: err.Error()}, &SmbError{Message: err.Error()}
	}

	if err := SessionSetup(client.Conn, req); err != nil {
		return Response{Status: "error", Message: err.Error()}, &SmbError{Message: err.Error()}
	}

	return Response{Status: "success", Message: "auth relayed"}, nil
}
