package brain

import (
"time"
)

type TimingControl struct {
BaseSleep    int
MaxSleep     int
CurrentSleep int
Jitter       int
}

func NewTimingControl(baseSleep int) *TimingControl {
return &TimingControl{
BaseSleep:    baseSleep,
MaxSleep:     baseSleep * 3,
CurrentSleep: baseSleep,
Jitter:       2,
}
}

func (t *TimingControl) CalculateSleep() int {
sleep := t.CurrentSleep
// Add jitter
sleep += t.Jitter * (time.Now().Nanosecond() % 3)
if sleep > t.MaxSleep {
sleep = t.MaxSleep
}
return sleep
}

func (t *TimingControl) Adjust(environment string) {
switch environment {
case "suspicious":
t.CurrentSleep = t.CurrentSleep + 5
case "safe":
t.CurrentSleep = t.BaseSleep
case "sandbox":
t.CurrentSleep = t.MaxSleep
}
}

func (t *TimingControl) Sleep() {
sleep := t.CalculateSleep()
time.Sleep(time.Duration(sleep) * time.Second)
}

func (t *TimingControl) AdaptiveSleep(indicators []string) {
suspiciousCount := 0
for _, indicator := range indicators {
if indicator == "edr" || indicator == "sandbox" {
suspiciousCount++
}
}
if suspiciousCount > 0 {
t.Adjust("suspicious")
} else {
t.Adjust("safe")
}
}
