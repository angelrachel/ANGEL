package osint

import (
"encoding/json"
"fmt"
"net"
"os/exec"
"strings"
)

type DNSOSINT struct{}

func NewDNSOSINT() *DNSOSINT {
return &DNSOSINT{}
}

func (d *DNSOSINT) SubdomainEnum(domain string) ([]string, error) {
cmd := exec.Command("subfinder", "-d", domain, "-silent")
output, err := cmd.Output()
if err != nil {
return nil, err
}
lines := strings.Split(string(output), "\n")
var subdomains []string
for _, line := range lines {
if line != "" {
subdomains = append(subdomains, line)
}
}
return subdomains, nil
}

func (d *DNSOSINT) ReverseDNS(ip string) (string, error) {
names, err := net.LookupAddr(ip)
if err != nil {
return "", err
}
if len(names) > 0 {
return names[0], nil
}
return "", nil
}

func (d *DNSOSINT) ZoneTransfer(domain, ns string) ([]string, error) {
cmd := exec.Command("dig", "AXFR", domain, "@"+ns)
output, err := cmd.Output()
if err != nil {
return nil, err
}
lines := strings.Split(string(output), "\n")
var records []string
for _, line := range lines {
if strings.Contains(line, "A") || strings.Contains(line, "CNAME") {
records = append(records, line)
}
}
return records, nil
}
