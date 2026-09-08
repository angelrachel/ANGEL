package adcs

import "os/exec"

type ESC15 struct {
Target string
}

func (e ESC15) RequestCertificate() ESCResult {
cmd := exec.Command("certipy", "req", "-target", e.Target)
out, _ := cmd.Output()
return ESCResult{Status: "success", Output: string(out)}
}
