# ANGEL — Authorized Red-Team Attack Platform

ANGEL adalah platform modular untuk **authorized red-team operations** dan security research di lingkungan yang memiliki izin tertulis. Repository ini merepresentasikan project pure-attack yang terpisah dari project defensive/safety control-plane.

> **Authorization required:** gunakan hanya pada lab, aset milik sendiri, atau engagement dengan ruang lingkup dan Rules of Engagement tertulis. Repository ini tidak memberikan izin untuk menguji sistem pihak ketiga.

## Project identity

ANGEL dirancang sebagai ekosistem multi-layer yang menggabungkan C2 framework, attack modules, orchestrator, infrastructure, evidence, reporting, dan cleanup lifecycle. Blueprint utama project ini adalah `struktursealangel_ANGEL(3).pdf`.

Repository ini masih dalam tahap pengembangan. Nama file atau folder tidak otomatis berarti seluruh capability telah selesai. Setiap domain harus dinilai berdasarkan status aktualnya: **implemented**, **partial**, **skeleton**, atau **planned**.

## Target architecture

| Layer | Domain | Target responsibility |
|---|---|---|
| Layer 1 | Infrastructure | Provisioning, separation, redirector, VPN, firewall, and deployment configuration |
| Layer 2 | C2 Framework | Implant, teamserver, listeners, profiles, cryptography, task/result flow, and SMB beacon |
| Layer 3 | Orchestrator | Routing, workflow coordination, fireteam execution, approvals, state, and MCP integration |
| Layer 4 | API Gateway | Operator API, authentication, authorization, validation, rate limiting, and real-time transport |
| Layer 5 | Frontend | Operator dashboard, agent console, monitoring, evidence, and report viewer |
| Layer 6 | Assessment Domains | Reconnaissance, exploitation, post-exploitation, persistence, collection, impact, evidence, reporting, and cleanup |

The target design follows functional separation between infrastructure components. Actual deployment readiness must be verified against the Terraform, Ansible, container, and reverse-proxy configuration before use.

## Blueprint coverage

| Blueprint domain | Repository area | Current audit status |
|---|---|---|
| C2 framework | `src/c2/` | Partial; core and module skeletons exist |
| Decoy/deception | `src/c2/server/gateway/`, `deploy/nginx/` | Partial; decoy implementation is not present |
| SQL injection | `src/c2/modules/exploitation/sqli/` | Partial; coverage and safety boundaries require verification |
| NoSQL injection | `src/c2/modules/exploitation/nosqli/` | Partial |
| Database post-exploitation | `src/c2/modules/post_exploitation/` | Skeleton/partial; implementation must be verified per DBMS |
| Evasion and stealth | `src/c2/modules/evasion/` | Skeleton/partial; contains high-risk paths that are not operationally supported |
| Kerberos and Active Directory | `src/c2/modules/ad/` | Skeleton/partial |
| Lateral movement | `src/c2/modules/lateral/` | Skeleton/partial |
| Persistence | `src/c2/modules/persistence/` | Partial; platform coverage is incomplete |
| Hardware rootkit | `src/c2/modules/rootkit/` | Skeleton/planned; hardware validation is not established |
| Credential access | `src/c2/modules/cred/`, `src/c2/modules/credential/` | Partial; platform-specific behavior requires verification |
| Collector/infostealer | `src/c2/modules/collector/` | Skeleton/partial |
| Destruction and impact | `src/c2/modules/destruct_impact/`, `src/c2/modules/destruction/` | Partial; use only in isolated, approved test environments |
| Orchestrator and brain | `src/c2/orchestrator/`, `src/c2/modules/orchestrator/`, `src/c2/modules/brain/` | Partial |
| Infrastructure | `deploy/` | Lab baseline; deployment references require validation |
| OSINT and reconnaissance | `src/c2/osint/` | Partial; passive and active coverage differ |
| Exploitation | `src/c2/modules/exploitation/` | Partial |
| Forensic evidence | `src/c2/evidence/` | Foundation implemented; end-to-end verification required |
| Reporting | `src/c2/report/`, `src/c2/modules/reporting/` | Partial |
| Cleanup and deletion | `src/c2/cleanup/`, `src/c2/modules/cleanup/` | Partial; cleanup behavior must be verified before use |

This table is an engineering status guide, not a claim that every blueprint item is complete.

## Repository layout

```text
ANGEL/
├── src/
│   ├── c2/                 # Go/Python C2, server, implant, and modules
│   ├── orchestrator/       # Workflow and policy primitives
│   ├── osint/              # Reconnaissance and inventory support
│   ├── evidence/           # Evidence normalization and integrity support
│   └── report/             # Report export support
├── deploy/                 # Terraform, Ansible, Nginx, and deployment scripts
├── tests/                  # Automated tests currently focused on foundation code
├── docs/                   # Operations and status documentation
└── go.mod                # Go module and tooling configuration
```

## Development status

The repository contains a mixture of working foundation code, partial modules, and placeholders. Before treating a domain as operational, verify its source implementation, tests, platform assumptions, error handling, authorization boundary, and deployment path.

The following items remain important engineering work:

1. Complete the missing or partial blueprint components and document their real status.
2. Validate all Go and Python build targets with pinned toolchain versions.
3. Add tests for C2 protocol behavior, module boundaries, deployment configuration, and failure handling.
4. Replace hardcoded development values with configured secrets and environment-specific settings.
5. Remove generated binaries, local databases, debug artifacts, and other build output from source control unless there is a documented release reason.
6. Verify Terraform, Ansible, Nginx, container, certificate, and network-separation configuration in an isolated lab.
7. Keep evidence, authorization records, target inventory, and engagement data outside the public source tree.

Generated status references:

- [`docs/MODULE_INVENTORY.md`](docs/MODULE_INVENTORY.md) — conservative package inventory based on source and test evidence.
- [`docs/ADR-001-architecture-boundary.md`](docs/ADR-001-architecture-boundary.md) — current Go-only architecture boundary.
- [`docs/ARCHITECTURE_SCAFFOLD.md`](docs/ARCHITECTURE_SCAFFOLD.md) — contract-first folders, toolchain status, and next coding tasks.
- `deploy/scripts/validate.sh` — static Terraform and shell validation; it does not apply or destroy infrastructure.
- `Makefile` — local format, test, build, and validation commands.

## Local validation

The repository contains a Go control plane plus safe application boundary scaffolds. Use the pinned Go toolchain and run the Go baseline checks:

```bash
go test ./src/...
go vet ./src/...
go build ./src/...
```

The Docker image is built with the same Go control-plane entrypoint. A successful documentation check must not be inferred unless the required tools are installed and the commands complete successfully.

## Security and operational boundary

ANGEL is intended for controlled, authorized red-team research. Keep test assets isolated, define an engagement scope before execution, use synthetic data where possible, preserve evidence according to the engagement policy, and stop when an asset falls outside the approved scope.

Do not commit credentials, tokens, private keys, customer data, target inventories, engagement evidence, generated binaries, or unreviewed destructive test artifacts. Review `SECURITY.md` and `docs/OPERATIONS.md` before any lab deployment.

## License and authorization

No production deployment or external assessment is authorized by this repository alone. Obtain explicit written permission from the system owner, define the permitted targets and techniques, and document rollback and emergency procedures before use.
