package windows

import (
"os"
)

type ADSPersistence struct{}

func NewADSPersistence() *ADSPersistence {
return &ADSPersistence{}
}

func (a *ADSPersistence) WriteADS(file, data string) error {
adsPath := file + ":angel"
return os.WriteFile(adsPath, []byte(data), 0644)
}

func (a *ADSPersistence) ReadADS(file string) (string, error) {
adsPath := file + ":angel"
data, err := os.ReadFile(adsPath)
if err != nil {
return "", err
}
return string(data), nil
}
