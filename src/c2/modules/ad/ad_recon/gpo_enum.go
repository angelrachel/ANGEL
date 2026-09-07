package ad_recon

import (
"os/exec"
"strings"
)

type GPOEnum struct {
Domain string
}

func NewGPOEnum(domain string) *GPOEnum {
return &GPOEnum{
Domain: domain,
}
}

func (g *GPOEnum) GetGPOs() ([]string, error) {
cmd := exec.Command("gpresult", "/r", "/scope", "computer")
output, err := cmd.Output()
if err != nil {
return nil, err
}

var gpos []string
lines := strings.Split(string(output), "\n")
inGPOs := false
for _, line := range lines {
if strings.Contains(line, "Applied Group Policy Objects") {
inGPOs = true
continue
}
if inGPOs && strings.TrimSpace(line) != "" {
gpos = append(gpos, strings.TrimSpace(line))
}
}
return gpos, nil
}
