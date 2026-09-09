package listener

import (
	"context"
	"net"
	"time"
)

type DNSResult struct {
	Query string
	IP    string
}

type DNSListener struct {
	Resolver *net.Resolver
	Server   string
}

func NewDNSListener(server string) *DNSListener {
	return &DNSListener{
		Resolver: &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				d := net.Dialer{Timeout: 5 * time.Second}
				return d.DialContext(ctx, "udp", server)
			},
		},
		Server: server,
	}
}

func (d *DNSListener) Lookup(ctx context.Context, domain string) []DNSResult {
	var results []DNSResult
	ips, err := d.Resolver.LookupHost(ctx, domain)
	if err != nil {
		return results
	}
	for _, ip := range ips {
		results = append(results, DNSResult{Query: domain, IP: ip})
	}
	return results
}

func (d *DNSListener) LookupTXT(ctx context.Context, domain string) []string {
	var results []string
	txts, err := d.Resolver.LookupTXT(ctx, domain)
	if err != nil {
		return results
	}
	return txts
}
