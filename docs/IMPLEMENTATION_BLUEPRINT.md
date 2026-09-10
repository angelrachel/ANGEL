# ANGEL Full Implementation Blueprint

## Status and authority

This document is the detailed implementation specification for the ANGEL real authorized red-team and security-validation platform. The earlier PDF remains an archived source document. This Markdown document is the active technical authority for repository structure, services, workflows, acceptance criteria, and release sequencing.

The platform operates against registered assets inside an isolated environment owned by the project or explicitly authorized by the project owner. Every action is authenticated, scoped, policy-checked, observable, evidence-linked, reversible where state changes occur, and covered by a release test.

The platform retains the depth of the original domain model while replacing operational mechanisms that are unsuitable as unrestricted deliverables with concrete security-engineering components that can be built, deployed, used, measured, and verified in the project-owned lab.

## 1. Introduction

### 1.1 Background

ANGEL is a multi-service security assessment platform for authorized engagements. It combines infrastructure, control-plane services, authenticated assessment workers, orchestration, operator experience, evidence, reporting, detection, recovery, and deployment automation.

The system is designed for P0 Critical and P1 Hard delivery. Completion means that the relevant component runs in the lab, communicates through typed contracts, handles failure, records evidence, integrates with dependent services, and passes its acceptance suite.

### 1.2 Objectives

The platform must:

1. Establish and enforce an explicit asset scope.
2. Authenticate operators, services, and assessment workers.
3. Execute approved assessments against registered lab services.
4. Produce findings with severity, confidence, evidence, and remediation.
5. Correlate task, policy, service, telemetry, and recovery events.
6. Generate technical and executive reports.
7. Restore changed lab state and verify cleanup.
8. Reproduce the complete deployment from a clean checkout.
9. Provide P0 and P1 evidence that is reviewable by engineering and HRD.

### 1.3 Scope

The operational scope includes:

- Private lab infrastructure.
- Web and API services.
- Relational and document data services.
- Identity and directory services.
- Network segmentation and telemetry.
- Authenticated assessment workers.
- Evidence and reporting.
- Detection engineering.
- Configuration, drift, recovery, and artifact integrity.
- CI, deployment, backup, restore, rollback, and teardown.

All assets are represented by stable identifiers in a registry. The control plane does not accept unregistered addresses, arbitrary command strings, or unscoped destinations.

### 1.4 Completion policy

| State | Meaning | Evidence |
|---|---|---|
| `planned` | Requirement is documented. | Blueprint section and issue identifier. |
| `structured` | Folder, contract, and configuration are present. | Repository tree and contract validation. |
| `implemented` | Source and local tests exist. | Build and unit test output. |
| `integrated` | Component communicates with dependencies. | Contract and integration tests. |
| `lab-verified` | Component operates against registered lab services. | Lab run, evidence records, and report. |
| `released` | Complete P0/P1 acceptance suite passes. | Release record with deployment and teardown evidence. |

## 2. Architecture

### 2.1 Design principles

The platform uses functional separation, typed contracts, explicit policy gates, bounded workflows, immutable evidence, deterministic lab state, and recoverable changes.

### 2.2 Seven-layer architecture

```text
┌──────────────────────────────────────────────────────────────────────┐
│ L7 VERIFICATION AND RELEASE                                          │
│ CI, contract tests, integration tests, lab acceptance, recovery,     │
│ rollback, integrity checks, dependency and artifact verification     │
└──────────────────────────────────────────────────────────────────────┘
                              │
┌──────────────────────────────────────────────────────────────────────┐
│ L6 EVIDENCE AND REPORTING                                            │
│ Hash-linked evidence, findings, severity, remediation, timeline,     │
│ executive report, technical report, export and chain verification    │
└──────────────────────────────────────────────────────────────────────┘
                              │
┌──────────────────────────────────────────────────────────────────────┐
│ L5 OPERATOR EXPERIENCE                                               │
│ Angular console, assets, scope, approvals, workers, tasks, findings, │
│ evidence, reports, telemetry, recovery and deployment status         │
└──────────────────────────────────────────────────────────────────────┘
                              │
┌──────────────────────────────────────────────────────────────────────┐
│ L4 ORCHESTRATION                                                     │
│ LangGraph workflows, dependency graph, approval pause, retries,      │
│ cancellation, timeout, risk score, compensation and audit            │
└──────────────────────────────────────────────────────────────────────┘
                              │
┌──────────────────────────────────────────────────────────────────────┐
│ L3 AUTHENTICATED ASSESSMENT WORKERS                                  │
│ HTTP, API, database, identity, network, telemetry, configuration,    │
│ artifact, recovery and evidence workers                              │
└──────────────────────────────────────────────────────────────────────┘
                              │
┌──────────────────────────────────────────────────────────────────────┐
│ L2 CONTROL PLANE AND GATEWAY                                         │
│ Go policy and task core, .NET API gateway, identity, RBAC, rate      │
│ limits, state machine, audit, emergency stop and worker registry     │
└──────────────────────────────────────────────────────────────────────┘
                              │
┌──────────────────────────────────────────────────────────────────────┐
│ L1 PRIVATE LAB INFRASTRUCTURE                                        │
│ Terraform, Ansible, containers, network segments, services, identity,│
│ databases, telemetry, secrets, backups and verified teardown         │
└──────────────────────────────────────────────────────────────────────┘
```

### 2.3 Communication model

