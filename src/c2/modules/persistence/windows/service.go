package windows

import (
"os/exec"
)

type ServicePersistence struct{}

func NewServicePersistence() *ServicePersistence {
return &ServicePersistence{}
}

func (s *ServicePersistence) CreateService(name, path string) error {
cmd := exec.Command("sc", "create", name, "binPath="+path, "start=auto")
return cmd.Run()
}

func (s *ServicePersistence) StartService(name string) error {
cmd := exec.Command("sc", "start", name)
return cmd.Run()
}

func (s *ServicePersistence) DeleteService(name string) error {
cmd := exec.Command("sc", "delete", name)
return cmd.Run()
}
