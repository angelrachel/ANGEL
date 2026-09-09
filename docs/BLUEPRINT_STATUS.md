# ANGEL Blueprint Status Matrix

This matrix tracks the implementation status of the **ANGEL authorized red-team attack platform** against the supplied pure-attack blueprint. It is an engineering status document, not a claim that every listed capability is complete or production-ready.

The detailed 2026-09-10 structure review is recorded in [`C2_BLUEPRINT_AUDIT.md`](C2_BLUEPRINT_AUDIT.md). The repository currently builds 55 Go packages under `src/c2`; eleven Go test files now cover control-plane foundations, while the blueprint's Angular, .NET 10, and LangGraph components are not present as implemented runtimes.

The architecture boundary is recorded in [`ADR-001-architecture-boundary.md`](ADR-001-architecture-boundary.md). Agent routing is explicitly simulation-only and evidence manifests now require a verified chain before they can be produced.

Status meanings:

- **Foundation:** supporting infrastructure or control-plane primitives exist.
- **Partial:** some source modules exist, but coverage, integration, or platform support is incomplete.
- **Skeleton:** files or interfaces exist, but substantial implementation and validation remain.
- **Planned/gap:** the blueprint area is absent, materially incomplete, or not demonstrated in the current repository.

| Blueprint section | Repository area | Status | Audit note |
|---|---|---|---|
| Architecture overview | `src/`, `src/c2/`, `deploy/` | Foundation/Partial | Multiple deployment and control-plane layers exist; the complete six-layer target architecture is not yet demonstrated end-to-end. |
| C2 framework | `src/c2/`, `src/c2/implant/`, `src/c2/server/` | Partial | Core server, agent, task/result flow, and module directories exist; the 98+ agent target, complete listener set, profile system, and production hardening require verification. |
| Malleable C2 profile | `src/c2/` and related module paths | Planned/gap | The full Teams, Office365, Google, custom profile loader and integrated traffic behavior described in the blueprint are not demonstrated as complete. |
| SMB beacon | `src/c2/modules/lateral/smb_beacon/` | Skeleton/Partial | Named-pipe and peer-to-peer files exist; integration, protocol correctness, and cross-host validation remain incomplete. |
| Decoy/deception layer | `src/c2/server/gateway/`, `deploy/nginx/` | Partial/gap | Reverse-proxy configuration and gateway primitives exist; a decoy implementation is not present. |
| SQL injection engine | `src/c2/modules/exploitation/sqli/` | Partial | Some SQL/auth-bypass source exists; full detector, DBMS exploit, OOB, target parser, and reporting coverage is not demonstrated. |
| NoSQL injection engine | `src/c2/modules/exploitation/nosqli/` | Partial/Skeleton | Several NoSQL-related files exist; the complete MongoDB, Elasticsearch, CouchDB, Redis, and Cassandra coverage is not demonstrated. |
| Database post-exploitation | `src/c2/modules/post_exploitation/` | Skeleton/Partial | DBMS-specific files exist, but many are short or placeholder-like and require implementation, tests, and integration. |
| C2 evasion and stealth | `src/c2/modules/evasion/` | Skeleton/Partial | High-risk evasion paths exist; completeness and platform correctness are unverified and no new operational evasion is supported here. |
| Kerberos and Active Directory | `src/c2/modules/ad/` | Skeleton/Partial | Kerberos, AD reconnaissance, and AD exploit directories exist; the Golden/Silver Ticket, ADCS, DCSync, and persistence coverage is not proven complete. |
| Lateral movement | `src/c2/modules/lateral/` | Skeleton/Partial | SMB, WMI, WinRM, and pivoting files exist; several files are small skeletons and the blueprint's broader protocol coverage is incomplete. |
| Persistence | `src/c2/modules/persistence/` | Partial | Windows, Linux, Darwin, and Android paths exist; platform coverage, rollback, testing, and reliable integration remain incomplete. |
| Hardware rootkit | `src/c2/modules/rootkit/` | Skeleton/Planned | UEFI, SMM, and firmware filenames exist, but real hardware/toolchain validation is not demonstrated. |
| Credential theft | `src/c2/modules/cred/`, `src/c2/modules/credential/`, `src/c2/modules/cred_crack/` | Partial | Browser, cloud, keychain, shadow/history, and cracking interfaces exist; platform correctness, extraction formats, and integration require verification. |
| Collector/infostealer | `src/c2/modules/collector/` | Skeleton/Partial | Collector files exist, but screen, keylog, Wi-Fi, webcam, browser, and evidence flows are not demonstrated as complete. |
| Destruction and impact | `src/c2/modules/destruction/`, `src/c2/modules/destruct_impact/` | Partial | Impact and destruction-related files exist; safe lab controls, recovery validation, and full blueprint coverage remain incomplete. |
| Orchestrator | `src/c2/orchestrator/`, `src/c2/modules/orchestrator/` | Partial | Workflow, policy, guarded routing, attack-graph, and fireteam-related files exist; end-to-end orchestration is not proven. |
| Autonomous decision-making | `src/c2/modules/brain/`, `src/c2/orchestrator/` | Skeleton/Partial | Decision, risk, behavior, and timing files exist; autonomous learning and auditable decision loops require implementation and validation. |
| Infrastructure | `deploy/terraform/`, `deploy/ansible/`, `deploy/nginx/`, `Dockerfile`, `docker-compose.yml` | Foundation/Partial | Infrastructure baseline exists; referenced playbooks/inventories, functional separation, certificate handling, and reproducible deployment require verification. |
| OSINT and reconnaissance | `src/c2/osint/` | Partial | Passive inventory and several reconnaissance modules exist; the blueprint's full DNS, port, web, person, company, and cloud coverage is incomplete. |
| Exploitation | `src/c2/modules/exploitation/` | Partial/Skeleton | Exploitation categories and module paths exist; XSS, SSRF, RCE, LFI/RFI, GraphQL/API, and CVE coverage is not demonstrated end-to-end. |
| Forensic evidence | `src/c2/evidence/` | Foundation/Partial | Hash chain, redaction, and manifest foundations exist; signature, encrypted storage, independent verification, and full collector integration require validation. |
| Reporting | `src/c2/report/`, `src/c2/modules/reporting/` | Partial | Technical/executive reporting support exists; PDF, encrypted delivery, metrics, timeline, and remediation coverage remain incomplete. |
| Cleanup and deletion | `src/c2/cleanup/`, `src/c2/modules/cleanup/` | Partial/Skeleton | Cleanup interfaces exist; safe artifact handling, revocation workflows, database cleanup, cache verification, and deletion manifest behavior require implementation and tests. |

## Repository-level blockers

Generated ELF executables have been removed from source control. The repository now builds the Go control-plane from the Dockerfile's pinned toolchain; release provenance and Terraform provider lockfile review remain required before packaging.

Several modules contain placeholder comments, short skeleton implementations, or behavior that does not match the function name. A file's presence in the tree must not be treated as proof that the corresponding blueprint capability is operational.

The project documentation and deployment references must remain synchronized with this matrix. When a module changes state, update this document, the relevant README section, tests, and release notes together.

## Verification policy

A domain may be marked **implemented** only when its source path, integration path, tests, platform assumptions, error handling, authorization boundary, and deployment behavior have been reviewed. Until then, use **foundation**, **partial**, **skeleton**, or **planned/gap** as appropriate.

All testing must occur within an explicitly authorized and isolated engagement scope. This status document does not itself authorize execution against any target.
