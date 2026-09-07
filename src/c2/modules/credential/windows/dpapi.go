package windows

import (
"os"
)

type DPAPI struct{}

func NewDPAPI() *DPAPI {
return &DPAPI{}
}

func (d *DPAPI) Decrypt(path string) ([]byte, error) {
return os.ReadFile(path)
}
