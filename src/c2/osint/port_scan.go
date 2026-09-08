package osint

import (
"fmt"
"net"
"sync"
"time"
)

type ScanResult struct {
Host    string
Port    int
Status  string
}

type Scanner struct {
Timeout time.Duration
Threads int
}

func NewScanner() *Scanner {
return &Scanner{Timeout: 5 * time.Second, Threads: 100}
}

func (s *Scanner) ScanTCP(host string, ports []int) []ScanResult {
sem := make([]struct{}, s.Threads)
var wg sync.WaitGroup
var results []ScanResult
var mu sync.Mutex

for _, port := range ports {
wg.Add(1)
sem = append(sem, struct{}{})
go func(p int) {
defer wg.Done()
addr := fmt.Sprintf("%s:%d", host, p)
conn, err := net.DialTimeout("tcp", addr, s.Timeout)
if err != nil {
return
}
conn.Close()
mu.Lock()
results = append(results, ScanResult{Host: host, Port: p, Status: "open"})
mu.Unlock()
}(port)
}
wg.Wait()
return results
}

func (s *Scanner) ScanUDP(host string, ports []int) []ScanResult {
var results []ScanResult
for _, port := range ports {
addr := fmt.Sprintf("%s:%d", host, port)
conn, err := net.DialTimeout("udp", addr, s.Timeout)
if err != nil {
continue
}
conn.Close()
results = append(results, ScanResult{Host: host, Port: port, Status: "open"})
}
return results
}

func (s *Scanner) ScanRange(host string, start, end int) []ScanResult {
var ports []int
for i := start; i <= end; i++ {
ports = append(ports, i)
}
return s.ScanTCP(host, ports)
}
