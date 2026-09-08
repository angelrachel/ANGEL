package osint

import (
"context"
"fmt"
"net"
"time"
)

type SubdomainResult struct {
Domain string
Sub    string
IP     string
}

type SubdomainScanner struct {
Timeout time.Duration
}

func NewSubdomainScanner() *SubdomainScanner {
return &SubdomainScanner{Timeout: 10 * time.Second}
}

func (s *SubdomainScanner) Resolve(ctx context.Context, domain string) []SubdomainResult {
var results []SubdomainResult
subs := []string{"www", "mail", "admin", "api", "dev", "test", "portal", "vpn"}
for _, sub := range subs {
select {
case <-ctx.Done():
return results
default:
}
fqdn := fmt.Sprintf("%s.%s", sub, domain)
ips, err := net.LookupHost(fqdn)
if err != nil {
continue
}
results = append(results, SubdomainResult{Domain: domain, Sub: fqdn, IP: ips[0]})
}
return results
}

func (s *SubdomainScanner) Bruteforce(ctx context.Context, domain string, wordlist []string) []SubdomainResult {
var results []SubdomainResult
for _, word := range wordlist {
select {
case <-ctx.Done():
return results
default:
}
fqdn := fmt.Sprintf("%s.%s", word, domain)
ips, err := net.LookupHost(fqdn)
if err != nil {
continue
}
results = append(results, SubdomainResult{Domain: domain, Sub: fqdn, IP: ips[0]})
}
return results
}
