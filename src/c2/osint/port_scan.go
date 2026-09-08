package osint

import (
"fmt"
"net"
"strings"
"time"
)

func PortScan(target string, startPort, endPort int) string {
var openPorts []string
for port := startPort; port <= endPort; port++ {
conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", target, port), 2*time.Second)
if err != nil {
continue
}
conn.Close()
openPorts = append(openPorts, fmt.Sprintf("%d", port))
}
if len(openPorts) == 0 {
return "No open ports found"
}
return "Open ports: " + strings.Join(openPorts, ", ")
}
