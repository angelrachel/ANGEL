package ad_recon

import (
"os/exec"
"strings"
)

type GroupEnum struct {
Domain string
}

func NewGroupEnum(domain string) *GroupEnum {
return &GroupEnum{
Domain: domain,
}
}

func (g *GroupEnum) GetGroups() ([]string, error) {
cmd := exec.Command("net", "group", "/domain")
output, err := cmd.Output()
if err != nil {
return nil, err
}

var groups []string
lines := strings.Split(string(output), "\n")
inGroups := false
for _, line := range lines {
if strings.Contains(line, "---") {
inGroups = true
continue
}
if inGroups && strings.TrimSpace(line) != "" && !strings.Contains(line, "command completed") {
groups = append(groups, strings.TrimSpace(line))
}
}
return groups, nil
}

func (g *GroupEnum) GetGroupMembers(group string) ([]string, error) {
cmd := exec.Command("net", "group", group, "/domain")
output, err := cmd.Output()
if err != nil {
return nil, err
}

var members []string
lines := strings.Split(string(output), "\n")
inMembers := false
for _, line := range lines {
if strings.Contains(line, "---") {
inMembers = true
continue
}
if inMembers && strings.TrimSpace(line) != "" && !strings.Contains(line, "command completed") {
members = append(members, strings.TrimSpace(line))
}
}
return members, nil
}
