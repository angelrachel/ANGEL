//go:build windows

package ad_recon

import (
"os/exec"
)

func EnumerateTrusts() string {
cmd := exec.Command("cmd", "/c", "nltest /domain_trusts")
output, err := cmd.Output()
if err != nil {
return ""
}
return string(output)
}

func EnumerateDomainTrusts() string {
cmd := exec.Command("cmd", "/c", "netdom query trust")
output, err := cmd.Output()
if err != nil {
return ""
}
return string(output)
}

func EnumerateForestTrusts() string {
cmd := exec.Command("cmd", "/c", "nltest /dsgetdc:forest")
output, err := cmd.Output()
if err != nil {
return ""
}
return string(output)
}

func DumpAllTrusts() string {
var result string
result += EnumerateTrusts()
result += EnumerateDomainTrusts()
result += EnumerateForestTrusts()
return result
}
