# ANGEL

**ANGEL** adalah platform **Authorized Security Validation and Exposure Management** untuk menemukan, memvalidasi, memprioritaskan, melaporkan, dan melakukan retest terhadap risiko keamanan secara terotorisasi, aman, dapat diaudit, dan dapat direproduksi.

> **Blueprint alignment:** lihat [`docs/BLUEPRINT_ALIGNMENT.md`](docs/BLUEPRINT_ALIGNMENT.md) untuk pemetaan blueprint 25 layer, status implementasi, safe equivalent, dan capability yang sengaja ditolak.

## Prinsip operasi

ANGEL menerapkan **policy before execution**, typed actions, evidence-first, human-controlled risk, least privilege, reproducibility, fail-closed behavior, separation of duties, dan production-safe defaults. Platform ini tidak menyediakan covert malware, credential theft, persistence tersembunyi, process injection, log deletion, arbitrary command execution, destructive actions, atau data exfiltration.

Proof berisiko dilakukan hanya pada controlled lab menggunakan synthetic data, canary marker, callback marker, dan disposable fixtures. Target produksi harus tercantum dalam engagement yang memiliki otorisasi tertulis, scope eksplisit, time window, capability, rate budget, reviewer, dan emergency stop.

## Cakupan 25 layer

| Layer | Domain | Status implementasi |
|---:|---|---|
| 1 | Authorization & Engagement Control | Control plane, authorization, scope, emergency stop |
| 2 | Identity, RBAC & Operator Security | Operator security dan audit boundary |
| 3 | Scope Enforcement & Policy Engine | Deny-by-default matcher dan capability gate |
| 4 | Assessment Orchestration & Execution Control | Signed lifecycle, queue, timeout, cancellation |
| 5 | Platform Configuration & Secrets | Environment validation dan redacted configuration |
| 6–10 | Asset, passive recon, network, web, cloud exposure | Safe collectors dan typed adapters |
| 11–15 | Authentication, authorization, input, server-side, API logic | Non-destructive checks dan fixture-backed proof |
| 16–20 | Identity posture, cloud IAM, endpoint, detection, supply chain | Read-only posture dan control validation |
| 21–22 | Controlled lab dan safe proof | Disposable HTTP fixture, synthetic markers, reset flow |
| 23 | Evidence ledger & chain of custody | Hash chain, redaction, manifest, replay verification |
| 24 | Attack path, risk & P0/P1 engine | Deterministic risk scoring berbasis evidence dan impact |
| 25 | Reporting, remediation & retest | JSON/Markdown report, remediation lifecycle, retest |

## Repository map

- `src/assessment/` — domain engine, governance, policy, orchestration, checks, evidence, risk, reports, API, persistence, dan simulation.
- `contracts/` — OpenAPI dan JSON Schema untuk authorization, task lifecycle, evidence, cleanup, dan reporting.
- `apps/frontend-angular/` — dashboard operator untuk engagement, job, module catalog, findings, dan readiness.
- `apps/gateway-dotnet/` — gateway adapter dengan contract tests.
- `apps/orchestrator-langgraph/` — planner terkontrol untuk typed assessment workflow.
- `lab/` — fixture HTTP disposable dan acceptance tests.
- `deploy/` — deployment, validation, dan infrastructure manifests.
- `simulation/` — synthetic task/result fixtures untuk replay dan negative tests.
- `src/assessment/agent/` — lab agent simulator fixture-only dengan typed task/result dan digest evidence.
- `contracts/task/lab-agent.schema.json` — contract untuk task simulator yang menolak target produksi dan arbitrary command.
- `docs/BLUEPRINT_ALIGNMENT.md` — matriks alignment blueprint 25 layer dan definition of done.

## Arsitektur yang didukung

ANGEL mempertahankan pembagian komponen yang sejalan dengan blueprint—frontend Angular, gateway .NET 10, orchestrator LangGraph, brain, infrastructure manifests, evidence, reporting, dan lab—dengan batas capability yang aman. Istilah seperti **agent**, **task runner**, dan **deception lab** berarti komponen terkontrol untuk assessment yang diotorisasi; istilah tersebut tidak berarti implant tersembunyi, covert C2, persistence, credential theft, evasion, arbitrary command execution, data exfiltration, atau destructive payload.

Komponen yang berbahaya untuk operasi produksi diganti dengan disposable fixture, synthetic marker, typed adapter, dan controlled callback. Semua execution tetap melewati authorization, scope, time window, rate budget, approval, audit, dan emergency stop.

## Jalankan lokal

```bash
cp .env.example .env
# isi secret development lokal; jangan commit file .env
make validate
make app-check
```

Mode container:

```bash
docker compose up -d
```

Container ANGEL berjalan non-root, read-only, tanpa capability tambahan, dengan `no-new-privileges`. Port service di-bind ke loopback pada mode lokal.

## Quality gates

```bash
make fmt
make test
make vet
make build
make safe-contracts
python3 scripts/repo_integrity_check.py
```

Test negatif wajib mencakup target di luar scope, expired job, signature invalid, replay, duplicate nonce, missing approval, capability terlarang, secret leakage, tenant isolation, dan path traversal. Finding Critical/P0/P1 hanya sah apabila evidence, confidence, reachability, business criticality, dan blast radius memenuhi policy; sistem tidak menaikkan severity untuk memenuhi target jumlah temuan.

## Otorisasi

Gunakan ANGEL hanya pada target yang dimiliki atau memiliki izin tertulis. Evidence, token, hostname customer, credential, private key, dan data pribadi tidak boleh dimasukkan ke repository. Lihat [`SECURITY.md`](SECURITY.md) untuk prosedur pelaporan keamanan.
