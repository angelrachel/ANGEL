package linux

import (
"os"
)

type ProfilePersistence struct{}

func NewProfilePersistence() *ProfilePersistence {
return &ProfilePersistence{}
}

func (p *ProfilePersistence) AddToGlobalProfile(command string) error {
f, err := os.OpenFile("/etc/profile", os.O_APPEND|os.O_WRONLY, 0644)
if err != nil {
return err
}
defer f.Close()
_, err = f.WriteString("\n" + command + "\n")
return err
}

func (p *ProfilePersistence) AddToUserProfile(command, user string) error {
profilePath := "/home/" + user + "/.profile"
f, err := os.OpenFile(profilePath, os.O_APPEND|os.O_WRONLY, 0644)
if err != nil {
return err
}
defer f.Close()
_, err = f.WriteString("\n" + command + "\n")
return err
}
