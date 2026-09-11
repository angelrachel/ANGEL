# ANGEL Phased Implementation Plan

**Tanggal:** 2026-09-11  
**Berdasarkan:** IMPLEMENTATION-GAP-REPORT.md  
**Prinsip:** Workstream kecil 3-8 item, owner/test/contract per fase, tidak ada 609 file terpisah.

---

## Overview

Implementasi dibagi menjadi **6 fase** berdasarkan dependency graph dan critical path:

| Fase | Nama | Fokus | Estimasi | Status |
|---|---|---|---|---|
| 0 | Audit & Plan | Dokumentasi gap dan rencana | 1 hari | **DONE** |
| 1 | Core Runtime | main.go + database + job execution | 3-5 hari | PLANNED |
| 2 | Assessment Pipeline | OSINT + checks + execution | 4-6 hari | PLANNED |
| 3 | Lab & Evidence | Lab services + evidence pipeline | 3-4 hari | PLANNED |
| 4 | Frontend & Gateway | Angular dashboard + .NET gateway | 5-7 hari | PLANNED |
| 5 | Integration & E2E | End-to-end testing + deployment | 3-4 hari | PLANNED |

Total estimasi: **19-27 hari kerja**

---

## Fase 1 — Core Runtime

### Goal
Platform bisa start, persist state, mengeksekusi job sederhana end-to-end.

### Workstream (6 items)

| # | Item | Owner Path | Test | Dependencies |
|---|---|---|---|---|
| 1 | `main.go` entry point | `src/assessment/main.go` | integration test | - |
| 2 | SQLite database layer | `src/assessment/server/database/` | `*_test.go` | - |
| 3 | Database migrations | `src/assessment/server/database/migrations/` | `*_test.go` | #2 |
| 4 | Job execution pipeline | `src/assessment/server/task/executor.go` | `*_test.go` | #1, #3 |
| 5 | State persistence (FileStore) | `src/assessment/storage/state_store.go` | existing | #3 |
| 6 | HTTP server bootstrap | `src/assessment/server/listener/` + `server/api/` | `*_test.go` | #1, #4 |

### Acceptance Criteria
- [ ] `go run ./src/assessment` starts server on :8001
- [ ] `GET /healthz` returns 200
- [ ] `POST /v1/engagements` creates engagement persisted to SQLite
- [ ] `POST /v1/jobs` creates signed job and queues for execution
- [ ] State survives restart (evidence, jobs, audit)
- [ ] All new code has unit + integration tests

### Blocker
- **OPEN-DECISIONS #1:** SQLite untuk local/lab (sudah diputuskan, tetap SQLite)

### Estimasi
- 3-5 hari kerja
- ~8-12 file baru/diubah
- ~800-1200 LOC

---

## Fase 2 — Assessment Pipeline

### Goal
Platform bisa menjalankan assessment checks (OSINT, web, TLS) dan menyimpan evidence.

### Workstream (7 items)

| # | Item | Owner Path | Test | Dependencies |
|---|---|---|---|---|
| 1 | Check registry + catalog | `src/assessment/checks/` | existing + expand | - |
| 2 | OSINT collectors | `src/assessment/osint/` | `*_test.go` | - |
| 3 | HTTP/TLS collectors | `src/assessment/collectors/` | `*_test.go` | - |
| 4 | Execution adapters | `src/assessment/execution/` | `*_test.go` | #1, #2, #3 |
| 5 | Evidence pipeline | `src/assessment/evidence/` | expand | #4 |
| 6 | Risk scoring integration | `src/assessment/risk/` + `correlation/` | expand | #5 |
| 7 | Finding generation | `src/assessment/query/findings.go` | `*_test.go` | #6 |

### Acceptance Criteria
- [ ] `POST /v1/executions` runs check against fixture target
- [ ] Evidence generated with hash chain and signature
- [ ] Findings generated from evidence with severity
- [ ] Risk score deterministic based on evidence
- [ ] All checks reject non-fixture targets

