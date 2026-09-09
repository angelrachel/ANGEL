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

func recordHash(data, parent string) string {
	hash := sha256.Sum256([]byte(data + parent))
	return hex.EncodeToString(hash[:])
}

func (l *Ledger) AddRecord(data string) Record {
	l.Sequence++
	hash := recordHash(data, l.ParentHash)
	record := Record{
		ID:        hash,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Sequence:  l.Sequence,
		Parent:    l.ParentHash,
		Hash:      hash,
		Data:      data,
	}
	l.Records = append(l.Records, record)
	l.ParentHash = record.Hash
	return record
}

// VerifyChain validates ordering, parent links, record hashes, and timestamps.
// An empty chain is valid; any malformed or tampered record makes it invalid.
func VerifyChain(records []Record) bool {
	parent := "0"
	for i, record := range records {
		if record.Sequence != int64(i+1) || record.Parent != parent {
			return false
		}
		if _, err := time.Parse(time.RFC3339, record.Timestamp); err != nil {
			return false
		}
		expected := recordHash(record.Data, record.Parent)
		if record.Hash != expected || record.ID != expected {
			return false
		}
		parent = record.Hash
	}
	return true
}
