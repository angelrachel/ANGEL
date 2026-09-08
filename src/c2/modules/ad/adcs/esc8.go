package adcs

import "os/exec"

type ESC8 struct {
Target string
}

func (e ESC8) RequestCertificate() ESCResult {
cmd := exec.Command("certipy", "req", "-target", e.Target)
out, _ := cmd.Output()
return ESCResult{Status: "success", Output: string(out)}
}
