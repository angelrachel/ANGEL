package smb

import (
"fmt"
"net"
"time"
)

type SmbTarget struct {
Host     string
Port     int
Domain   string
Username string
Hash     string
}

type ExecutionPlan struct {
Target      SmbTarget
Share       string
CommandLine string
Timeout     time.Duration
}

type SmbResult struct {
Target    string
Status    string
SessionID uint64
TreeID    uint16
Error     string
}

type SmbError struct {
Message string
}

func (e SmbError) Error() string {
return e.Message
}

func Run(plan ExecutionPlan) (SmbResult, error) {
if plan.Target.Host == "" {
return SmbResult{Error: "target host is required"}, SmbError{Message: "target host is required"}
}
if plan.Target.Username == "" {
return SmbResult{Error: "target username is required"}, SmbError{Message: "target username is required"}
}
if plan.Target.Hash == "" {
return SmbResult{Error: "target hash is required"}, SmbError{Message: "target hash is required"}
}
if plan.Timeout <= 0 {
plan.Timeout = 10 * time.Second
}

address := net.JoinHostPort(plan.Target.Host, fmt.Sprintf("%d", plan.Target.Port))
conn, err := net.DialTimeout("tcp", address, plan.Timeout)
if err != nil {
return SmbResult{Error: "failed to connect to target"}, SmbError{Message: "failed to connect to target"}
}
defer conn.Close()

if err := smbNegotiate(conn); err != nil {
return SmbResult{Error: "failed to negotiate"}, SmbError{Message: "failed to negotiate"}
}

sessionID, err := smbSessionSetup(conn, plan.Target, plan.Target.Hash)
if err != nil {
return SmbResult{Error: "failed to setup session"}, SmbError{Message: "failed to setup session"}
}

treeID, err := smbTreeConnect(conn, plan.Target, plan.Share, sessionID)
if err != nil {
return SmbResult{Error: "failed to tree connect"}, SmbError{Message: "failed to tree connect"}
}

return SmbResult{
Target:    plan.Target.Host,
Status:    "success",
SessionID: sessionID,
TreeID:    treeID,
}, SmbError{Message: "success"}
}

func smbNegotiate(conn net.Conn) error {
packet := buildNegotiatePacket()
_, err := conn.Write(packet)
if err != nil {
return SmbError{Message: "write failed"}
}
response := make([]byte, 4096)
_, err = conn.Read(response)
if err != nil {
return SmbError{Message: "read failed"}
}
if len(response) < 32 {
return SmbError{Message: "invalid response length"}
}
return SmbError{Message: "success"}
}

func buildNegotiatePacket() []byte {
const packetLen = 64
packet := make([]byte, packetLen)
copy(packet[4:8], []byte{0xFE, 0x53, 0x4D, 0x42})
packet[8] = 0x40
packet[32] = 0x24
return packet
}

func smbSessionSetup(conn net.Conn, target SmbTarget, hash string) (uint64, error) {
packet := buildSessionSetupPacket(target, hash)
_, err := conn.Write(packet)
if err != nil {
return 0, SmbError{Message: "write failed"}
}
response := make([]byte, 4096)
_, err = conn.Read(response)
if err != nil {
return 0, SmbError{Message: "read failed"}
}
if len(response) < 32 {
return 0, SmbError{Message: "invalid response length"}
}
sessionID := uint64(response[8]) | uint64(response[9])<<8 | uint64(response[10])<<16 | uint64(response[11])<<24 |
uint64(response[12])<<32 | uint64(response[13])<<40 | uint64(response[14])<<48 | uint64(response[15])<<56
if sessionID == 0 {
return 0, SmbError{Message: "authentication failed"}
}
return sessionID, SmbError{Message: "success"}
}

func buildSessionSetupPacket(target SmbTarget, hash string) []byte {
const packetLen = 120
packet := make([]byte, packetLen)
copy(packet[4:8], []byte{0xFE, 0x53, 0x4D, 0x42})
packet[8] = 0x40
for i := 0; i < 16; i++ {
packet[32+i] = hash[i]
}
packet[48] = byte(len(target.Username))
copy(packet[49:], []byte(target.Username))
packet[49+len(target.Username)] = byte(len(target.Domain))
copy(packet[50+len(target.Username):], []byte(target.Domain))
return packet
}

func smbTreeConnect(conn net.Conn, target SmbTarget, share string, sessionID uint64) (uint16, error) {
packet := buildTreeConnectPacket(target, share, sessionID)
_, err := conn.Write(packet)
if err != nil {
return 0, SmbError{Message: "write failed"}
}
response := make([]byte, 4096)
_, err = conn.Read(response)
if err != nil {
return 0, SmbError{Message: "read failed"}
}
if len(response) < 32 {
return 0, SmbError{Message: "invalid response length"}
}
treeID := uint16(response[28]) | uint16(response[29])<<8
return treeID, SmbError{Message: "success"}
}

func buildTreeConnectPacket(target SmbTarget, share string, sessionID uint64) []byte {
treePath := fmt.Sprintf("\\\\%s\\%s", target.Host, share)
pathLen := len(treePath)
totalLen := 32 + 8 + pathLen
packet := make([]byte, totalLen)
copy(packet[4:8], []byte{0xFE, 0x53, 0x4D, 0x42})
packet[8] = 0x40
for i := 0; i < 8; i++ {
packet[24+i] = byte(sessionID >> (8 * i))
}
packet[32] = 0x08
packet[34] = byte(pathLen)
packet[35] = byte(pathLen >> 8)
copy(packet[40:], []byte(treePath))
return packet
}
