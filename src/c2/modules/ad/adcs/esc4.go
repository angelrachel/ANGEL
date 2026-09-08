package adcs

import "os/exec"

type ESC4 struct {
Target string
User   string
}

func (e ESC4) RequestCertificate() ESCResult {
cmd := exec.Command("certipy", "req", "-u", e.User, "-target", e.Target)
out, _ := cmd.Output()
return ESCResult{Status: "success", Output: string(out)}
}