### Blocker
- Tidak ada blocker teknis
- Perlu fixture data baru di `simulation/`

### Estimasi
- 4-6 hari kerja
- ~10-15 file baru/diubah
- ~1000-1500 LOC

---

## Fase 3 — Lab & Evidence

### Goal
Platform memiliki lab infrastructure yang lengkap untuk detection validation dan recovery validation.

### Workstream (5 items)

| # | Item | Owner Path | Test | Dependencies |
|---|---|---|---|---|
| 1 | Detection validation lab | `lab/services/detection/` | `*_test.go` | - |
| 2 | Recovery validation lab | `lab/services/recovery/` | `*_test.go` | - |
| 3 | Simulation fixtures | `simulation/fixtures/` | `*_test.go` | - |
| 4 | Replay verification | `src/assessment/simulation/replay/` | `*_test.go` | #1, #2, #3 |
| 5 | Lab acceptance scripts | `lab/acceptance/` | shell scripts | #4 |

### Acceptance Criteria
- [ ] Detection validation: generate alert, verify SIEM receives it
- [ ] Recovery validation: corrupt data, restore from backup, verify integrity
- [ ] Simulation fixtures: replay produces same findings
- [ ] Lab reset: `make lab-reset` cleans all state
- [ ] All lab services run in isolated containers

### Blocker
- Perlu docker-compose tambahan untuk lab services

### Estimasi
- 3-4 hari kerja
- ~8-10 file baru/diubah
- ~600-900 LOC

---

## Fase 4 — Frontend & Gateway

### Goal
Operator bisa menggunakan dashboard untuk mengelola engagements, jobs, findings.

### Workstream (6 items)

| # | Item | Owner Path | Test | Dependencies |
|---|---|---|---|---|
| 1 | Angular app shell | `apps/frontend-angular/src/app/` | `*.spec.ts` | - |
| 2 | Engagement management UI | `apps/frontend-angular/src/app/engagements/` | `*.spec.ts` | #1 |
| 3 | Job monitoring UI | `apps/frontend-angular/src/app/jobs/` | `*.spec.ts` | #1 |
| 4 | Findings dashboard | `apps/frontend-angular/src/app/findings/` | `*.spec.ts` | #1 |
| 5 | .NET Gateway expansion | `apps/gateway-dotnet/` | `*_test.go` | - |
| 6 | Gateway integration tests | `apps/gateway-dotnet.tests/` | existing | #5 |

### Acceptance Criteria
- [ ] `npm run build` succeeds
- [ ] Dashboard shows engagement list
- [ ] Operator can create engagement via UI
- [ ] Job status updates in real-time
- [ ] Findings displayed with severity badges
- [ ] Gateway passes all contract tests

### Blocker
- Tidak ada blocker teknis

### Estimasi
- 5-7 hari kerja
- ~15-20 file baru/diubah
- ~1500-2000 LOC

---

## Fase 5 — Integration & E2E

### Goal
Semua komponen terintegrasi dan platform bisa digunakan end-to-end.

### Workstream (5 items)

| # | Item | Owner Path | Test | Dependencies |
|---|---|---|---|---|
| 1 | E2E test suite | `tests/e2e/` | test files | All phases |
| 2 | CI pipeline expansion | `.github/workflows/ci.yml` | - | #1 |
| 3 | Production deployment docs | `docs/DEPLOYMENT.md` | - | - |
| 4 | Operational runbooks | `docs/RUNBOOKS.md` | - | - |
| 5 | Performance testing | `tests/perf/` | test files | All |

### Acceptance Criteria
- [ ] E2E test: create engagement → run assessment → view findings → retest
- [ ] CI runs all tests (Go, Python, Node, .NET)
- [ ] Deployment documentation complete
- [ ] Runbooks for bootstrap, emergency stop, backup/restore
- [ ] Performance baseline established

