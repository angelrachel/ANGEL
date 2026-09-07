package osint

import (
"context"
"net"
"strconv"
"sync"
"time"
)

type PortScanResult struct {
IP       string
Port     int
Open     bool
Protocol string
}

type PortScanner struct {
Target   string
Timeout  time.Duration
Protocol string
}

func NewPortScanner(target string, timeout time.Duration, protocol string) *PortScanner {
return &PortScanner{
Target:   target,
Timeout:  timeout,
Protocol: protocol,
}
}

func (p *PortScanner) ScanPort(ctx context.Context, port int) PortScanResult {
address := net.JoinHostPort(p.Target, strconv.Itoa(port))
conn, err := net.DialTimeout(p.Protocol, address, p.Timeout)
if err != nil {
return PortScanResult{IP: p.Target, Port: port, Open: false, Protocol: p.Protocol}
}
conn.Close()
return PortScanResult{IP: p.Target, Port: port, Open: true, Protocol: p.Protocol}
}

func (p *PortScanner) ScanRange(ctx context.Context, startPort, endPort int) ([]PortScanResult, error) {
var results []PortScanResult
var wg sync.WaitGroup
resultChan := make(chan PortScanResult, endPort-startPort+1)

for i := startPort; i <= endPort; i++ {
wg.Add(1)
go func(port int) {
defer wg.Done()
resultChan <- p.ScanPort(ctx, port)
}(i)
}

wg.Wait()
close(resultChan)

for res := range resultChan {
if res.Open {
results = append(results, res)
}
}
return results, nil
}

func (p *PortScanner) ScanCommonPorts(ctx context.Context) ([]PortScanResult, error) {
commonPorts := []int{21, 22, 23, 25, 53, 80, 110, 111, 135, 139, 143, 443, 445, 993, 995, 1723, 3306, 3389, 5432, 5900, 6379, 8080, 8443, 9200, 27017}
return p.ScanRange(ctx, 1, 1000)
}

func (p *PortScanner) GetTarget() string {
return p.Target
}

func (p *PortScanner) GetTimeout() time.Duration {
return p.Timeout
}
