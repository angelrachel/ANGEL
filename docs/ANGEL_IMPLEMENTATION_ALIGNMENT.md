# ANGEL Implementation Alignment

Dokumen ini adalah jembatan operasional antara blueprint ANGEL dan struktur repository yang sudah ada. **Struktur repository aktual adalah canonical implementation path**; nama path logis pada blueprint dipetakan ke path aktual di sini. Hermes wajib membaca dokumen ini sebelum memindahkan atau membuat file.

## Keputusan struktur

Repository ini mempertahankan monorepo aktual: `src/assessment/` untuk core Go, `apps/` untuk frontend/gateway/orchestrator, `contracts/` untuk schema, `lab/` untuk fixture dan acceptance, `simulation/` untuk fixture/replay, `deploy/` untuk deployment, dan `scripts/` untuk quality gates. Jangan membuat duplikasi `packages/`, `domains/`, `migrations/`, `tests/`, atau `deployments/` hanya untuk mencocokkan nama blueprint. Gunakan path aktual dan catat mapping.

## Canonical map

| Blueprint logical area | Repository canonical path |
|---|---|
| apps/api | `src/assessment/server/api/` dan `apps/gateway-dotnet/` |
| apps/worker | `src/assessment/agent/`, `src/assessment/execution/`, `src/assessment/server/task/` |
| apps/scheduler | `src/assessment/scheduler/` |
| apps/report-renderer | `src/assessment/report/`, `src/assessment/reporting/` |
| packages/identity | `src/assessment/governance/`, `src/assessment/domain/` |
| packages/governance | `src/assessment/governance/`, `src/assessment/control/` |
| packages/policy-engine | `src/assessment/policy/` |
| packages/evidence | `src/assessment/evidence/` |
| packages/risk-engine | `src/assessment/risk/`, `src/assessment/correlation/` |
| packages/storage | `src/assessment/storage/`, `src/assessment/server/database/`, `src/assessment/persistence/` |
| domains/asset | `src/assessment/assets/` |
| domains/network-posture | `src/assessment/osint/` dan `src/assessment/checks/` |
| domains/web-posture | `src/assessment/osint/web.go` dan `src/assessment/checks/` |
| domains/api-posture | `src/assessment/server/api/`, `contracts/` |
| domains/recovery | `src/assessment/persistence/`, `src/assessment/cleanup/` |
| controlled lab | `lab/`, `simulation/` |
| migrations | `src/assessment/server/database/` dan migration mechanism yang harus ditetapkan |
| tests | test files berdampingan dengan package, `apps/*tests`, `lab/acceptance/` |
| deployments | `deploy/` dan `docker-compose.yml` |

## Status interpretation

- `Implemented`: artifact dan test sudah ada, tetapi tetap perlu command verification.
- `Partial`: sebagian contract atau domain ada; jangan klaim coverage penuh.
- `Safe equivalent`: tujuan validasi tersedia dengan batas aman/lab.
- `Planned`: belum ada artifact implementasi.
- `Denied`: capability tidak boleh ditambahkan.

## Rule Hermes

1. Pertahankan path aktual bila sudah memenuhi boundary.
2. Buat file baru hanya pada owner path yang tercantum.
3. Jangan memindahkan 154+ file Go atau membuat duplicate package tanpa alasan.
4. Setiap perubahan update traceability, test matrix, dan alignment.
5. Test harus menunjukkan status nyata; keberadaan file bukan bukti implementasi.
6. Semua capability berbahaya tetap `Denied` atau `Safe equivalent`.
