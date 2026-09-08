//go:build windows

package network_evasion

import (
"context"
"net"
"time"
)

func DNSRotation(domain string) bool {
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
_, err := net.DefaultResolver.LookupHost(ctx, domain)
return err == nil
}

func DNSRandomQuery(domain string) string {
return "random" + time.Now().Format("150405") + "." + domain
}
