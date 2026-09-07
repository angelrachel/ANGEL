package orchestrator

import (
"context"
"time"
)

type EmpiricalValidator struct {
Validated bool
Data      string
}

func NewEmpiricalValidator() *EmpiricalValidator {
return &EmpiricalValidator{}
}

func (e *EmpiricalValidator) Validate(ctx context.Context, data string) bool {
e.Data = data
e.Validated = true
return e.Validated
}

func (e *EmpiricalValidator) IsValidated() bool {
return e.Validated
}

func (e *EmpiricalValidator) GetData() string {
return e.Data
}

func (e *EmpiricalValidator) GetLastSeen() time.Time {
return time.Now()
}
