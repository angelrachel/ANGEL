package orchestrator

import (
"context"
"time"
)

type AttackSurface struct {
Targets []string
Vectors []string
}

func NewAttackSurface() *AttackSurface {
return &AttackSurface{
Targets: []string{},
Vectors: []string{},
}
}

func (a *AttackSurface) AddTarget(target string) {
a.Targets = append(a.Targets, target)
}

func (a *AttackSurface) AddVector(vector string) {
a.Vectors = append(a.Vectors, vector)
}

func (a *AttackSurface) GetSurface() []string {
return a.Targets
}

func (a *AttackSurface) GetVectors() []string {
return a.Vectors
}

func (a *AttackSurface) Scan(ctx context.Context) bool {
return true
}

func (a *AttackSurface) GetLastSeen() time.Time {
return time.Now()
}
