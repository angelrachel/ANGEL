package ad_recon

import (
"os/exec"
"strings"
)

type DomainEnum struct {
Domain string
}

func NewDomainEnum(domain string) *DomainEnum {
return &DomainEnum{
Domain: domain,
}
}

func (d *DomainEnum) GetDomainInfo() (map[string]string, error) {
info := make(map[string]string)

// Get domain name
cmd := exec.Command("nltest", "/dsgetdc:", d.Domain)
output, err := cmd.Output()
if err != nil {
return nil, err
}

lines := strings.Split(string(output), "\n")
for _, line := range lines {
if strings.Contains(line, "Domain") {
info["domain"] = strings.TrimSpace(strings.Split(line, ":")[1])
}
if strings.Contains(line, "Dns Forest") {
info["forest"] = strings.TrimSpace(strings.Split(line, ":")[1])
}
}

return info, nil
}

func (d *DomainEnum) GetDomainControllers() ([]string, error) {
cmd := exec.Command("nltest", "/dclist:", d.Domain)
output, err := cmd.Output()
if err != nil {
return nil, err
}

var dcs []string
lines := strings.Split(string(output), "\n")
for _, line := range lines {
if strings.Contains(line, "\\") {
dc := strings.TrimSpace(strings.Split(line, "\\")[1])
dcs = append(dcs, dc)
}
}
return dcs, nil
}
