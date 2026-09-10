# ANGEL Legacy Blueprint Mapping and Backlog

## 1. Purpose

This document preserves the old PDF as a reference while mapping every major domain into the new real authorized red-team lab blueprint. The old document remains archived. The new P0/P1 blueprint is the implementation authority.

The mapping prevents domain loss while making completion measurable. Each legacy area is classified as **promoted**, **transformed**, **boundary**, or **backlog**. A source directory, filename, or successful compile does not change a domain to implemented.

## 2. Classification rules

| Classification | Meaning |
|---|---|
| Promoted | Directly implemented as a real lab platform or control-plane deliverable in P0/P1. |
| Transformed | Implemented as a safe owned-lab assessment, detector, validator, evidence, or recovery capability. |
| Boundary | Preserved for architecture traceability but not an operational deliverable in the new blueprint. |
| Backlog | Requires a separate project-owner decision, specialized lab, or additional acceptance criteria before work starts. |
| Legacy partial | Existing files are present but incomplete, unverified, or not integrated. |

## 3. Complete legacy-domain mapping

| Legacy PDF domain | Current GitHub area | New blueprint destination | Classification | Current state and gap |
|---|---|---|---|---|
| Architecture overview | `src/`, `src/c2/`, `deploy/` | L1–L7 target architecture | Promoted | Foundation exists; complete P0 deployment and P1 integration remain. |
| Infrastructure | `deploy/terraform/`, `deploy/ansible/`, `deploy/nginx/`, Docker files | P0-01 and P1-07 | Promoted | Baseline exists; Terraform tool validation and repeatable lab deployment remain. |
| C2 framework | `src/c2/`, `src/c2/server/`, `src/c2/implant/` | P0-03, P0-04, P1-05 | Transformed | Control plane and lab adapters are promoted; operational implant capability is not the P0/P1 target. |
| Malleable C2 profiles | Related C2/profile paths | P1 telemetry and protocol contracts | Boundary | Full traffic-shaping profile system is not a P0/P1 deliverable. |
| HTTP/DNS listeners | `src/c2/server/listener/` | P0 gateway and lab protocol adapters | Promoted | Core listener foundation is tested; deployment and multi-service integration remain. |
| SMB beacon | `src/c2/modules/lateral/smb_beacon/` | P1 network segmentation and protocol fixture | Transformed | Use protocol fixture and detection/segmentation tests; current files are incomplete. |
| Decoy/deception | `src/c2/server/gateway/`, `deploy/nginx/` | P1 service isolation and decoy test service | Transformed | Gateway primitives exist; complete decoy service is not present. |
| SQL injection | `src/c2/modules/exploitation/sqli/` | P1-02 owned database assessment lab | Transformed | Use synthetic database and controlled fixture checks; existing source is partial/unverified. |
| NoSQL injection | `src/c2/modules/exploitation/nosqli/` | P1-02 document database posture lab | Transformed | Use owned document database and validation fixtures; existing coverage is incomplete. |
| Database post-exploitation | `src/c2/modules/post_exploitation/` | P1-02 database posture, audit, and recovery checks | Transformed | DB-specific files exist; integration and tests are missing. |
| Evasion and stealth | `src/c2/modules/evasion/` | P1-04 telemetry, detection, and tamper validation | Transformed | Replace operational bypass goals with measurable detection and audit-integrity outcomes. |
| Kerberos and Active Directory | `src/c2/modules/ad/` | P1-03 identity and directory lab | Transformed | Synthetic identity and authorization posture are promoted; operational abuse paths remain legacy. |
| Lateral movement | `src/c2/modules/lateral/` | P1-04 network segmentation and denied-path tests | Transformed | Network graph and segmentation outcomes are promoted; uncontrolled host movement is not. |
| Persistence | `src/c2/modules/persistence/` | P1-08 resilience, rollback, and persistence detection | Transformed | Detection and rollback are deliverables; installer behavior remains outside P0/P1. |
| Hardware rootkit | `src/c2/modules/rootkit/` | Specialized research boundary | Boundary/backlog | UEFI, SMM, SPI, and JTAG files are unverified; no hardware lab exists. |
| Credential theft | `src/c2/modules/cred/`, `credential/`, `cred_crack/` | P1-03 synthetic identity, secret-handling, and audit checks | Transformed/boundary | Use synthetic credentials and access-policy evidence; extraction behavior is not a deliverable. |
| Collector/InfoStealer | `src/c2/modules/collector/` | P0/P1 evidence pipeline and synthetic collection | Transformed | Collection schema, redaction, provenance, and evidence links are promoted; invasive collection is not. |
| Destruction and impact | `src/c2/modules/destruction/`, `destruct_impact/` | P1 recovery, rollback, and impact assessment | Transformed | Recovery verification is promoted; destructive execution is not. |
| Orchestrator | `src/c2/orchestrator/`, `src/c2/modules/orchestrator/` | P0-03 and P1-05 | Promoted | Scope, guarded routing, and simulation planner foundations exist; full adapter workflow remains. |
| Autonomous brain | `src/c2/modules/brain/` | P1-05 auditable workflow planner | Transformed | Auditable bounded planning is promoted; unrestricted autonomous attack decisions are not. |
| MCP integration | `apps/orchestrator-langgraph/`, related contracts | P0/P1 orchestrator contracts | Promoted | Fixture-only planner exists; real multi-service integration remains. |
| API gateway | `src/c2/server/api/`, `apps/gateway-dotnet/` | P0-02, P0-03, P0-07 | Promoted | Fail-closed boundary exists; authentication, persistence, and full UI integration remain. |
| Frontend | `apps/frontend-angular/` | P0-07 | Promoted | Simulation UI scaffold exists; gateway-backed operator workflow remains. |
| Evidence/forensics | `src/c2/evidence/` | P0-05 | Promoted | Chain, redaction, timestamp, and manifest foundations are tested; collector integration remains. |
| Reporting | `src/c2/report/`, `src/c2/modules/reporting/` | P0-06 and P1-06 | Promoted | Report foundation is tested; complete multi-service findings workflow remains. |
| Cleanup/deletion | `src/c2/cleanup/`, `src/c2/modules/cleanup/` | P0-01 and P1-07 | Transformed | Non-destructive manifests and rollback verification are promoted; destructive cleanup remains legacy. |
| OSINT/recon | `src/c2/osint/` | P0-04 owned-lab inventory and P1 assessment adapters | Transformed | Passive inventory primitives exist; full coverage and integration remain. |
| Exploitation framework | `src/c2/modules/exploitation/` | P1 owned-lab assessment adapters | Transformed/boundary | Controlled fixture assessments are promoted; generalized target exploitation is not. |
| Teamserver and deployment | `Dockerfile`, compose, deploy scripts | P0-01, P0-08, P1-07 | Promoted | Build and guarded scripts exist; clean lab deployment still needs full validation. |

