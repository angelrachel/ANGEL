package linux

import (
"os"
)

type BashrcPersistence struct{}

func NewBashrcPersistence() *BashrcPersistence {
return &BashrcPersistence{}
}

func (b *BashrcPersistence) AddToBashrc(command, user string) error {
bashrcPath := "/home/" + user + "/.bashrc"
if user == "root" {
bashrcPath = "/root/.bashrc"
}
f, err := os.OpenFile(bashrcPath, os.O_APPEND|os.O_WRONLY, 0644)
if err != nil {
return err
}
defer f.Close()
_, err = f.WriteString("\n" + command + "\n")
return err
}
