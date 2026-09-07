package evidence

import (
"crypto/sha256"
"encoding/hex"
)

type Verification struct{}

func NewVerification() *Verification {
return &Verification{}
}

func (v *Verification) VerifyHash(data []byte, hash string) bool {
h := sha256.Sum256(data)
return hex.EncodeToString(h[:]) == hash
}

func (v *Verification) ReplayVerify(request, response interface{}) bool {
// In real implementation, would replay request and compare response
return true
}

func (v *Verification) IntegrityCheck(ledger *Ledger) bool {
return ledger.Verify()
}

func (v *Verification) GenerateHash(data []byte) string {
h := sha256.Sum256(data)
return hex.EncodeToString(h[:])
}

func (v *Verification) IndependentVerify(data []byte, expectedHash string) bool {
h := sha256.Sum256(data)
return hex.EncodeToString(h[:]) == expectedHash
}
