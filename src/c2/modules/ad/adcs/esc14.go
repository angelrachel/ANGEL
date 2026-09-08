package adcs

import "os/exec"

type ESC14 struct {
Target string
}

func (e ESC14) RequestCertificate() ESCResult {
cmd := exec.Command("certipy", "req", "-target", e.Target)
out, _ := cmd.Output()
return ESCResult{Status: "success", Output: string(out)}
}
