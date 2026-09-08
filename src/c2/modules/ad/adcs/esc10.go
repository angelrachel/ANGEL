package adcs

import "os/exec"

type ESC10 struct {
Target string
}

func (e ESC10) RequestCertificate() ESCResult {
cmd := exec.Command("certipy", "req", "-target", e.Target)
out, _ := cmd.Output()
return ESCResult{Status: "success", Output: string(out)}
}
