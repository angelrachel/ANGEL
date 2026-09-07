package ad_recon

import (
"os/exec"
"strings"
)

type OUEnum struct {
Domain string
}

func NewOUEnum(domain string) *OUEnum {
return &OUEnum{
Domain: domain,
}
}

func (o *OUEnum) GetOUs() ([]string, error) {
cmd := exec.Command("dsquery", "ou", "-domain", o.Domain)
output, err := cmd.Output()
if err != nil {
return nil, err
}

var ous []string
lines := strings.Split(string(output), "\n")
for _, line := range lines {
if strings.Contains(line, "OU=") {
ous = append(ous, strings.TrimSpace(line))
}
}
return ous, nil
}
