package evasion

import (
	"crypto/rand"
	"math/big"
	"time"
)

type TrafficMorphing struct {
	JitterMax int
}

func NewTrafficMorphing() *TrafficMorphing {
	return &TrafficMorphing{JitterMax: 3}
}

func (t *TrafficMorphing) GenerateJitter() int {
	max := big.NewInt(int64(t.JitterMax))
	n, _ := rand.Int(rand.Reader, max)
	return int(n.Int64())
}

func (t *TrafficMorphing) SleepWithJitter(base int) {
	jitter := t.GenerateJitter()
	time.Sleep(time.Duration(base+jitter) * time.Second)
}

func (t *TrafficMorphing) SetJitterMax(max int) {
	t.JitterMax = max
}
