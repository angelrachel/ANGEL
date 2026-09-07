package cred_crack

import (
"os/exec"
"strconv"
"strings"
)

type HashcatWrapper struct {
BinaryPath string
Wordlist   string
Rules      string
}

func NewHashcatWrapper(binaryPath, wordlist, rules string) *HashcatWrapper {
return &HashcatWrapper{
BinaryPath: binaryPath,
Wordlist:   wordlist,
Rules:      rules,
}
}

func (h *HashcatWrapper) SetBinaryPath(path string) {
h.BinaryPath = path
}

func (h *HashcatWrapper) SetWordlist(path string) {
h.Wordlist = path
}

func (h *HashcatWrapper) SetRules(path string) {
h.Rules = path
}

func (h *HashcatWrapper) Crack(hashFile, mode int, attackMode int) (bool, error) {
var args []string
args = append(args, "-m", strconv.Itoa(mode))
args = append(args, "-a", strconv.Itoa(attackMode))
args = append(args, hashFile)
args = append(args, h.Wordlist)

if h.Rules != "" {
args = append(args, "-r", h.Rules)
}

args = append(args, "--force", "-O", "--status")

cmd := exec.Command(h.BinaryPath, args...)
err := cmd.Run()
if err != nil {
return false, err
}
return true, nil
}

func (h *HashcatWrapper) CrackWithBruteForce(hashFile string, mode int, mask string) (bool, error) {
args := []string{
"-m", strconv.Itoa(mode),
"-a", "3",
hashFile,
mask,
"--force",
}

cmd := exec.Command(h.BinaryPath, args...)
err := cmd.Run()
if err != nil {
return false, err
}
return true, nil
}

func (h *HashcatWrapper) GetStatus(hashFile string) (string, error) {
args := []string{
"--show",
hashFile,
}
cmd := exec.Command(h.BinaryPath, args...)
output, err := cmd.Output()
if err != nil {
return "", err
}
return string(output), nil
}

func (h *HashcatWrapper) Cleanup() bool {
args := []string{
"--potfile-disable",
}
cmd := exec.Command(h.BinaryPath, args...)
err := cmd.Run()
if err != nil {
return false
}
return true
}

func (h *HashcatWrapper) GetCrackedHashes(hashFile string) ([]string, error) {
output, err := h.GetStatus(hashFile)
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

func (h *HashcatWrapper) TestConnection() bool {
cmd := exec.Command(h.BinaryPath, "-I")
err := cmd.Run()
return err == nil
}
