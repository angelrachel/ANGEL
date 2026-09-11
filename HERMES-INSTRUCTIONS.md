# HERMES-INSTRUCTIONS

## Instruksi Eksekusi ANGEL dari Nol sampai Verifikasi

Gunakan `/home/ubuntu/ANGEL.md` sebagai blueprint utama dan jangan gunakan PDF sebagai instruksi eksekusi langsung. PDF/source hanya dipakai untuk traceability internal.

## 1. Tujuan

Bangun repository ANGEL dari keadaan kosong sampai menjadi platform authorized security validation yang dapat dijalankan lokal, diuji dalam controlled lab, menghasilkan evidence, finding P2/P1/P0 bila bukti memenuhi threshold, report, dan retest.

Jangan membuat severity secara artifisial. P0/P1/P2 hanya boleh dihasilkan dari evidence valid, exposure/reachability, criticality, impact, confidence, dan review sesuai ANGEL.md.

## 2. Aturan kerja wajib

- Gunakan bahasa Indonesia untuk dokumentasi, komentar penjelasan, dan laporan progres.
- Baca seluruh `ANGEL.md` sebelum membuat kode.
- Baca seluruh repository dan konfigurasi sebelum mengubah file.
- Jangan membuat folder kosong, stub, skeleton, dummy success, atau endpoint hardcoded.
- Jangan menonaktifkan test, lint, type-check, security check, atau lab assertion.
- Jangan menghapus requirement sulit dari traceability.
- Semua perubahan harus memiliki test dan alasan teknis.
- Jika dependency eksternal tidak tersedia, selesaikan jalur offline/lab dan catat blocker.
- Jangan memasukkan secret, credential, token, data klien, atau target nyata.

## 3. Batas kemampuan

Implementasikan authorized worker, typed action, policy engine, posture assessment, evidence, graph analysis, detection validation, recovery validation, reporting, dan retest.

Jangan membuat atau mengaktifkan malware, ransomware, rootkit, backdoor, credential extraction, keylogger, hidden persistence, evasion, process injection, arbitrary command execution, destructive payload, covert channel, data exfiltration, atau implant/C2 aktif.

Untuk requirement berisiko, pertahankan tujuan evaluasi melalui posture analysis, synthetic simulation, benign marker, controlled lab, detection validation, recovery validation, atau canary tanpa data nyata. Tandai sebagai `SUBSTITUTED-SAFE`, bukan `DONE` palsu.

## 4. Artefak wajib

Buat dan rawat:

- `ANGEL.md`
- `REQUIREMENTS-TRACEABILITY.md`
- `MODULE-TRACEABILITY.csv`
- `TYPED-ACTION-CATALOG.md`
- `PERMISSION-MATRIX.md`
- `API-SPEC.md`
- `DATABASE-SCHEMA.md`
- `LAB-ARCHITECTURE.md`
- `TEST-MATRIX.md`
- `ACCEPTANCE-CRITERIA.md`
- `OPEN-DECISIONS.md`
- `VERIFICATION-REPORT.md`

## 5. Status yang diizinkan

Gunakan hanya status berikut:

`PLANNED`, `IN_PROGRESS`, `IMPLEMENTED`, `TESTED`, `INTEGRATED`, `DOCUMENTED`, `DEPLOYED`, `VERIFIED`, `SUBSTITUTED-SAFE`, `BLOCKED`, `NEEDS-HUMAN-DECISION`.

Status `VERIFIED` membutuhkan bukti artifact dan test. Blueprint, nama file, atau endpoint saja tidak cukup.

## 6. Urutan eksekusi

### Tahap 0 — Baseline

Inventaris semua file, pilih stack lokal, buat migration convention, config schema, CI command, dan `OPEN-DECISIONS.md`. Jangan memilih deployment production tanpa keputusan manusia.

### Tahap 1 — Fondasi

Implementasikan config validation, database, migration, health/readiness, structured logging, trace ID, error taxonomy, secret reference, dan local compose.

### Tahap 2 — Governance

Implementasikan organization isolation, user, RBAC, session, engagement, authorization record, scope allowlist, schedule, approval, audit event, legal hold, dan emergency stop.

### Tahap 3 — Policy dan orchestration

Implementasikan typed action registry, capability gate, scope matcher, rate/concurrency limit, signed job, mTLS worker identity, queue, scheduler, timeout, retry, cancellation, idempotency, replay protection, dan result validation.

### Tahap 4 — Evidence, finding, dan risk

Implementasikan collector, redaction, hashing, manifest, storage, observation, finding normalization, deduplication, confidence, asset criticality, impact, reachability, severity gate, remediation, dan retest.

