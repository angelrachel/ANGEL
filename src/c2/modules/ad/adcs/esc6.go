package adcs

import "os/exec"

type ESC6 struct {
Target string
}

func (e ESC6) RequestCertificate() ESCResult {
cmd := exec.Command("certipy", "req", "-target", e.Target)
out, _ := cmd.Output()
return ESCResult{Status: "success", Output: string(out)}
}
