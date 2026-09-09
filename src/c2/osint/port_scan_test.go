package osint

import (
	"net"
	"testing"
	"time"
)

func TestScanTCPFindsLocalListenerAndSkipsInvalidPorts(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer listener.Close()
	port := listener.Addr().(*net.TCPAddr).Port
	scanner := &Scanner{Timeout: 200 * time.Millisecond, Threads: 2}
	results := scanner.ScanTCP("127.0.0.1", []int{port, 0, 65536})
	if len(results) != 1 || results[0].Port != port || results[0].Status != "open" {
		t.Fatalf("unexpected scan results: %#v", results)
	}
}

func TestScanRangeRejectsReversedRange(t *testing.T) {
	scanner := &Scanner{Timeout: time.Millisecond, Threads: 1}
	if got := scanner.ScanRange("127.0.0.1", 20, 10); got != nil {
		t.Fatalf("expected nil for reversed range, got %#v", got)
	}
}
