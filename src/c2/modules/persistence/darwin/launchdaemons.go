package darwin

import (
"os"
)

type LaunchDaemonsPersistence struct{}

func NewLaunchDaemonsPersistence() *LaunchDaemonsPersistence {
return &LaunchDaemonsPersistence{}
}

func (l *LaunchDaemonsPersistence) CreateDaemon(name, execPath string) error {
plistContent := `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>` + name + `</string>
    <key>ProgramArguments</key>
    <array>
        <string>` + execPath + `</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
</dict>
</plist>`
return os.WriteFile("/Library/LaunchDaemons/"+name+".plist", []byte(plistContent), 0644)
}

func (l *LaunchDaemonsPersistence) LoadDaemon(name string) error {
return os.Symlink("/Library/LaunchDaemons/"+name+".plist", "/Library/LaunchDaemons/"+name+".plist")
}
