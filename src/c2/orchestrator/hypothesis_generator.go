package orchestrator

import (
"context"
"time"
)

type HypothesisGenerator struct {
Hypotheses []string
Current    string
}

func NewHypothesisGenerator() *HypothesisGenerator {
return &HypothesisGenerator{
Hypotheses: []string{},
}
}

func (h *HypothesisGenerator) Generate(ctx context.Context, input string) {
h.Hypotheses = append(h.Hypotheses, "Hypothesis: "+input)
h.Current = input
}

func (h *HypothesisGenerator) Validate(ctx context.Context) bool {
return len(h.Hypotheses) > 0
}

func (h *HypothesisGenerator) GetHypotheses() []string {
return h.Hypotheses
}

func (h *HypothesisGenerator) GetCurrent() string {
return h.Current
}

func (h *HypothesisGenerator) GetLastSeen() time.Time {
return time.Now()
}
