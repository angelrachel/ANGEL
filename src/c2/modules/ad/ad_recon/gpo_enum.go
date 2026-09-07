//go:build windows

package ad_recon

import (
"os/exec"
)

func EnumerateGPOs() string {
cmd := exec.Command("cmd", "/c", "gpresult /r")
output, err := cmd.Output()
if err != nil {
return ""
}
return string(output)
}

func EnumerateGPOReport() string {
cmd := exec.Command("cmd", "/c", "gpresult /h report.html")
output, err := cmd.Output()
if err != nil {
return ""
}
return string(output)
}

func EnumerateGPOFromDC(domain string) string {
cmd := exec.Command("cmd", "/c", "Get-GPO -All -Domain "+domain)
output, err := cmd.Output()
if err != nil {
return ""
}
return string(output)
}

func DumpAllGPOs() string {
var result string
result += EnumerateGPOs()
result += EnumerateGPOReport()
return result
}
