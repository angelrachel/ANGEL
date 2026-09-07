package kerberos

import (
"fmt"
)

type ASREPRoast struct {
Domain string
}

func NewASREPRoast(domain string) *ASREPRoast {
return &ASREPRoast{
Domain: domain,
}
}

func (a *ASREPRoast) RequestASREP(username string) (string, error) {
// Request AS-REP for users with DONT_REQ_PREAUTH
// No credentials needed
return fmt.Sprintf("$krb5asrep$%s$%s$%s", a.Domain, username, "HASH"), nil
}

func (a *ASREPRoast) CheckPreAuth(username string) bool {
// Check if user has DONT_REQ_PREAUTH flag set
// In real implementation, would query LDAP
return true
}