| Communication | Contract | Controls | Evidence |
|---|---|---|---|
| Angular to gateway | OpenAPI and typed client | Authentication, role, rate limit, schema | Request ID, response status, audit record |
| Gateway to Go control plane | Internal API contract | Service identity, scope, policy version | Decision record and task record |
| Control plane to worker | Worker task contract | Worker identity, capability, asset scope, expiry | Dispatch, receipt, result signature |
| Worker to lab service | Registered asset contract | Asset allowlist, operation policy, timeout | Service response, finding, telemetry |
| Worker to evidence service | Evidence contract | Provenance, redaction, hash chain | Evidence record and parent hash |
| Orchestrator to services | Workflow contract | Approval pause, retry budget, cancellation | Workflow state and decision trace |
| Deployment to services | Infrastructure contract | Private network, secret injection, health checks | Deployment state and teardown record |

## 3. Repository structure

### 3.1 Target repository tree

```text
ANGEL/
├── apps/
│   ├── frontend-angular/
│   │   ├── src/app/
│   │   │   ├── core/
│   │   │   │   ├── auth/
│   │   │   │   ├── api/
│   │   │   │   ├── guards/
│   │   │   │   └── interceptors/
│   │   │   ├── features/
│   │   │   │   ├── dashboard/
│   │   │   │   ├── assets/
│   │   │   │   ├── scope/
│   │   │   │   ├── approvals/
│   │   │   │   ├── workers/
│   │   │   │   ├── tasks/
│   │   │   │   ├── findings/
│   │   │   │   ├── evidence/
│   │   │   │   ├── reports/
│   │   │   │   ├── telemetry/
│   │   │   │   └── recovery/
│   │   │   ├── models/
│   │   │   ├── app.config.ts
│   │   │   ├── app.routes.ts
│   │   │   ├── app.ts
│   │   │   ├── app.html
│   │   │   └── app.scss
│   │   ├── angular.json
│   │   ├── package.json
│   │   └── README.md
│   ├── gateway-dotnet/
│   │   ├── Authentication/
│   │   ├── Authorization/
│   │   ├── Contracts/
│   │   ├── Endpoints/
│   │   ├── Middleware/
│   │   ├── Persistence/
│   │   ├── Program.cs
│   │   ├── SimulationTask.cs
│   │   └── AngelGateway.csproj
│   ├── gateway-dotnet.tests/
│   └── orchestrator-langgraph/
│       ├── nodes/
│       ├── policies/
│       ├── workflows/
│       ├── planner.py
│       ├── test_planner.py
│       └── requirements.txt
├── contracts/
│   ├── api/openapi.yaml
│   ├── authorization/
│   ├── assets/
│   ├── evidence/
│   ├── findings/
│   ├── identity/
│   ├── reports/
│   ├── tasks/
│   ├── telemetry/
│   └── workers/
├── lab/
│   ├── compose/
│   ├── services/
│   │   ├── web-app/
│   │   ├── api-service/
│   │   ├── relational-db/
│   │   ├── document-db/
│   │   ├── identity-service/
│   │   ├── telemetry-service/
│   │   ├── evidence-service/
│   │   └── report-service/
│   ├── network/
│   ├── seed/
│   ├── scenarios/
│   └── acceptance/
├── src/c2/
│   ├── cleanup/
│   ├── evidence/
│   ├── modules/
│   │   ├── assessment_workers/
│   │   ├── asset_intelligence/
│   │   ├── configuration/
│   │   ├── database_security/
│   │   ├── deception/
│   │   ├── identity_security/
│   │   ├── network_security/
│   │   ├── recovery/
│   │   ├── reporting/
│   │   └── telemetry/
│   ├── orchestrator/
│   ├── osint/
│   ├── report/
│   ├── server/
│   │   ├── api/
│   │   ├── crypto/
│   │   ├── database/
│   │   ├── gateway/
│   │   ├── listener/
│   │   └── task/
│   └── simulation/replay/
├── deploy/
│   ├── terraform/
│   ├── ansible/
│   ├── nginx/
│   ├── compose/
│   └── scripts/
├── scripts/
├── tests/
│   ├── contract/
│   ├── integration/
│   ├── lab/
│   ├── recovery/
│   └── release/
├── docs/
│   ├── BLUEPRINT_STATUS.md
│   ├── ARCHITECTURE_SCAFFOLD.md
│   ├── ACCEPTANCE_MATRIX.md
│   ├── DEPLOYMENT.md
│   ├── OPERATIONS.md
│   └── RELEASE_NOTES.md
├── Dockerfile
├── docker-compose.yml
├── Makefile
└── go.mod
```

### 3.2 Repository file responsibilities

| Area | File responsibility | Required state |
|---|---|---|
| `contracts/` | Canonical request, response, state, evidence, worker, identity, and report schemas | Versioned and contract-tested |
| `src/c2/server/` | Go API, task store, crypto, database, listeners, and gateway primitives | Unit, integration, and failure tests |
| `src/c2/orchestrator/` | Scope, policy, workflow state, routing, dependency graph, and compensation | Policy and workflow acceptance tests |
| `src/c2/modules/assessment_workers/` | Typed workers for registered lab assessments | Operational lab integration |
| `src/c2/evidence/` | Chain, redaction, provenance, verification, and manifest | Tamper and redaction tests |
| `src/c2/report/` | Technical and executive report generation | Deterministic report tests |
| `apps/gateway-dotnet/` | External API, authentication, authorization, rate limiting, and typed contract boundary | API integration tests |
| `apps/frontend-angular/` | Operator console and workflow views | Browser and gateway integration tests |
| `apps/orchestrator-langgraph/` | Bounded workflow graph and explainable decisions | Planner and failure-path tests |
| `lab/` | Registered services, state seed, network, telemetry, and acceptance procedures | Clean provision, run, restore, teardown |
| `deploy/` | Private lab deployment and rollback | Static validation and lab deployment pass |
| `tests/` | Cross-component and release acceptance | P0/P1 release gates |
| `docs/` | Operating procedures and status | Synchronized with code and releases |

