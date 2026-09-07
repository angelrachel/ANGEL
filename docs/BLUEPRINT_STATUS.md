# Blueprint Status Matrix

This matrix records what is implemented in the repository and what remains intentionally outside the safe control-plane scope.

| Blueprint section | Status | Notes |
|---|---|---|
| Architecture overview | Implemented | Python control plane, SQLite state, event stream, authenticated HTTP API |
| C2 framework | Safe foundation only | Synthetic agent registration, task queue, encrypted envelopes, replay protection, and allowlisted simulator tasks |
| Decoy/deception layer | Not implemented | Deception and traffic mimicry are not required for the defensive control plane |
| SQL/NoSQL injection | Defensive review only | Schema and API contract review; no exploit or exfiltration engine |
| Database post-exploitation | Not implemented | Owner-only decision requiring a separately authorized lab |
| Evasion and stealth | Not implemented | No EDR bypass, anti-analysis, process injection, or log clearing |
| Kerberos/Active Directory attacks | Not implemented | Passive inventory may be added; ticket abuse and credential extraction remain excluded |
| Lateral movement | Not implemented | No SMB/WMI/WinRM execution or tunneling |
| Persistence | Not implemented | No registry, scheduled-task, service, cron, UEFI, or Android persistence |
| Hardware rootkit | Not implemented | Explicitly excluded |
| Credential theft | Not implemented | Explicitly excluded |
| Collector/infostealer | Not implemented | No keylogging, browser theft, screen capture, webcam, or Wi-Fi extraction |
| Destruction and impact | Not implemented | No wiper, ransomware, sabotage, or destructive actions |
| Orchestrator | Implemented safely | Bounded workflows, approvals, replay, cancellation, circuit breaker, and MCP policy |
| Autonomous decision-making | Safe deterministic controls only | No autonomous offensive action selection |
| Infrastructure | Implemented and hardened | Terraform, Ansible, Docker, Nginx, local binding, non-root, read-only runtime |
| OSINT and reconnaissance | Implemented passively | Scope-aware inventory, certificate metadata, response fingerprinting, passive sources |
| Exploitation | Defensive validation only | BOLA canary checks and schema drift analysis; no arbitrary exploit delivery |
| Forensic evidence | Implemented | Redaction, hash chain, backup integrity, manifest verification |
| Reporting | Implemented | Technical findings, executive summary, Markdown/JSON export, passive finding adapters |
| Cleanup and deletion | Safe lifecycle only | Credential rotation and approved lab cleanup remain owner procedures; forensic concealment is excluded |

The owner must provide written authorization, target inventory, retention rules, production secrets, infrastructure approval, and human review before any future active assessment module is considered.
