package brain

import "time"

type TimingConfig struct {
	BaseSleep int
	MaxSleep  int
}

type TimingControl struct {
	Config TimingConfig
}

func NewTimingControl(baseSleep, maxSleep int) *TimingControl {
	return &TimingControl{Config: TimingConfig{BaseSleep: baseSleep, MaxSleep: maxSleep}}
}

func (t *TimingControl) Sleep(seconds int) {
	if seconds < 1 {
		seconds = t.Config.BaseSleep
	}
	time.Sleep(time.Duration(seconds) * time.Second)
}

func (t *TimingControl) Jitter() int {
	return t.Config.BaseSleep + t.Config.MaxSleep/2
}

func (t *TimingControl) Adjust(environment string) int {
	if environment == "suspicious" {
		return t.Config.MaxSleep
	}
	return t.Config.BaseSleep
}