## 4. Backend control plane

### 4.1 Go package structure

```text
src/c2/
├── server/
│   ├── api/
│   │   ├── routes.go
│   │   ├── handlers.go
│   │   ├── middleware.go
│   │   ├── request_validation.go
│   │   └── api_test.go
│   ├── crypto/
│   │   ├── aes.go
│   │   ├── certificate.go
│   │   ├── signing.go
│   │   └── crypto_test.go
│   ├── database/
│   │   ├── sqlite.go
│   │   ├── models.go
│   │   ├── migrations.go
│   │   └── database_test.go
│   ├── gateway/
│   │   ├── route_selector.go
│   │   ├── request_policy.go
│   │   └── gateway_test.go
│   ├── listener/
│   │   ├── http.go
│   │   ├── websocket.go
│   │   ├── dns.go
│   │   └── listener_test.go
│   └── task/
│       ├── queue.go
│       ├── state_machine.go
│       ├── result.go
│       └── task_test.go
├── orchestrator/
│   ├── scope.go
│   ├── approval.go
│   ├── router.go
│   ├── graph.go
│   ├── compensation.go
│   └── orchestrator_test.go
├── modules/
│   ├── assessment_workers/
│   │   ├── http_posture.go
│   │   ├── api_contract.go
│   │   ├── database_policy.go
│   │   ├── identity_policy.go
│   │   ├── network_policy.go
│   │   ├── telemetry_delivery.go
│   │   └── worker_test.go
│   ├── asset_intelligence/
│   │   ├── inventory.go
│   │   ├── ownership.go
│   │   ├── service_metadata.go
│   │   └── inventory_test.go
│   ├── configuration/
│   │   ├── baseline.go
│   │   ├── drift.go
│   │   ├── remediation.go
│   │   └── configuration_test.go
│   ├── database_security/
│   │   ├── relational_policy.go
│   │   ├── document_policy.go
│   │   ├── privilege_graph.go
│   │   ├── audit_correlation.go
│   │   └── database_security_test.go
│   ├── deception/
│   │   ├── decoy_service.go
│   │   ├── canary_event.go
│   │   ├── alert_correlation.go
│   │   └── deception_test.go
│   ├── identity_security/
│   │   ├── account_lifecycle.go
│   │   ├── role_graph.go
│   │   ├── stale_identity.go
│   │   └── identity_test.go
│   ├── network_security/
│   │   ├── segment_graph.go
│   │   ├── allowed_path.go
│   │   ├── denied_path.go
│   │   └── network_test.go
│   ├── recovery/
│   │   ├── backup.go
│   │   ├── restore.go
│   │   ├── rollback.go
│   │   ├── teardown.go
│   │   └── recovery_test.go
│   ├── reporting/
│   │   ├── findings.go
│   │   ├── severity.go
│   │   ├── remediation.go
│   │   └── reporting_test.go
│   └── telemetry/
│       ├── events.go
│       ├── correlation.go
│       ├── detection.go
│       └── telemetry_test.go
├── evidence/
│   ├── ledger.go
│   ├── manifest.go
│   ├── redaction.go
│   ├── verification.go
│   └── evidence_test.go
├── report/
│   ├── report_generator.go
│   ├── executive.go
│   ├── delivery.go
│   └── report_test.go
└── simulation/replay/
    ├── replay.go
    ├── fixture_test.go
    └── replay_test.go
```

### 4.2 Backend task lifecycle

```text
CREATE
  │
  ▼
VALIDATE CONTRACT
  │
  ▼
CHECK IDENTITY → CHECK ASSET SCOPE → CHECK APPROVAL → CHECK POLICY
  │                         │
  │                         └── REJECT + EVIDENCE
  ▼
QUEUE
  │
  ▼
DISPATCH AUTHENTICATED WORKER
  │
  ├── RUNNING → RESULT VERIFY → EVIDENCE → FINDING → REPORT
  ├── TIMEOUT → CANCEL → COMPENSATE → EVIDENCE
  ├── FAILURE → RETRY BUDGET → FAILED → EVIDENCE
  └── EMERGENCY STOP → STOPPED → RECOVERY RECORD
```

### 4.3 Backend input, process, output, and evidence

| Component | Input | Process | Output | Evidence |
|---|---|---|---|---|
| Scope policy | Asset ID, role, technique, time, approval | Validate policy and version | Allow or deny decision | Policy decision record |
| Task queue | Approved task | Validate state, idempotency, retry budget | Queue state | Task state record |
| Worker registry | Worker identity and capability | Verify credential and capability | Registered or quarantined worker | Enrollment and health record |
| Assessment worker | Registered asset ID and typed operation | Execute approved service check | Structured result | Request, response, finding, timestamp |
| Evidence ledger | Event and provenance | Hash, link, redact, verify | Record and chain status | Parent hash and verification result |
| Report service | Findings, evidence, timeline | Aggregate, score, render | Technical and executive report | Report hash and export record |
| Recovery service | Resource inventory and version | Restore, rollback, verify | Recovery status | Recovery and teardown record |

