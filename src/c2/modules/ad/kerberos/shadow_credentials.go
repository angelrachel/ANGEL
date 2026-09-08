package kerberos

import "os/exec"

type ShadowCredentialResult struct {
Status string
}

type ShadowCredential struct {
TargetUser string
}

func (s ShadowCredential) AddKeyCredential() ShadowCredentialResult {
cmd := exec.Command("certipy", "shadow", "-target", s.TargetUser)
cmd.Run()
return ShadowCredentialResult{Status: "success"}
}

func (s ShadowCredential) RemoveKeyCredential() ShadowCredentialResult {
cmd := exec.Command("certipy", "shadow", "-remove", "-target", s.TargetUser)
cmd.Run()
return ShadowCredentialResult{Status: "success"}
}
