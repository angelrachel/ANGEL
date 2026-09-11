# ANGEL API Specification

Sumber kebenaran utama adalah `contracts/api/openapi.yaml`. Dokumen ini adalah ringkasan operator.

## Health & observability

- `GET /healthz` — status service (tidak butuh auth).
- `GET /readyz` — readiness probe (200/503).
- `GET /metrics` — operational metrics snapshot.

## Platform & modules

- `GET /api/v1/platform` — manifest platform, layer count, dan denied capabilities.
- `GET /api/v1/modules` — daftar modul ANGEL terdaftar.

## Task & job lifecycle

- `POST /api/v1/tasks` — buat task simulasi (simulation).
- `GET /api/v1/tasks` — list task.
- `POST /api/v1/assessment-jobs` — buat signed assessment job.
- `GET /api/v1/jobs` — list job dengan filter dan pagination.

## Assessment execution

- `POST /api/v1/executions` — jalankan satu safe adapter.
- `GET /api/v1/checks` — list check yang dapat dieksekusi.
- `POST /api/v1/checks` — evaluasi satu check.

## Evidence & findings

- `GET /api/v1/evidence` — list redacted evidence.
- `GET /api/v1/findings` — list finding dengan filter.
- `GET /api/v1/observations` — query durable observations.

## Engagement & governance

- `GET /api/v1/engagements` — list policy state.
- `POST /api/v1/engagements` — register engagement.

## Remediation & retest

- `GET /api/v1/remediations` — list remediation plans.
- `POST /api/v1/remediations` — buat remediation.
- `PATCH /api/v1/remediations` — transisi status.
- `POST /api/v1/remediations/retest` — record retest.

## Reporting

- `GET /api/v1/reports` — list report.
- `POST /api/v1/reports` — generate report.
- `GET /api/v1/reports/export` — export JSON/Markdown.
- `GET /api/v1/reports/artifact` — render JSON/Markdown/HTML.

## Audit & snapshot

- `GET /api/v1/audit` — list hash-chained audit events.
- `GET /api/v1/audit/verify` — verify audit chain.
- `GET /api/v1/snapshot` — export signed snapshot.
- `POST /api/v1/snapshot/validate` — validate snapshot sebelum restore.

## Simulation

- `POST /v1/simulation/tasks` — submit fixture-backed simulation task.

## Security

Semua endpoint kecuali healthz/readyz/metrics memerlukan Bearer token via `bearerAuth`.

## Sumber

- OpenAPI: `contracts/api/openapi.yaml`
- Schema task: `contracts/task/task.schema.json`
- Schema result: `contracts/task/result.schema.json`
- Schema lifecycle: `contracts/task/lifecycle.schema.json`
- Schema lab agent: `contracts/task/lab-agent.schema.json`
