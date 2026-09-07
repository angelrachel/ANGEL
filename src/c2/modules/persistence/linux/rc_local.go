package linux

import (
"os"
)

type RCLocalPersistence struct{}

func NewRCLocalPersistence() *RCLocalPersistence {
return &RCLocalPersistence{}
}

func (r *RCLocalPersistence) AddToRCLocal(command string) error {
f, err := os.OpenFile("/etc/rc.local", os.O_APPEND|os.O_WRONLY, 0755)
if err != nil {
return err
}
defer f.Close()
_, err = f.WriteString("\n" + command + " &\n")
return err
}