## 5. API gateway and contracts

### 5.1 Gateway structure

```text
apps/gateway-dotnet/
├── Authentication/
│   ├── OperatorAuthenticator.cs
│   ├── WorkerAuthenticator.cs
│   └── CredentialRevocation.cs
├── Authorization/
│   ├── ScopePolicy.cs
│   ├── RolePolicy.cs
│   ├── ApprovalPolicy.cs
│   └── EmergencyStopPolicy.cs
├── Contracts/
│   ├── AssetContracts.cs
│   ├── TaskContracts.cs
│   ├── WorkerContracts.cs
│   ├── EvidenceContracts.cs
│   └── ReportContracts.cs
├── Endpoints/
│   ├── HealthEndpoints.cs
│   ├── AssetEndpoints.cs
│   ├── ScopeEndpoints.cs
│   ├── ApprovalEndpoints.cs
│   ├── TaskEndpoints.cs
│   ├── WorkerEndpoints.cs
│   ├── EvidenceEndpoints.cs
│   └── ReportEndpoints.cs
├── Middleware/
│   ├── RequestIdMiddleware.cs
│   ├── RateLimitMiddleware.cs
│   ├── ErrorMiddleware.cs
│   └── AuditMiddleware.cs
├── Persistence/
│   ├── TaskRepository.cs
│   ├── ScopeRepository.cs
│   └── AuditRepository.cs
├── Program.cs
└── AngelGateway.csproj
```

### 5.2 Required API groups

| API group | Required operations | P0/P1 result |
|---|---|---|
| Health | Health, version, dependency status | Service readiness and release gate |
| Assets | Register, inspect, approve, retire | Stable scope identifiers |
| Scope | Create, activate, pause, close | Enforced engagement boundary |
| Approvals | Request, approve, reject, expire | Policy-controlled authorization |
| Tasks | Create, inspect, cancel, retry | Durable lifecycle |
| Workers | Enroll, capability, health, quarantine, revoke | Authenticated worker control |
| Evidence | Query, verify, export | Chain-of-custody validation |
| Findings | List, triage, remediate, regress | Assessment intelligence |
| Reports | Generate, inspect, download | Technical and executive output |
| Recovery | Backup, restore, rollback, teardown | Verified environment state |

### 5.3 Contract requirements

Every API request and response requires a schema identifier, version, request identifier, actor, timestamp, and error shape. Unknown fields are rejected where the contract is closed. Duplicate requests are idempotent where state is created. Error responses identify policy, validation, dependency, timeout, or internal failure categories.

## 6. Frontend Angular application

### 6.1 Frontend structure

```text
apps/frontend-angular/src/app/
├── core/
│   ├── auth/auth.service.ts
│   ├── api/api-client.ts
│   ├── guards/auth.guard.ts
│   ├── guards/scope.guard.ts
│   └── interceptors/request.interceptor.ts
├── features/
│   ├── dashboard/
│   │   ├── dashboard.component.ts
│   │   ├── dashboard.component.html
│   │   └── dashboard.component.scss
│   ├── assets/
│   ├── scope/
│   ├── approvals/
│   ├── workers/
│   ├── tasks/
│   ├── findings/
│   ├── evidence/
│   ├── reports/
│   ├── telemetry/
│   └── recovery/
├── models/
│   ├── asset.model.ts
│   ├── task.model.ts
│   ├── finding.model.ts
│   ├── evidence.model.ts
│   └── report.model.ts
├── app.config.ts
├── app.routes.ts
├── app.ts
├── app.html
└── app.scss
```

### 6.2 Frontend workflow

```text
LOGIN
  ▼
DASHBOARD HEALTH
  ▼
SELECT REGISTERED ASSET
  ▼
SELECT APPROVED WORKFLOW
  ▼
REVIEW SCOPE AND APPROVAL
  ▼
SUBMIT TASK
  ▼
FOLLOW WORKER AND TASK STATE
  ▼
REVIEW FINDINGS AND EVIDENCE
  ▼
GENERATE REPORT
  ▼
REVIEW RECOVERY AND TEARDOWN
```

### 6.3 Frontend acceptance

- The UI reads registered assets from the gateway.
- The UI never accepts an unregistered destination.
- The UI displays policy denial reasons.
- The UI shows worker health and task state.
- The UI links findings to evidence.
- The UI requests report generation through the gateway.
- The UI displays backup, rollback, and teardown state.
- API failures, expired sessions, rate limits, and emergency stop are visible.

## 7. LangGraph orchestrator

### 7.1 Orchestrator structure

```text
apps/orchestrator-langgraph/
├── nodes/
│   ├── load_scope.py
│   ├── validate_approval.py
│   ├── build_plan.py
│   ├── dispatch_worker.py
│   ├── collect_result.py
│   ├── write_evidence.py
│   ├── build_finding.py
│   ├── generate_report.py
│   └── compensate.py
├── policies/
│   ├── scope_policy.py
│   ├── capability_policy.py
│   ├── retry_policy.py
│   └── emergency_policy.py
├── workflows/
│   ├── assessment_workflow.py
│   ├── recovery_workflow.py
│   └── teardown_workflow.py
├── planner.py
├── state.py
├── test_planner.py
└── requirements.txt
```

