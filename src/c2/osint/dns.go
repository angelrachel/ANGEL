package osint

import (
"context"
"fmt"
"net"
"strings"
"time"
)

func LookupA(domain string) ([]string, error) {
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
return net.DefaultResolver.LookupHost(ctx, domain)
}

func LookupMX(domain string) (string, error) {
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
records, err := net.DefaultResolver.LookupMX(ctx, domain)
if err != nil {
return "", err
}
var recordsString []string
for _, record := range records {
recordsString = append(recordsString, record.Host)
}
return strings.Join(recordsString, ", "), nil
}

func LookupTXT(domain string) (string, error) {
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
records, err := net.DefaultResolver.LookupTXT(ctx, domain)
if err != nil {
return "", err
}
return strings.Join(records, ", "), nil
}

func ScanPorts(target string, ports []int) map[int]bool {
results := make(map[int]bool)
for _, port := range ports {
conn, err := net.DialTimeout("tcp", net.JoinHostPort(target, fmt.Sprintf("%d", port)), 2*time.Second)
if err != nil {
results[port] = false
continue
}
conn.Close()
results[port] = true
}
return results
}
