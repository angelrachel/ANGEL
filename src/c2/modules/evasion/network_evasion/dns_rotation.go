//go:build windows

package network_evasion

import (
"net"
"strings"
"time"
)

func DNSRotation() bool {
resolver := &net.Resolver{}
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
ips, err := resolver.LookupHost(ctx, "example.com")
if err != nil {
return false
}
return len(ips) > 0
}

func DNSRandomQuery() string {
domains := []string{"update.microsoft.com", "www.google.com", "login.live.com"}
return domains[time.Now().UnixNano()%int64(len(domains))]
}

func DNSSleep() {
time.Sleep(2 * time.Second)
}

func DNSGetIPs(domain string) string {
resolver := &net.Resolver{}
ips, err := resolver.LookupHost(context.Background(), domain)
if err != nil {
return "Error: " + err.Error()
}
return strings.Join(ips, ",")
}
