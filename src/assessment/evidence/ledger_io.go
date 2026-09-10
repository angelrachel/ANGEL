package evidence

import (
	"encoding/json"
	"fmt"
)

func (l *Ledger) Export() ([]byte, error) {
	if l == nil || !VerifyChain(l.Records) {
		return nil, fmt.Errorf("cannot export invalid evidence ledger")
	}
	return json.MarshalIndent(struct {
		Records []Record `json:"records"`
	}{Records: l.Records}, "", "  ")
}

func Import(data []byte) (*Ledger, error) {
	var payload struct {
		Records []Record `json:"records"`
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		return nil, err
	}
	if !VerifyChain(payload.Records) {
		return nil, fmt.Errorf("invalid evidence ledger")
	}
	ledger := NewLedger()
	ledger.Records = append([]Record(nil), payload.Records...)
	ledger.Sequence = int64(len(ledger.Records))
	if len(ledger.Records) > 0 {
		ledger.ParentHash = ledger.Records[len(ledger.Records)-1].Hash
	}
	return ledger, nil
}

func (l *Ledger) Head() (Record, bool) {
	if l == nil || len(l.Records) == 0 {
		return Record{}, false
	}
	return l.Records[len(l.Records)-1], true
}
