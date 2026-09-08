package implant

import (
"crypto/rand"
"crypto/rsa"
"crypto/x509"
"encoding/pem"
"os"
)

type GeneratorError struct {
Message string
}

func (e GeneratorError) Error() string {
return e.Message
}

type Generator struct {
Key *rsa.PrivateKey
}

func NewGenerator() *Generator {
key, err := rsa.GenerateKey(rand.Reader, 2048)
if err != nil {
panic(err)
}
return &Generator{Key: key}
}

func (g *Generator) GeneratePrivateKey() ([]byte, error) {
return x509.MarshalPKCS1PrivateKey(g.Key), nil
}

func (g *Generator) GeneratePublicKey() ([]byte, error) {
return x509.MarshalPKCS1PublicKey(&g.Key.PublicKey), nil
}

func (g *Generator) GenerateConfig(server, agentID string) (map[string]string, error) {
if server == "" {
return map[string]string{}, GeneratorError{Message: "server required"}
}
if agentID == "" {
return map[string]string{}, GeneratorError{Message: "agent_id required"}
}
return map[string]string{"server": server, "agent_id": agentID, "key": "ANGEL_C2_MASTER_KEY_SL9X_2026"}, nil
}

func (g *Generator) WritePrivateKey(filename string) error {
key, err := g.GeneratePrivateKey()
if err != nil {
return err
}
block := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: key}
return os.WriteFile(filename, pem.EncodeToMemory(block), 0600)
}
