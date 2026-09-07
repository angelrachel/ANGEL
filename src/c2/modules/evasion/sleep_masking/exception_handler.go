//go:build windows

package sleep_masking

import (
"time"
)

func SleepWithExceptionHandler(duration time.Duration) bool {
time.Sleep(duration)
return true
}

func SleepWithExceptionMask(duration time.Duration, masker *SleepMasker) bool {
masker.EncryptRegion()
time.Sleep(duration)
masker.DecryptRegion()
return true
}