### Tahap 5 — Domain assessment

Implementasikan asset, DNS/TLS, web/API posture, SBOM/dependency, identity/access posture, cloud/container/CI posture, endpoint posture, detection validation, logging/tamper validation, recovery validation, compliance baseline, dan graph correlation.

### Tahap 6 — Reporting dan lab

Implementasikan technical report, executive report, JSON/Markdown/HTML output, evidence references, lab fixtures, isolated network, synthetic data, expected results, reset, teardown, timeout, cancellation, tamper, restore, dan deterministic rerun.

## 7. Prosedur per requirement

Untuk setiap requirement dan module ID:

1. Baca source reference dan tujuan.
2. Tetapkan owner package/service.
3. Definisikan input, output, schema, error, permission, dependency, policy, dan mode.
4. Implementasikan handler/service nyata.
5. Hubungkan API/event/database bila relevan.
6. Buat unit test.
7. Buat contract/integration test.
8. Buat security negative test.
9. Buat lab scenario dengan synthetic fixture.
10. Jalankan setup, run, assertion, evidence check, report check, reset, dan rerun.
11. Perbaiki failure.
12. Update status dan artifact ID pada `MODULE-TRACEABILITY.csv`.

## 8. Lab wajib

Lab harus terisolasi dan memiliki:

- `compose.yaml` atau equivalent;
- target fixture web/API/identity/cloud/container/detection/recovery;
- synthetic data;
- policy fixture;
- expected findings;
- expected evidence hashes/schema;
- failure fixtures;
- reset dan teardown;
- no route to production or arbitrary external target.

Scenario minimum:

1. scope deny;
2. policy deny;
3. invalid signature;
4. replay detection;
5. timeout/cancel;
6. asset/posture finding;
7. web/API validation;
8. identity/trust graph;
9. SBOM/cloud/container posture;
10. synthetic detection marker;
11. evidence tamper;
12. recovery restore and RTO/RPO;
13. severity evidence gate;
14. report/retest;
15. emergency stop;
16. reset and deterministic rerun.

## 9. Quality gates

Jalankan secara berurutan:

```text
format → lint → type-check/compile → migration test → unit test
→ contract test → integration test → security test → lab test
→ evidence verification → report verification → reset/rerun
→ traceability audit → secret scan
```

Jangan melewati gate gagal. Jika tool tidak tersedia, catat `BLOCKED` dan alasan, jangan mengubah hasil menjadi lulus.

## 10. Verification gate 100%

`VERIFIED` hanya boleh digunakan jika:

- setiap source requirement memiliki ANGEL ID;
- setiap ANGEL ID memiliki owner dan contract;
- implementation artifact ada;
- unit, integration, security, dan lab test yang relevan lulus;
- evidence dan expected result tersedia;
- report dapat membaca evidence;
- lab dapat di-reset dan rerun;
- tidak ada secret;
- tidak ada test disabled;
- semua blocker dan human decision tercatat.

Buat `VERIFICATION-REPORT.md` dengan jumlah berikut:

- total source records;
- total module records;
- implemented;
- tested;
- integrated;
- substituted-safe;
- blocked;
- verified;
- test command dan exit code;
- lab scenario dan result;
- artifact path dan hash.

## 11. Final response

Gunakan bahasa Indonesia dan hanya nyatakan pekerjaan yang benar-benar terbukti. Jelaskan:

1. fitur yang diimplementasikan;
2. command yang dijalankan;
3. test yang lulus/gagal;
4. lab scenario yang lulus/gagal;
5. jumlah traceability;
6. status P2/P1/P0 evidence gate;
7. blocker eksternal;
8. keputusan manusia yang masih diperlukan;
9. cara menjalankan dan mereset lab.

Jangan mengatakan “selesai 100%” jika satu requirement masih `PLANNED`, `IN_PROGRESS`, `BLOCKED`, atau `NEEDS-HUMAN-DECISION`.

## Repository aktual — aturan wajib

Sebelum implementasi, baca `docs/ANGEL_IMPLEMENTATION_ALIGNMENT.md`. Repository ini memakai path canonical yang sudah ada (`src/assessment/`, `apps/`, `contracts/`, `lab/`, `simulation/`, `deploy/`); jangan membuat duplikasi `packages/`, `domains/`, `migrations/`, `tests/`, atau `deployments/` hanya untuk mencocokkan nama logical blueprint. Isi `REQUIREMENTS-TRACEABILITY.md`, `TEST-MATRIX.md`, `ACCEPTANCE-CRITERIA.md`, dan `OPEN-DECISIONS.md` berdasarkan artifact aktual. Jangan mengklaim `VERIFIED` tanpa command dan test lulus.
