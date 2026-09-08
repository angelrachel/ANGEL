package osint

import (
"context"
"fmt"
"net"
"time"
)

type DNSResult struct {
Domain string
IP     string
}

type DNSEnum struct {
Timeout time.Duration
}

func NewDNSEnum() *DNSEnum {
return &DNSEnum{Timeout: 5 * time.Second}
}

func (d *DNSEnum) Lookup(ctx context.Context, domain string) DNSResult {
var result DNSResult
ips, err := net.LookupHost(domain)
if err != nil {
return result
}
result.Domain = domain
result.IP = ips[0]
return result
}

func (d *DNSEnum) ReverseLookup(ip string) DNSResult {
var result DNSResult
names, err := net.LookupAddr(ip)
if err != nil {
return result
}
result.Domain = names[0]
result.IP = ip
return result
}

func (d *DNSEnum) SubdomainBruteforce(ctx context.Context, domain string, words []string) []DNSResult {
var results []DNSResult
for _, word := range words {
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
results = append(results, DNSResult{Domain: fqdn, IP: ips[0]})
}
return results
}