### Blocker
- OPEN-DECISIONS #2-8 harus diselesaikan sebelum production deployment

### Estimasi
- 3-4 hari kerja
- ~8-12 file baru/diubah
- ~500-800 LOC

---

## Urutan Implementasi (Dependency Order)

```
Fase 1.1: main.go
    │
    ├── Fase 1.2: SQLite layer
    │       │
    │       └── Fase 1.3: Migrations
    │               │
    │               └── Fase 1.4: Job executor
    │                       │
    │                       └── Fase 1.5: State store
    │                               │
    │                               └── Fase 1.6: HTTP bootstrap
    │                                       │
    │                                       ▼
Fase 2: Assessment Pipeline ─────────────────┤
    │                                       │
    ▼                                       │
Fase 3: Lab & Evidence ──────────────────────┤
    │                                       │
    ▼                                       │
Fase 4: Frontend & Gateway ──────────────────┘
    │
    ▼
Fase 5: Integration & E2E
```

---

## Test Strategy per Fase

| Fase | Unit Test | Integration Test | Contract Test | Lab Test |
|---|---|---|---|---|
| 1 | ✅ Required | ✅ Required | - | - |
| 2 | ✅ Required | ✅ Required | ✅ Required | ✅ Required |
| 3 | ✅ Required | - | - | ✅ Required |
| 4 | ✅ Required | - | ✅ Required | - |
| 5 | - | ✅ Required | ✅ Required | ✅ Required |

---

## File Changes per Fase

| Fase | New Files | Modified Files | Total |
|---|---|---|---|
| 1 | 6-8 | 4-6 | ~12 |
| 2 | 8-12 | 6-8 | ~18 |
| 3 | 6-8 | 4-6 | ~12 |
| 4 | 12-16 | 6-8 | ~22 |
| 5 | 6-10 | 4-6 | ~14 |
| **Total** | **38-54** | **24-34** | **~78** |

---

## Lab Scenarios Required

| Scenario | Phase | Service |
|---|---|---|
| Fixture HTTP assessment | 2 | `lab/services/fixture-http` |
| Detection validation | 3 | `lab/services/detection` |
| Recovery validation | 3 | `lab/services/recovery` |
| Canary callback | 3 | `lab/services/canary` |
| Full engagement E2E | 5 | All services |

---

## Acceptance Criteria Summary

Platform dianggap **selesai per fase** ketika:

1. Semua test hijau (`make check` passes)
2. Semua acceptance criteria terpenuhi
3. Traceability terupdate (MODULE-TRACEABILITY.csv)
4. Tidak ada TODO/FIXME/dummy/stub
5. Tidak ada capability DENIED yang diaktifkan
6. Dokumentasi terupdate

---

## Blocker & Decisions Needing Human

| Item | Type | Status |
|---|---|---|
| Database production (SQLite vs Postgres) | Decision | OPEN-DECISIONS #1 (default SQLite OK) |
| Queue production (in-memory vs broker) | Decision | OPEN-DECISIONS #2 |
| Object storage (local vs cloud) | Decision | OPEN-DECISIONS #3 |
| Identity provider + MFA | Decision | OPEN-DECISIONS #4 |
| Retention/legal hold | Decision | OPEN-DECISIONS #5 |
| Connector cloud/SIEM | Decision | OPEN-DECISIONS #6 |
| Production SLO/RTO/RPO | Decision | OPEN-DECISIONS #7 |
| Deployment target + CA | Decision | OPEN-DECISIONS #8 |

---

## Notes

- Setiap fase maksimal 8 item workstream
- Item dengan owner/test/contract yang sama digabung
- 609 capability tidak diperlakukan sebagai 609 file terpisah
- Blueprint di-group menjadi ~78 file aktual
- Safe equivalent dijaga sesuai SAFE_EQUIVALENT_ROADMAP.md

---

**Status:** FASE 0 COMPLETE — Menunggu instruksi Fase 1
