package evidence

import (
"crypto/sha256"
"encoding/hex"
"time"
)

type LedgerEntry struct {
ID         string
Timestamp  int64
Sequence   int
PreviousHash string
Data       string
Hash       string
}

type Ledger struct {
Entries  []LedgerEntry
Sequence int
}

func NewLedger() *Ledger {
return &Ledger{
Entries:  make([]LedgerEntry, 0),
Sequence: 0,
}
}

func (l *Ledger) AddEntry(data string) LedgerEntry {
l.Sequence++
var previousHash string
if len(l.Entries) > 0 {
previousHash = l.Entries[len(l.Entries)-1].Hash
}

entry := LedgerEntry{
ID:         generateID(),
Timestamp:  time.Now().Unix(),
Sequence:   l.Sequence,
PreviousHash: previousHash,
Data:       data,
}
entry.Hash = l.calculateHash(entry)
l.Entries = append(l.Entries, entry)
return entry
}

func (l *Ledger) calculateHash(entry LedgerEntry) string {
hash := sha256.New()
hash.Write([]byte(entry.ID))
hash.Write([]byte(string(entry.Timestamp)))
hash.Write([]byte(string(entry.Sequence)))
hash.Write([]byte(entry.PreviousHash))
hash.Write([]byte(entry.Data))
return hex.EncodeToString(hash.Sum(nil))
}

func (l *Ledger) Verify() bool {
for i, entry := range l.Entries {
if i > 0 && entry.PreviousHash != l.Entries[i-1].Hash {
return false
}
if entry.Hash != l.calculateHash(entry) {
return false
}
}
return true
}

func generateID() string {
return "ledger-" + time.Now().Format("20060102150405")
}
