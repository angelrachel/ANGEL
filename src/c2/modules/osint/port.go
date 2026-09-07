package osint

import (
"net"
"strconv"
"time"
)

type PortOSINT struct{}

func NewPortOSINT() *PortOSINT {
return &PortOSINT{}
}

func (p *PortOSINT) Scan(target string, ports []int) ([]int, error) {
var openPorts []int
for _, port := range ports {
conn, err := net.DialTimeout("tcp", target+":"+strconv.Itoa(port), 2*time.Second)
if err == nil {
conn.Close()
openPorts = append(openPorts, port)
}
}
return openPorts, nil
}

func (p *PortOSINT) ScanCommon(target string) ([]int, error) {
commonPorts := []int{21, 22, 23, 25, 53, 80, 110, 111, 135, 139, 143, 443, 445, 993, 995, 1723, 3306, 3389, 5900, 8080}
return p.Scan(target, commonPorts)
}

func (p *PortOSINT) GrabBanner(target string, port int) (string, error) {
conn, err := net.DialTimeout("tcp", target+":"+strconv.Itoa(port), 3*time.Second)
if err != nil {
return "", err
}
defer conn.Close()
buf := make([]byte, 1024)
n, err := conn.Read(buf)
if err != nil {
return "", err
}
return string(buf[:n]), nil
}
