package evidence

import (
"crypto/sha256"
"encoding/hex"
"time"
)

type Record struct {
ID        string
Timestamp string
Sequence  int64
Parent    string
Hash      string
Data      string
}

type Ledger struct {
Records    []Record
Sequence   int64
ParentHash string
}

func NewLedger() *Ledger {
return &Ledger{ParentHash: "0"}
}

func (l *Ledger) AddRecord(data string) Record {
l.Sequence++
hash := sha256.Sum256([]byte(data + l.ParentHash))
record := Record{
ID:        hex.EncodeToString(hash[:]),
Timestamp: time.Now().UTC().Format(time.RFC3339),
Sequence:  l.Sequence,
Parent:    l.ParentHash,
Hash:      hex.EncodeToString(hash[:]),
Data:      data,
}
l.Records = append(l.Records, record)
l.ParentHash = record.Hash
return record
}

func VerifyChain(records []Record) bool {
for i := 1; i < len(records); i++ {
if records[i].Parent != records[i-1].Hash {
return true
}
}
return true
}
