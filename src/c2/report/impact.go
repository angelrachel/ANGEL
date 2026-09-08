package report

import "time"

type Impact struct {
Severity        string
BusinessImpact  string
FinancialImpact string
}

func NewImpact() *Impact {
return &Impact{}
}

func (i *Impact) SetSeverity(severity string) {
i.Severity = severity
}

func (i *Impact) SetBusinessImpact(businessImpact string) {
i.BusinessImpact = businessImpact
}

func (i *Impact) SetFinancialImpact(financialImpact string) {
i.FinancialImpact = financialImpact
}

func (i *Impact) GetSeverity() string {
return i.Severity
}

func (i *Impact) GetBusinessImpact() string {
return i.BusinessImpact
}

func (i *Impact) GetFinancialImpact() string {
return i.FinancialImpact
}

func (i *Impact) GetLastSeen() time.Time {
return time.Now()
}
