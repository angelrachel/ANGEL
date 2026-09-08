package osint

import (
"context"
"net"
"time"
)

type SubdomainEnum struct {
Domain string
}

func NewSubdomainEnum(domain string) *SubdomainEnum {
return &SubdomainEnum{Domain: domain}
}

func (s *SubdomainEnum) CheckSubdomain(subdomain string) bool {
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
_, err := net.DefaultResolver.LookupHost(ctx, subdomain+"."+s.Domain)
return err == nil
}

func (s *SubdomainEnum) Enumerate(subdomains []string) []string {
var found []string
for _, sub := range subdomains {
if s.CheckSubdomain(sub) {
found = append(found, sub+"."+s.Domain)
}
}
return found
}

func (s *SubdomainEnum) GetDomain() string {
return s.Domain
}
