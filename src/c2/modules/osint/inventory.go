package osint

import (
"context"
"fmt"
"net"
"strconv"
"time"
)

type Inventory struct {
Target string
IPs    []string
Ports  []int
}

func NewInventory(target string) *Inventory {
return &Inventory{
Target: target,
}
}

func (i *Inventory) ResolveIPs() error {
ips, err := net.LookupHost(i.Target)
if err != nil {
return err
}
i.IPs = ips
return nil
}

func (i *Inventory) ScanPorts(timeout time.Duration) error {
scanner := NewPortScanner(i.Target, timeout, "tcp")
results, err := scanner.ScanCommonPorts(context.Background())
if err != nil {
return err
}
for _, result := range results {
i.Ports = append(i.Ports, result.Port)
}
return nil
}

func (i *Inventory) GetIPs() []string {
return i.IPs
}

func (i *Inventory) GetPorts() []int {
return i.Ports
}

func (i *Inventory) Print() string {
return fmt.Sprintf("Target: %s\nIPs: %v\nPorts: %v", i.Target, i.IPs, i.Ports)
}
