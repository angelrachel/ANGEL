package osint

import (
	"net"
	"strconv"
	"sync"
	"time"
)

type ScanResult struct {
	Host   string
	Port   int
	Status string
}

type Scanner struct {
	Timeout time.Duration
	Threads int
}

func NewScanner() *Scanner {
	return &Scanner{Timeout: 5 * time.Second, Threads: 100}
}

func (s *Scanner) ScanTCP(host string, ports []int) []ScanResult {
	threads := s.Threads
	if threads < 1 {
		threads = 1
	}
	timeout := s.Timeout
	if timeout <= 0 {
		timeout = time.Second
	}
	sem := make(chan struct{}, threads)
	var wg sync.WaitGroup
	var results []ScanResult
	var mu sync.Mutex

	for _, port := range ports {
		if port < 1 || port > 65535 {
			continue
		}
		wg.Add(1)
		go func(p int) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			addr := net.JoinHostPort(host, strconv.Itoa(p))
			conn, err := net.DialTimeout("tcp", addr, timeout)
			if err != nil {
				return
			}
			_ = conn.Close()
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
	timeout := s.Timeout
	if timeout <= 0 {
		timeout = time.Second
	}
	for _, port := range ports {
		if port < 1 || port > 65535 {
			continue
		}
		addr := net.JoinHostPort(host, strconv.Itoa(port))
		conn, err := net.DialTimeout("udp", addr, timeout)
		if err != nil {
			continue
		}
		_ = conn.Close()
		results = append(results, ScanResult{Host: host, Port: port, Status: "open"})
	}
	return results
}

func (s *Scanner) ScanRange(host string, start, end int) []ScanResult {
	if start < 1 {
		start = 1
	}
	if end > 65535 {
		end = 65535
	}
	if end < start {
		return nil
	}
	ports := make([]int, 0, end-start+1)
	for i := start; i <= end; i++ {
		ports = append(ports, i)
	}
	return s.ScanTCP(host, ports)
}
