//go:build windows

package network_evasion

import (
	"crypto/rand"
	"fmt"
	"time"
)

func GenerateRandomIP() string {
	b := make([]byte, 4)
	rand.Read(b)
	return fmt.Sprintf("%d.%d.%d.%d", b[0], b[1], b[2], b[3])
}

func RotateIP() string {
	return GenerateRandomIP()
}

func GenerateJitter(min, max time.Duration) time.Duration {
	diff := max - min
	return min + time.Duration(time.Now().UnixNano()%int64(diff))
}

func SleepWithJitter(min, max time.Duration) {
	time.Sleep(GenerateJitter(min, max))
}