### 7.2 Workflow state

```text
scope_id
approval_id
asset_ids
workflow_id
current_node
completed_nodes
pending_nodes
worker_assignments
retry_count
risk_score
policy_version
evidence_refs
finding_refs
compensation_plan
error_state
```

### 7.3 Orchestrator acceptance

- Workflow nodes have typed input and output.
- A policy denial prevents downstream dispatch.
- A failed dependency prevents dependent nodes from running.
- Retry budgets are finite and recorded.
- Cancellation propagates to workers.
- Approval pauses preserve state.
- Compensation steps restore affected lab state.
- Every decision has an explanation and policy reference.

## 8. Lab infrastructure and deployment

### 8.1 Lab service topology

```text
                    ┌──────────────────┐
                    │ Operator Console │
                    └────────┬─────────┘
                             │
                    ┌────────▼─────────┐
                    │ .NET API Gateway │
                    └────────┬─────────┘
                             │
                    ┌────────▼─────────┐
                    │ Go Control Plane │
                    └───┬────────┬─────┘
                        │        │
             ┌──────────▼─┐  ┌───▼──────────┐
             │ Orchestrator│  │ Evidence     │
             │ LangGraph   │  │ and Reports  │
             └──────┬──────┘  └──────┬──────┘
                    │                │
          ┌─────────▼────────────────▼─────────┐
          │ Authenticated Assessment Workers   │
          └───┬──────────┬──────────┬──────────┘
              │          │          │
       ┌──────▼───┐ ┌────▼─────┐ ┌──▼─────────┐
       │ Web/API   │ │ Databases│ │ Identity   │
       │ Services  │ │ Services │ │ Services   │
       └──────┬───┘ └────┬─────┘ └──┬─────────┘
              │          │          │
              └──────────▼──────────┘
                    Telemetry Layer
```

### 8.2 Deployment files

```text
deploy/
├── terraform/
│   ├── main.tf
│   ├── variables.tf
│   ├── outputs.tf
│   ├── modules/vpc/
│   ├── modules/control_plane/
│   ├── modules/assessment_workers/
│   ├── modules/data_services/
│   ├── modules/telemetry/
│   └── environments/lab/
├── ansible/
│   ├── inventories/lab/hosts.ini
│   ├── group_vars/lab.yml
│   └── playbooks/
│       ├── base.yml
│       ├── control_plane.yml
│       ├── workers.yml
│       ├── data_services.yml
│       ├── telemetry.yml
│       └── recovery.yml
├── compose/
│   ├── control-plane.yml
│   ├── lab-services.yml
│   └── observability.yml
├── nginx/
│   └── nginx.conf
└── scripts/
    ├── deploy.sh
    ├── validate.sh
    ├── backup.sh
    ├── restore.sh
    ├── rollback.sh
    └── destroy.sh
```

### 8.3 Deployment acceptance

- Terraform syntax and provider lock validation pass.
- Ansible syntax and inventory validation pass.
- Container images build with pinned dependencies.
- Private network segmentation is verified.
- Health checks pass after deployment.
- Backup and restore produce matching integrity records.
- Rollback restores the previous service version.
- Destroy removes all lab resources and records the result.

## 9. Domain implementation specifications

## 9.1 Layer 1: Infrastructure

### Operational replacement

The infrastructure domain becomes a reproducible private assessment environment with control, worker, data, identity, telemetry, and observability separation.

### File structure

```text
deploy/terraform/modules/
├── vpc/
├── control_plane/
├── gateway/
├── workers/
├── web_services/
├── databases/
├── identity/
├── telemetry/
├── evidence/
└── reports/
```

### Input, process, output, evidence

| Input | Process | Output | Evidence |
|---|---|---|---|
| Environment config, service versions, network policy | Provision and configure | Running lab and asset registry | Deployment manifest, health results, versions |
| Backup request | Snapshot and hash state | Backup artifact | Backup hash and timestamp |
| Rollback request | Restore approved version | Previous service state | Rollback verification |

### P0/P1 result

P0 requires clean provisioning, health, segmentation, backup, restore, and teardown. P1 adds controlled upgrades, drift detection, and recovery measurement.

## 9.2 Layer 2: Control plane and authenticated workers

### Operational replacement

The legacy implant/teamserver domain becomes an authenticated assessment-worker platform. Workers are registered services with declared capabilities, certificate lifecycle, task expiry, health monitoring, quarantine, result signatures, and evidence links.

### File structure

```text
src/c2/modules/assessment_workers/
├── registry.go
├── enrollment.go
├── capability.go
├── heartbeat.go
├── dispatch.go
├── result_verifier.go
├── quarantine.go
└── worker_test.go
```

### Worker lifecycle

```text
REGISTER → VERIFY IDENTITY → DECLARE CAPABILITY → HEALTHY
    │                                      │
    ├── REVOKED ← QUARANTINED ← UNHEALTHY  │
    │                                      ▼
    └────────────── RECEIVE SCOPED TASK → SIGNED RESULT
```

### P0/P1 result

P0 requires registration, scoped task delivery, signed result, expiry, health, quarantine, and revocation. P1 adds worker pools, bounded concurrency, rolling upgrades, and recovery.

## 9.3 Layer 3: Protocol and service assessment

