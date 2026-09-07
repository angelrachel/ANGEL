package ad_recon

import (
"os/exec"
"strings"
)

type SPNEnum struct {
Domain string
}

func NewSPNEnum(domain string) *SPNEnum {
return &SPNEnum{
Domain: domain,
}
}

func (s *SPNEnum) GetSPNs() ([]string, error) {
cmd := exec.Command("setspn", "-T", s.Domain, "-F", "-Q", "*/*")
output, err := cmd.Output()
if err != nil {
return nil, err
}

var spns []string
lines := strings.Split(string(output), "\n")
for _, line := range lines {
if strings.Contains(line, "CN=") {
parts := strings.Fields(line)
for _, part := range parts {
if strings.Contains(part, "/") {
spns = append(spns, part)
}
}
}
}
return spns, nil
}

func (s *SPNEnum) GetUserSPNs() ([]string, error) {
cmd := exec.Command("setspn", "-T", s.Domain, "-F", "-Q", "*/")
output, err := cmd.Output()
if err != nil {
return nil, err
}

var userSPNs []string
lines := strings.Split(string(output), "\n")
for _, line := range lines {
if strings.Contains(line, "CN=") && !strings.Contains(line, "krbtgt") {
parts := strings.Fields(line)
for _, part := range parts {
if strings.Contains(part, "/") {
userSPNs = append(userSPNs, part)
}
}
}
}
return userSPNs, nil
}
