package ad_recon

import (
"os/exec"
"strings"
)

type UserEnum struct {
Domain string
}

func NewUserEnum(domain string) *UserEnum {
return &UserEnum{
Domain: domain,
}
}

func (u *UserEnum) GetUsers() ([]string, error) {
cmd := exec.Command("net", "user", "/domain")
output, err := cmd.Output()
if err != nil {
return nil, err
}

var users []string
lines := strings.Split(string(output), "\n")
inUsers := false
for _, line := range lines {
if strings.Contains(line, "---") {
inUsers = true
continue
}
if inUsers && strings.TrimSpace(line) != "" && !strings.Contains(line, "command completed") {
fields := strings.Fields(line)
for _, f := range fields {
if len(f) > 0 {
users = append(users, f)
}
}
}
}
return users, nil
}

func (u *UserEnum) GetAdminUsers() ([]string, error) {
cmd := exec.Command("net", "group", "Domain Admins", "/domain")
output, err := cmd.Output()
if err != nil {
return nil, err
}

var admins []string
lines := strings.Split(string(output), "\n")
inMembers := false
for _, line := range lines {
if strings.Contains(line, "---") {
inMembers = true
continue
}
if inMembers && strings.TrimSpace(line) != "" && !strings.Contains(line, "command completed") {
admins = append(admins, strings.TrimSpace(line))
}
}
return admins, nil
}
