package destruct_impact

import (
	"fmt"
	"time"
)

type ImpactCalculator struct {
	Targets       []string
	DataSize      int64
	Downtime      time.Duration
	AffectedUsers int
}

func NewImpactCalculator() *ImpactCalculator {
	return &ImpactCalculator{}
}

func (i *ImpactCalculator) CalculateDataImpact(dataSize int64) string {
	if dataSize > 1000000 {
		return "CRITICAL: Large data breach"
	} else if dataSize > 100000 {
		return "HIGH: Significant data breach"
	} else {
		return "MEDIUM: Minor data breach"
	}
}

func (i *ImpactCalculator) CalculateDowntimeImpact(duration time.Duration) string {
	if duration > 24*time.Hour {
		return "CRITICAL: Extended downtime"
	} else if duration > 4*time.Hour {
		return "HIGH: Significant downtime"
	} else {
		return "MEDIUM: Short downtime"
	}
}

func (i *ImpactCalculator) CalculateUserImpact(users int) string {
	if users > 10000 {
		return "CRITICAL: Widespread user impact"
	} else if users > 1000 {
		return "HIGH: Major user impact"
	} else {
		return "MEDIUM: Limited user impact"
	}
}

func (i *ImpactCalculator) TotalImpact() string {
	return fmt.Sprintf("Data: %s\nDowntime: %s\nUsers: %s",
		i.CalculateDataImpact(i.DataSize),
		i.CalculateDowntimeImpact(i.Downtime),
		i.CalculateUserImpact(i.AffectedUsers))
}

func (i *ImpactCalculator) GetTargets() []string {
	return i.Targets
}

func (i *ImpactCalculator) AddTarget(target string) {
	i.Targets = append(i.Targets, target)
}
