package ad_recon

import (
"os/exec"
"strings"
)

type SiteEnum struct {
Domain string
}

func NewSiteEnum(domain string) *SiteEnum {
return &SiteEnum{
Domain: domain,
}
}

func (s *SiteEnum) GetSites() ([]string, error) {
cmd := exec.Command("nltest", "/dsgetsite")
output, err := cmd.Output()
if err != nil {
return nil, err
}

var sites []string
lines := strings.Split(string(output), "\n")
for _, line := range lines {
if strings.Contains(line, "Site") {
parts := strings.Split(line, ":")
if len(parts) > 1 {
sites = append(sites, strings.TrimSpace(parts[1]))
}
}
}
return sites, nil
}