### Operational replacement

The legacy traffic-profile and listener requirements become protocol-validation workers for HTTP, HTTPS, WebSocket, DNS service health, and gateway contract behavior against registered lab services.

### File structure

```text
src/c2/modules/assessment_workers/protocol/
├── http_posture.go
├── api_contract.go
├── websocket_health.go
├── dns_service.go
├── tls_posture.go
├── request_policy.go
└── protocol_test.go
```

### P0/P1 result

P0 validates health, authentication, schema, headers, TLS configuration, and error handling. P1 adds traffic baselines, service dependency maps, detection coverage, and regression comparison.

## 9.4 Layer 4: Deception and detection

### Operational replacement

The deception domain becomes registered decoy services, canary assets, event correlation, alert routing, and response verification.

### File structure

```text
src/c2/modules/deception/
├── decoy_service.go
├── canary_asset.go
├── event_capture.go
├── alert_correlation.go
├── containment.go
└── deception_test.go
```

### P1 result

A registered decoy access creates a correlated event, alert, evidence record, timeline entry, and reversible containment action.

## 9.5 Layer 5: Database security

### Operational replacement

The SQL and NoSQL domains become real database security assessment workers for owned relational and document services. The workers test access policy, schema exposure, audit configuration, backup integrity, configuration posture, and remediation state.

### File structure

```text
src/c2/modules/database_security/
├── relational_policy.go
├── document_policy.go
├── privilege_graph.go
├── schema_inventory.go
├── audit_correlation.go
├── backup_integrity.go
├── remediation.go
└── database_security_test.go
```

### Input, process, output, evidence

| Input | Process | Output | Evidence |
|---|---|---|---|
| Registered database ID, role, policy | Verify connection and access policy | Access decision and finding | Query audit, role, timestamp |
| Configuration snapshot | Compare to approved baseline | Drift finding | Snapshot hash and diff |
| Backup ID | Restore and verify | Recovery status | Restore log and integrity hash |
| Remediation change | Re-run control | Resolved or persistent finding | Before/after evidence |

## 9.6 Layer 6: Identity and directory security

### Operational replacement

The AD and Kerberos domain becomes an isolated identity and directory security lab with users, groups, role bindings, service identities, approval state, expiry, access review, and lifecycle verification.

### File structure

```text
src/c2/modules/identity_security/
├── account_lifecycle.go
├── role_graph.go
├── group_review.go
├── service_identity.go
├── stale_identity.go
├── access_review.go
└── identity_test.go
```

### P1 result

The platform detects excessive access, stale identities, expired approvals, mismatched role bindings, and unreviewed service identities, then records remediation and re-verifies the state.

## 9.7 Layer 7: Network segmentation and path validation

### Operational replacement

The lateral movement domain becomes a network and identity path validator. It builds a registered asset graph, evaluates allowed and denied paths, verifies segmentation, correlates denied events, and supports reversible containment.

### File structure

```text
src/c2/modules/network_security/
├── segment_graph.go
├── asset_path.go
├── allowed_path.go
├── denied_path.go
├── firewall_policy.go
├── containment.go
└── network_test.go
```

### P1 result

Every approved path succeeds only under policy. Every denied path produces a reason, event, evidence record, and alert. The graph shows the dependency and trust relationship without uncontrolled host movement.

## 9.8 Layer 8: Configuration, drift, and recovery

### Operational replacement

The persistence domain becomes configuration inventory, approved baseline comparison, drift detection, recovery, rollback, and cleanup verification.

### File structure

```text
src/c2/modules/configuration/
├── baseline.go
├── snapshot.go
├── drift.go
├── remediation.go
├── rollback.go
└── configuration_test.go
```

### P1 result

The platform detects unapproved service, schedule, package, identity, and configuration changes in registered lab assets, restores the approved baseline, and verifies the result.

## 9.9 Layer 9: Credential and secret governance

### Operational replacement

The credential domain becomes a secret-management and access-governance service using project-owned credentials. It handles issuance, access review, rotation, expiry, revocation, exposure detection in registered stores, and report redaction.

### File structure

```text
src/c2/modules/secret_governance/
├── secret_registry.go
├── access_review.go
├── exposure_detection.go
├── rotation.go
├── expiry.go
├── revocation.go
└── secret_governance_test.go
```

### P1 result

Access outside role policy is denied, exposure is detected in registered stores, rotation invalidates previous values, expiry blocks authentication, and reports contain no raw secret values.

## 9.10 Layer 10: Collection and privacy-safe evidence

### Operational replacement

The collector domain becomes a consented data-observation and evidence pipeline for registered service events, configuration records, telemetry, and test artifacts. It applies provenance, redaction, retention, and evidence-chain controls.

### File structure

```text
src/c2/modules/telemetry/
├── events.go
├── collector.go
├── correlation.go
├── redaction.go
├── retention.go
└── telemetry_test.go
```

### P0/P1 result

Every collected record has source, asset, task, timestamp, classification, redaction state, and hash. Retention and deletion produce verification records.

## 9.11 Layer 11: Impact and resilience

### Operational replacement

The destruction and impact domain becomes blast-radius analysis, backup validation, recovery drills, RTO/RPO measurement, rollback, and service dependency verification.

### File structure

```text
src/c2/modules/recovery/
├── dependency_graph.go
├── blast_radius.go
├── backup.go
├── restore.go
├── rollback.go
├── rto_rpo.go
└── recovery_test.go
```

