//go:build windows

package ad_recon

import (
"os/exec"
)

func EnumerateSites() string {
cmd := exec.Command("cmd", "/c", "dsquery site -name *")
output, err := cmd.Output()
if err != nil {
return ""
}
return string(output)
}

func EnumerateSiteSubnets() string {
cmd := exec.Command("cmd", "/c", "dsquery subnet -name *")
output, err := cmd.Output()
if err != nil {
return ""
}
return string(output)
}

func EnumerateSiteServers(site string) string {
cmd := exec.Command("cmd", "/c", "dsquery server -site "+site)
output, err := cmd.Output()
if err != nil {
return ""
}
return string(output)
}

func EnumerateAllSites() string {
return EnumerateSites()
}
