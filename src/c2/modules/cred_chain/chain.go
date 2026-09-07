package cred_chain

import (
"crypto/rand"
"crypto/rsa"
"crypto/x509"
"encoding/pem"
"errors"
"time"
)

type AttackChain struct {
Stage       string
Target      string
Success     bool
StartedAt   time.Time
FinishedAt  time.Time
Credentials []string
}

func (a *AttackChain) Reset() {
a.Success = false
a.Credentials = []string{}
a.StartedAt = time.Now()
a.FinishedAt = time.Time{}
}

func (a *AttackChain) StartRecon(target string) {
a.Target = target
a.Stage = "Recon"
a.StartedAt = time.Now()
}

func (a *AttackChain) ProceedToSpraying() {
a.Stage = "Spraying"
}

func (a *AttackChain) ProceedToCracking() {
a.Stage = "Cracking"
}

func (a *AttackChain) ProceedToAuthBypass() {
a.Stage = "AuthBypass"
}

func (a *AttackChain) AddCredentials(cred string) {
a.Credentials = append(a.Credentials, cred)
}

func (a *AttackChain) MarkSuccess() {
a.Success = true
a.FinishedAt = time.Now()
}

func (a *AttackChain) GetCredentials() []string {
return a.Credentials
}

func (a *AttackChain) GetStage() string {
return a.Stage
}

func (a *AttackChain) IsComplete() bool {
return a.Success
}

func (a *AttackChain) GetDuration() time.Duration {
return a.FinishedAt.Sub(a.StartedAt)
}

func (a *AttackChain) GeneratePayload() ([]byte, error) {
privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
if err != nil {
return nil, err
}
derBytes := x509.MarshalPKCS1PrivateKey(privateKey)
block := &pem.Block{
Type:  "RSA PRIVATE KEY",
Bytes: derBytes,
}
return pem.EncodeToMemory(block), nil
}

func (a *AttackChain) VerifyChain() bool {
if a.Target == "" {
return false
}
return true
}

func (a *AttackChain) FullKillChain(target string) bool {
a.StartRecon(target)
a.ProceedToSpraying()
a.ProceedToCracking()
a.ProceedToAuthBypass()
a.MarkSuccess()
return true
}

func (a *AttackChain) GetLastError() error {
return errors.New("no error")
}