### P1 result

The platform measures affected services, verifies backup integrity, restores registered state, records recovery duration, and produces a recovery report.

## 9.12 Layer 12: Artifact and firmware integrity

### Operational replacement

The hardware and firmware domain becomes artifact inventory, image metadata, hash verification, signature verification, secure-boot posture, release comparison, and approved recovery.

### File structure

```text
src/c2/modules/artifact_integrity/
├── inventory.go
├── metadata.go
├── hash_verify.go
├── signature_verify.go
├── diff.go
├── recovery.go
└── artifact_test.go
```

### P1 result

Unapproved artifact changes are detected, invalid signatures block promotion, differences are reviewable, and recovery returns the registered asset to its approved artifact state.

## 9.13 Layer 13: Asset intelligence and reconnaissance

### Operational replacement

The reconnaissance domain becomes registered asset intelligence: ownership, service metadata, DNS records for registered zones, dependency mapping, exposure inventory, and change detection.

### File structure

```text
src/c2/modules/asset_intelligence/
├── inventory.go
├── ownership.go
├── dns_inventory.go
├── service_metadata.go
├── dependency_map.go
├── exposure.go
└── inventory_test.go
```

### P0/P1 result

P0 builds an accurate asset registry. P1 correlates services, identities, data stores, network segments, dependencies, exposure changes, and evidence.

## 9.14 Layer 14: Controlled application assessment

### Operational replacement

The exploitation domain becomes an owned vulnerable-application assessment environment with controlled security checks, request validation, authorization testing, configuration review, patch verification, and regression evidence.

### File structure

```text
src/c2/modules/application_security/
├── authentication.go
├── authorization.go
├── input_validation.go
├── api_schema.go
├── configuration.go
├── remediation.go
└── application_security_test.go
```

### P1 result

The platform identifies controlled policy violations in registered applications, links findings to evidence, validates remediation, and prevents requests outside the registered application contract.

## 9.15 Layer 15: Orchestration and auditable decision-making

### Operational replacement

The autonomous brain becomes bounded, explainable workflow planning with risk scoring, dependency ordering, approval pauses, policy references, and compensation.

### File structure

```text
src/c2/orchestrator/
├── scope.go
├── approval.go
├── graph.go
├── risk.go
├── behavior.go
├── timing.go
├── compensation.go
└── orchestrator_test.go
```

### P1 result

Every decision is reproducible from inputs, policy version, risk score, dependencies, and approval state. The planner cannot dispatch an operation outside the registered workflow.

## 9.16 Layer 16: Evidence chain of custody

### Existing and target files

```text
src/c2/evidence/
├── ledger.go
├── manifest.go
├── redaction.go
├── verification.go
├── provenance.go
└── evidence_test.go
```

### P0 result

Evidence records are immutable by hash chain, ordered by time, redacted, linked to asset/task/finding, independently verifiable, and required for report finalization.

## 9.17 Layer 17: Reporting and intelligence

### Existing and target files

```text
src/c2/report/
├── report_generator.go
├── executive.go
├── technical.go
├── timeline.go
├── remediation.go
├── delivery.go
└── report_test.go
```

### P0/P1 result

Reports include scope, methodology, assets, findings, severity, confidence, evidence, timeline, remediation, residual risk, recovery status, and limitations.

## 9.18 Layer 18: Cleanup, rollback, and teardown

### Existing and target files

```text
src/c2/cleanup/
├── manifest.go
├── artifact_cleanup.go
├── rollback.go
├── verification.go
└── cleanup_test.go
```

### P0 result

Every state-changing workflow produces a cleanup or compensation plan, executes the approved recovery sequence, and verifies the resulting state.

## 9.19 Layer 19: API gateway and operator access

The gateway controls authentication, authorization, request validation, rate limiting, error handling, audit, task persistence, and report access. The Angular console consumes only typed gateway contracts.

## 9.20 Layer 20: Database and state persistence

The state layer stores assets, scopes, approvals, workers, tasks, results, findings, evidence references, reports, deployments, backups, and recovery records. Migrations are versioned, rollback-aware, and covered by tests.

## 9.21 Layer 21: Telemetry and detection

Telemetry records task dispatch, worker health, policy decisions, service requests, denials, evidence writes, finding changes, report generation, backup, restore, rollback, and teardown. Detection rules have regression tests and coverage metrics.

## 9.22 Layer 22: CI and quality automation

```text
make fmt
make test
make vet
make build
make safe-contracts
make app-check
make lab-check
make recovery-check
make release-check
```

The CI pipeline blocks release on formatting, build, contract, unit, integration, lab, evidence, report, recovery, dependency, or repository-integrity failure.

## 9.23 Layer 23: Deployment and operations

Operations documentation must cover environment setup, identity provisioning, secret configuration, service startup, health checks, scope registration, task creation, evidence review, report generation, backup, restore, rollback, teardown, and incident response.

## 9.24 Layer 24: P0/P1 acceptance matrix

