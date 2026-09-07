package cred_crack

import (
"os/exec"
"strconv"
"strings"
)

type JohnWrapper struct {
BinaryPath string
Wordlist   string
Rules      string
}

func NewJohnWrapper(binaryPath, wordlist, rules string) *JohnWrapper {
return &JohnWrapper{
BinaryPath: binaryPath,
Wordlist:   wordlist,
Rules:      rules,
}
}

func (j *JohnWrapper) Crack(hashFile string, format string) (bool, error) {
var args []string
args = append(args, "--wordlist="+j.Wordlist)
args = append(args, "--format="+format)
if j.Rules != "" {
args = append(args, "--rules="+j.Rules)
}
args = append(args, hashFile)

cmd := exec.Command(j.BinaryPath, args...)
err := cmd.Run()
if err != nil {
return false, err
}
return true, nil
}

func (j *JohnWrapper) Show(hashFile string, format string) (string, error) {
args := []string{
"--show",
"--format=" + format,
hashFile,
}
cmd := exec.Command(j.BinaryPath, args...)
output, err := cmd.Output()
if err != nil {
return "", err
}
return string(output), nil
}

func (j *JohnWrapper) GetCrackedHashes(hashFile string, format string) ([]string, error) {
output, err := j.Show(hashFile, format)
if err != nil {
return nil, err
}
var results []string
lines := strings.Split(string(output), "\n")
for _, line := range lines {
if strings.Contains(line, ":") {
results = append(results, line)
}
}
return results, nil
}

func (j *JohnWrapper) TestConnection() bool {
cmd := exec.Command(j.BinaryPath, "--help")
err := cmd.Run()
return err == nil
}
