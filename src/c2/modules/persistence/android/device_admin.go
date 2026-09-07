package android

import (
"os"
)

type DeviceAdmin struct{}

func NewDeviceAdmin() *DeviceAdmin {
return &DeviceAdmin{}
}

func (d *DeviceAdmin) EnableDeviceAdmin() error {
adminContent := `<?xml version="1.0" encoding="utf-8"?>
<device-admin xmlns:android="http://schemas.android.com/apk/res/android">
    <uses-policies>
        <limit-password />
        <watch-login />
        <reset-password />
        <force-lock />
        <wipe-data />
    </uses-policies>
</device-admin>`
return os.WriteFile("/data/local/tmp/device_admin.xml", []byte(adminContent), 0644)
}