| ID | Acceptance test | Required result | Priority |
|---|---|---|---|
| P0-01 | Clean lab provision | All services healthy and assets registered | Critical |
| P0-02 | Scope denial | Unregistered asset rejected with evidence | Critical |
| P0-03 | Approval denial | Unapproved task rejected with policy reason | Critical |
| P0-04 | Worker enrollment | Registered worker healthy and capable | Critical |
| P0-05 | Task lifecycle | Create through result with valid state transitions | Critical |
| P0-06 | Failure lifecycle | Timeout, retry, failure, and compensation verified | Critical |
| P0-07 | Evidence tamper | Modified record fails chain verification | Critical |
| P0-08 | Report generation | Report includes findings, evidence, timeline, remediation | Critical |
| P0-09 | Gateway workflow | Angular to gateway to control plane succeeds | Critical |
| P0-10 | Backup restore | State restored with matching integrity record | Critical |
| P0-11 | Teardown | Lab resources removed and clean state verified | Critical |
| P1-01 | Web/API assessment | Registered services produce findings and remediation evidence | Hard |
| P1-02 | Database assessment | Relational and document policy checks pass | Hard |
| P1-03 | Identity assessment | Role, expiry, stale access, and lifecycle findings pass | Hard |
| P1-04 | Segmentation | Allowed and denied paths behave according to policy | Hard |
| P1-05 | Deception | Canary event produces correlated alert and timeline | Hard |
| P1-06 | Telemetry | Event coverage and tamper detection pass | Hard |
| P1-07 | Orchestration | Multi-step workflow handles approval, failure, retry, and compensation | Hard |
| P1-08 | Secret governance | Rotation, expiry, revocation, and redaction pass | Hard |
| P1-09 | Drift recovery | State is restored to approved baseline | Hard |
| P1-10 | Artifact integrity | Hash, signature, diff, and recovery checks pass | Hard |
| P1-11 | Resilience | RTO/RPO, backup, recovery, and rollback are measured | Hard |
| P1-12 | Release | Clean checkout passes complete release suite | Hard |

## 9.25 Layer 25: Delivery and completion record

Each release must include:

```text
release_id
commit_sha
service_versions
toolchain_versions
asset_registry_hash
scope_id
approval_ids
task_ids
worker_versions
evidence_root_hash
finding_summary
report_hash
backup_hash
recovery_result
teardown_result
unit_test_result
integration_test_result
lab_acceptance_result
security_review_result
operator_signoff
```

A P0 release is complete when every P0 acceptance row passes. A P1 release is complete when every P0 row remains green and every P1 row passes. The release record is stored outside generated build output and linked to the evidence chain.

## 10. Implementation flow from zero to release

```text
[START]
   │
   ▼
[Review contracts and repository status]
   │
   ▼
[Create lab network and service registry]
   │
   ▼
[Deploy gateway, control plane, identity, data, telemetry]
   │
   ▼
[Run health, segmentation, and dependency tests]
   │
   ▼
[Implement scope, approval, task, worker, and evidence contracts]
   │
   ▼
[Connect Angular, .NET, Go, and LangGraph workflows]
   │
   ▼
[Implement P0 assessment workers]
   │
   ▼
[Run P0 unit, contract, integration, lab, recovery, and teardown tests]
   │
   ▼
[Implement P1 web, database, identity, network, telemetry, recovery, and artifact domains]
   │
   ▼
[Run P1 failure, resilience, remediation, and regression tests]
   │
   ▼
[Generate evidence package and reports]
   │
   ▼
[Run clean-checkout release gate]
   │
   ▼
[Record release and operator signoff]
   │
   ▼
[END]
```

## 11. Current repository alignment

The current repository already contains reusable foundations in `src/c2/orchestrator`, `src/c2/modules/orchestrator`, `src/c2/evidence`, `src/c2/report`, `src/c2/cleanup`, `src/c2/server`, `src/c2/osint`, `src/c2/simulation/replay`, `apps/frontend-angular`, `apps/gateway-dotnet`, `apps/orchestrator-langgraph`, `contracts`, `simulation`, `deploy`, and `scripts`.

The required implementation work is to connect these foundations into the target repository structure, add persistent control-plane state, add authenticated workers, add registered lab services, connect the Angular client to the gateway, run deployment acceptance, replace false-success paths in promoted packages, add integration and recovery suites, and synchronize the status documents.

The current GitHub source tree has 230 Go files and 18 Go test files. The new blueprint requires every promoted P0/P1 package to gain direct tests, integration coverage, failure handling, evidence linkage, and lab acceptance before it can be marked released.

## 12. Final release requirements

The platform is accepted as complete only when:

1. The lab provisions from a clean checkout.
2. Every service reports health and version.
3. Registered assets and scopes are enforced.
4. Operators and workers are authenticated.
5. Approvals and emergency stop work.
6. Tasks are durable, idempotent, cancellable, and recoverable.
7. Workers execute registered assessment operations and return verified results.
8. Evidence is chained, redacted, and independently verified.
9. Findings and reports are generated from actual records.
10. The Angular console completes the operator workflow.
11. Web, API, database, identity, network, telemetry, recovery, and artifact checks pass their P1 suites.
12. Backup, restore, rollback, and teardown are verified.
13. CI blocks release on any required failure.
14. The release package contains deployment, test, evidence, report, recovery, and teardown records.

## References

[1]: https://github.com/angelrachel/ANGEL "ANGEL repository"
[2]: https://owasp.org/www-project-application-security-verification-standard/ "OWASP Application Security Verification Standard"
[3]: https://owasp.org/www-project-web-security-testing-guide/ "OWASP Web Security Testing Guide"
[4]: https://www.nist.gov/cyberframework "NIST Cybersecurity Framework"
[5]: https://json-schema.org/ "JSON Schema documentation"
