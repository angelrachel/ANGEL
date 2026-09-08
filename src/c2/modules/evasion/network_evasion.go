package evasion

import (
"crypto/rand"
"math/big"
)

type NetworkEvasion struct {
JitterMax int
}

func NewNetworkEvasion() *NetworkEvasion {
return &NetworkEvasion{JitterMax: 3}
}

func (n *NetworkEvasion) GenerateJitter() int {
max := big.NewInt(int64(n.JitterMax))
val, _ := rand.Int(rand.Reader, max)
return int(val.Int64())
}

func (n *NetworkEvasion) SetJitterMax(max int) {
n.JitterMax = max
}
