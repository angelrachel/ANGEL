//go:build windows

package sleep_masking

import (
"time"
)

func SleepWithCallback(duration time.Duration, callback func()) {
if callback == nil {
return
}
time.Sleep(duration)
callback()
}

func SleepWithCallbackAndMask(duration time.Duration, callback func(), masker *SleepMasker) {
if callback == nil {
return
}
masker.EncryptRegion()
time.Sleep(duration)
masker.DecryptRegion()
callback()
}
