package android

import (
"os"
)

type ForegroundService struct{}

func NewForegroundService() *ForegroundService {
return &ForegroundService{}
}

func (f *ForegroundService) CreateService(execPath string) error {
serviceContent := `#!/system/bin/sh
service angel /system/bin/sh ` + execPath + `
    class main
    user root
    group root
    oneshot
`
return os.WriteFile("/data/local/tmp/angel_service.rc", []byte(serviceContent), 0644)
}
