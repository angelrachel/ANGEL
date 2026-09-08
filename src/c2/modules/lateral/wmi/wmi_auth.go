package wmi

import (
"fmt"
"os/exec"
"strings"
)

type Result struct {
Command string
Output  string
Error   string
}

type Auth struct {
Host     string
User     string
Pass     string
}

func (a Auth) ExecuteQuery(query string) Result {
cmd := exec.Command("cmd", "/c", "wmic", "/node:"+a.Host, "/user:"+a.User, "/password:"+a.Pass, "process", "list", "brief")
out, err := cmd.Output()
result := Result{Command: query}
if err != nil {
result.Error = err.Error()
return result
}
result.Output = string(out)
return result
}

func (a Auth) GetSystemInfo() Result {
cmd := exec.Command("cmd", "/c", "wmic", "/node:"+a.Host, "/user:"+a.User, "/password:"+a.Pass, "os", "get", "caption,version")
out, err := cmd.Output()
result := Result{Command: "os"}
if err != nil {
result.Error = err.Error()
return result
}
result.Output = string(out)
return result
}

func (a Auth) EnumerateProcesses() Result {
cmd := exec.Command("cmd", "/c", "wmic", "/node:"+a.Host, "/user:"+a.User, "/password:"+a.Pass, "process", "get", "name,processid")
out, err := cmd.Output()
result := Result{Command: "process"}
if err != nil {
result.Error = err.Error()
return result
}
result.Output = string(out)
return result
}

func FormatOutput(output string) string {
return strings.TrimSpace(output)
}

func FormatAddress(host string, port int) string {
return fmt.Sprintf("%s:%d", host, port)
}
