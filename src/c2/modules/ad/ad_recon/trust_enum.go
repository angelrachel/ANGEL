package ad_recon

import (
"os/exec"
"strings"
)

type TrustEnum struct {
Domain string
}

func NewTrustEnum(domain string) *TrustEnum {
return &TrustEnum{
Domain: domain,
}
}

func (t *TrustEnum) GetTrusts() ([]string, error) {
cmd := exec.Command("nltest", "/domain_trusts")
output, err := cmd.Output()
if err != nil {
return nil, err
}

var trusts []string
lines := strings.Split(string(output), "\n")
for _, line := range lines {
if strings.Contains(line, "Trusted") {
parts := strings.Split(line, ":")
if len(parts) > 1 {
trusts = append(trusts, strings.TrimSpace(parts[1]))
}
}
}
return trusts, nil
}