## 4. Current GitHub status by implementation state

### 4.1 Implemented or foundation-tested

- Go control-plane packages and selected tests.
- Scope and authorization primitives.
- Simulation-only guarded agent router.
- Evidence ledger, manifest, redaction, and chain verification.
- Report generator and delivery validation.
- Fixture-only replay runner.
- Cleanup manifest validation.
- Angular simulation workspace scaffold.
- .NET fail-closed gateway boundary.
- LangGraph fixture-only planner.
- Gateway and planner unit tests.
- Safe contract validation.
- Makefile and CI source-boundary checks.
- Repository integrity and module inventory generation.

### 4.2 Partial

- C2 server and listener integration.
- OSINT and reconnaissance coverage.
- Infrastructure configuration.
- Reporting integration with live task records.
- Cleanup workflow integration.
- API gateway persistence and authentication.
- Angular gateway-backed task workflow.
- Orchestrator multi-step workflow.
- Deployment reproducibility.

### 4.3 Skeleton or placeholder-like

- Database post-exploitation packages.
- NoSQL package coverage.
- AD and Kerberos package coverage.
- Brain and autonomous decision packages.
- Reporting module wrapper.
- Hardware rootkit paths.
- Several lateral and persistence subpackages.
- Several exploitation subpackages.

### 4.4 High-risk legacy areas requiring separate review

- Implant source paths.
- Credential and secret extraction paths.
- Persistence installers.
- Evasion and log-cleanup paths.
- Lateral movement and pivoting paths.
- Rootkit and firmware paths.
- Destruction and impact paths.
- General exploitation paths.
- Exfiltration-related behavior if introduced later.

These areas remain represented in the mapping so the old PDF is not silently discarded. Their presence does not make them P0/P1 deliverables.

## 5. Known error and quality backlog

