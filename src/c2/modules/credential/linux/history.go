package linux

import (
"os"
)

type History struct{}

func NewHistory() *History {
return &History{}
}

func (h *History) ExtractBashHistory(user string) (string, error) {
path := "/home/" + user + "/.bash_history"
if user == "root" {
path = "/root/.bash_history"
}
data, err := os.ReadFile(path)
if err != nil {
return "", err
}
return string(data), nil
}

func (h *History) ExtractZshHistory(user string) (string, error) {
path := "/home/" + user + "/.zsh_history"
if user == "root" {
path = "/root/.zsh_history"
}
data, err := os.ReadFile(path)
if err != nil {
return "", err
}
return string(data), nil
}
