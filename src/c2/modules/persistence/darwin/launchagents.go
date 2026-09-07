package darwin

import (
"os"
)

type LaunchAgentsPersistence struct{}

func NewLaunchAgentsPersistence() *LaunchAgentsPersistence {
return &LaunchAgentsPersistence{}
}

func (l *LaunchAgentsPersistence) CreateAgent(name, execPath string) error {
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
return os.WriteFile("/Library/LaunchAgents/"+name+".plist", []byte(plistContent), 0644)
}

func (l *LaunchAgentsPersistence) LoadAgent(name string) error {
return os.Symlink("/Library/LaunchAgents/"+name+".plist", "/Library/LaunchAgents/"+name+".plist")
}