| Finding | Current evidence | Required resolution for the new blueprint |
|---|---|---|
| 167 success-literal findings | Static checker | Replace false-success paths in promoted packages with verified outcomes and structured errors. |
| 146 ignored-command-error findings | Static checker | Handle command failures explicitly in promoted tooling; do not promote ignored-error paths. |
| Most offensive packages lack local tests | Module inventory | Add tests only for transformed owned-lab adapters, detectors, validators, and policy boundaries. |
| Terraform binary unavailable in current environment | Toolchain check | Install or provide the pinned tool, then run static validation against the isolated lab only. |
| `BLUEPRINT_STATUS.md` contains stale scaffold wording | Documentation review | Synchronize status matrix with this mapping and the new P0/P1 blueprint. |
| Angular UI is not yet gateway-backed | Application review | Add typed API client, authentication context, polling/event updates, and negative integration tests. |
| .NET gateway has validation boundary but no persistent task store | Application review | Connect it to the control-plane contract and add idempotency and state-transition tests. |
| LangGraph planner is fixture-only | Application review | Add bounded lab workflow nodes only after P0 task/evidence contracts are integrated. |
| Deployment scripts are guarded but not fully lab-verified | Deployment review | Execute provision, health, rollback, and teardown acceptance tests in the isolated lab. |

## 6. P0/P1 traceability matrix

| New deliverable | Existing files to reuse | New work required | Acceptance evidence |
|---|---|---|---|
| P0-01 Lab | `deploy/`, Docker files, `simulation/` | Lab composition, seed, health, teardown, network tests | Clean provision/smoke/teardown run |
| P0-02 Policy | `src/c2/orchestrator/scope.go`, agent router | API persistence, approval records, emergency stop | Positive/negative authorization suite |
| P0-03 Lifecycle | `src/c2/server/task/`, result files | State machine, idempotency, retries, cancellation | Contract/integration suite |
| P0-04 Adapters | `src/c2/osint/`, server primitives | Allowlisted owned-lab adapters | Real synthetic service run |
| P0-05 Evidence | `src/c2/evidence/` | Finding linkage and service integration | Tamper/redaction/chain suite |
| P0-06 Reports | `src/c2/report/` | Full task-to-report pipeline | Deterministic report acceptance |
| P0-07 UI | `apps/frontend-angular/` | Gateway client and real workflow | Browser/integration workflow |
| P0-08 CI | `.github/workflows/`, Makefile, scripts | Lab smoke and release gates | Clean checkout release pass |
| P1-01 Web/API lab | `contracts/`, fixtures | Resettable synthetic services | Assessment scenario suite |
| P1-02 Database lab | Existing DB-related source and contracts | Seeded relational/document services | Posture/evidence/report suite |
| P1-03 Identity lab | AD/Kerberos paths and scope policy | Synthetic identity service and fixtures | Role/account lifecycle suite |
| P1-04 Telemetry | Evidence, report, server primitives | Event bus, correlation, alerts | Timeline and denied-path suite |
| P1-05 Orchestration | Go orchestrator and LangGraph planner | Bounded multi-step workflow | Failure/retry/cancel/rollback suite |
| P1-06 Findings | Report and evidence packages | Severity, dedupe, remediation state | Finding quality suite |
| P1-07 Deployment | Terraform/Ansible/Nginx | Tool validation, deploy, rollback, drift checks | Isolated lab deployment run |
| P1-08 Hardening | API, crypto, database, deploy paths | Secret injection, rate limits, rotation, resource limits | Resilience/security suite |

## 7. What is deliberately not omitted

The old domain families remain mapped: C2, profiles, listeners, SMB, deception, SQL, NoSQL, database post-exploitation, evasion, AD/Kerberos, lateral movement, persistence, rootkit, credential access, collectors, destruction, orchestrator, autonomous decision-making, infrastructure, OSINT, exploitation, evidence, reporting, and cleanup. Each one has a destination, status, and acceptance boundary in this document.

The new blueprint is therefore not a beginner scan-only plan. It is a hard/critical engineering plan for a real isolated assessment platform, with real service deployment, real controlled lab interactions, real evidence, real reporting, real rollback, and release gates. The distinction is that operational high-risk capabilities are not silently treated as completed P0/P1 work merely because they appear in the legacy tree.

## 8. Document relationship

- `ANGEL_REAL_LAB_P0_P1_BLUEPRINT.md` is the new implementation authority.
- This file is the traceability, migration, error, and backlog authority.
- The original PDF remains archived as historical input.
- GitHub should be changed only after the owner approves this new blueprint pair.

## 9. References

[1]: https://github.com/angelrachel/ANGEL "ANGEL repository"
[2]: https://owasp.org/www-project-application-security-verification-standard/ "OWASP Application Security Verification Standard"
[3]: https://owasp.org/www-project-web-security-testing-guide/ "OWASP Web Security Testing Guide"
[4]: https://www.nist.gov/cyberframework "NIST Cybersecurity Framework"
[5]: https://json-schema.org/ "JSON Schema documentation"
