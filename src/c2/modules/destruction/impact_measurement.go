package destruction

import (
"time"
)

type ImpactMeasurement struct {
RecoveryTime time.Duration
AffectedData int64
}

func NewImpactMeasurement() *ImpactMeasurement {
return &ImpactMeasurement{}
}

func (i *ImpactMeasurement) MeasureRecoveryTime(start time.Time) time.Duration {
return time.Since(start)
}

func (i *ImpactMeasurement) CalculateBlastRadius(dataSize int64) string {
if dataSize > 1000000 {
return "CRITICAL"
}
if dataSize > 100000 {
return "HIGH"
}
return "MEDIUM"
}

func (i *ImpactMeasurement) CalculateBusinessImpact(affectedAssets int) string {
if affectedAssets > 100 {
return "CRITICAL"
}
if affectedAssets > 10 {
return "HIGH"
}
return "MEDIUM"
}
