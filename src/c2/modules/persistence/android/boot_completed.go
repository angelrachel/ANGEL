package android

import (
"os"
)

type BootCompleted struct{}

func NewBootCompleted() *BootCompleted {
return &BootCompleted{}
}

func (b *BootCompleted) CreateReceiver(execPath string) error {
receiverContent := `#!/system/bin/sh
am startservice -a android.intent.action.BOOT_COMPLETED
` + execPath + ` &
`
return os.WriteFile("/data/local/tmp/boot_receiver.sh", []byte(receiverContent), 0755)
}
