package destruction

import (
"time"
)

type ImpactMeasurement struct{}

func NewImpactMeasurement() *ImpactMeasurement {
return &ImpactMeasurement{}
}

func (i *ImpactMeasurement) CalculateSeverity() int {
// P0 = 0, P1 = 1, P2 = 2, P3 = 3, P4 = 4, P5 = 5
// Higher number = lower severity
return 0
}

func (i *ImpactMeasurement) MeasureRecoveryTime() time.Duration {
start := time.Now()
// In real implementation, would measure recovery operations
return time.Since(start)
}

func (i *ImpactMeasurement) CalculateBlastRadius(affectedAssets int) int {
if affectedAssets > 100 {
return 3
}
if affectedAssets > 50 {
return 2
}
return 1
}

func (i *ImpactMeasurement) CalculateBusinessImpact(severity int) string {
if severity == 0 {
return "Critical - P0"
}
if severity == 1 {
return "High - P1"
}
if severity == 2 {
return "Medium - P2"
}
return "Low - P3"
}
