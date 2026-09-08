package adcs

import "os/exec"

type ESC11 struct {
Target string
}

func (e ESC11) RequestCertificate() ESCResult {
cmd := exec.Command("certipy", "req", "-target", e.Target)
out, _ := cmd.Output()
return ESCResult{Status: "success", Output: string(out)}
}
