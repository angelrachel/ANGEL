package orchestrator

import (
"strings"
)

type IntentRouter struct {
Type AgentType
}

func (i *IntentRouter) ClassifyIntent(input string) AgentType {
lower := strings.ToLower(input)
if strings.Contains(lower, "scan") || strings.Contains(lower, "recon") || strings.Contains(lower, "subdomain") {
return AgentRecon
}
if strings.Contains(lower, "sql") || strings.Contains(lower, "database") {
return AgentSQL
}
if strings.Contains(lower, "nosql") || strings.Contains(lower, "mongo") {
return AgentNoSQL
}
if strings.Contains(lower, "exploit") || strings.Contains(lower, "rce") {
return AgentPost
}
if strings.Contains(lower, "lateral") || strings.Contains(lower, "smb") || strings.Contains(lower, "wmi") {
return AgentLateral
}
if strings.Contains(lower, "destroy") || strings.Contains(lower, "wipe") || strings.Contains(lower, "delete") {
return AgentDestroy
}
if strings.Contains(lower, "credential") || strings.Contains(lower, "password") || strings.Contains(lower, "hash") {
return AgentCred
}
return AgentRecon
}

func (i *IntentRouter) GetAgentType() AgentType {
return i.Type
}
