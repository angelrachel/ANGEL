package linux

import (
"os"
)

type Shadow struct{}

func NewShadow() *Shadow {
return &Shadow{}
}

func (s *Shadow) ExtractShadow() (string, error) {
data, err := os.ReadFile("/etc/shadow")
if err != nil {
return "", err
}
return string(data), nil
}

func (s *Shadow) ExtractPasswd() (string, error) {
data, err := os.ReadFile("/etc/passwd")
if err != nil {
return "", err
}
return string(data), nil
}
