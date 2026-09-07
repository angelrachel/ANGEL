package main

import (
    "fmt"
    "os"
    "time"
)

func main() {
    hostname, _ := os.Hostname()
    fmt.Printf("[+] Linux implant started on %s\n", hostname)
    fmt.Println("[+] Waiting for tasks...")

    for {
        time.Sleep(5 * time.Second)
    }
}
