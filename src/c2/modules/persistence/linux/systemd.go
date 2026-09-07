package linux

import (
"os"
)

type SystemdPersistence struct{}

func NewSystemdPersistence() *SystemdPersistence {
return &SystemdPersistence{}
}

func (s *SystemdPersistence) CreateService(name, execStart string) error {
serviceContent := `[Unit]
Description=Angel Service
After=network.target

[Service]
ExecStart=` + execStart + `
Restart=always
User=root

[Install]
WantedBy=multi-user.target`
return os.WriteFile("/etc/systemd/system/"+name+".service", []byte(serviceContent), 0644)
}

func (s *SystemdPersistence) EnableService(name string) error {
return os.Symlink("/etc/systemd/system/"+name+".service", "/etc/systemd/system/multi-user.target.wants/"+name+".service")
}

func (s *SystemdPersistence) CreateTimer(name, command, schedule string) error {
timerContent := `[Unit]
Description=Angel Timer

[Timer]
OnCalendar=` + schedule + `
Persistent=true

[Install]
WantedBy=timers.target`
return os.WriteFile("/etc/systemd/system/"+name+".timer", []byte(timerContent), 0644)
}
