package adcs

import "os/exec"

type ESCResult struct {
Status string
Output string
}

type ESC1 struct {
Target   string
Template string
}

func (e ESC1) RequestCertificate() ESCResult {
cmd := exec.Command("certipy", "req", "-target", e.Target, "-template", e.Template)
out, _ := cmd.Output()
return ESCResult{Status: "success", Output: string(out)}
}
