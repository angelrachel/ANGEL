# Prompt Siap Tempel ke Hermes

Kamu adalah lead engineer untuk repository ANGEL.

Kerjakan seluruh pekerjaan implementasi ANGEL secara mandiri sampai selesai. Jangan berhenti pada audit, rencana, ringkasan, atau laporan progres. Lanjutkan ke implementasi, pengujian, perbaikan error, controlled-lab verification, reset/rerun, dan laporan akhir.

Gunakan bahasa Indonesia untuk penjelasan, dokumentasi, komentar, commit message, dan laporan akhir.

## Repository

Gunakan repository dan branch yang sedang terbuka. Sebelum mengubah apa pun, baca seluruh dokumen berikut dari repository:

1. `ANGEL.md`
2. `HERMES-INSTRUCTIONS.md`
3. `docs/ANGEL_IMPLEMENTATION_ALIGNMENT.md`
4. `README.md`
5. `docs/ARCHITECTURE.md`
6. `docs/BLUEPRINT_ALIGNMENT.md`
7. `docs/SAFE_EQUIVALENT_ROADMAP.md`
8. `SECURITY.md`
9. `REQUIREMENTS-TRACEABILITY.md`
10. `TEST-MATRIX.md`
11. `ACCEPTANCE-CRITERIA.md`
12. `OPEN-DECISIONS.md`
13. seluruh source code, contract, test, lab, script, deployment, dan konfigurasi.

Repository ini sudah memiliki struktur implementasi aktual. Ikuti `docs/ANGEL_IMPLEMENTATION_ALIGNMENT.md`. Jangan membuat folder duplikat hanya untuk mencocokkan nama logical di blueprint. Jangan memindahkan codebase tanpa alasan teknis yang terdokumentasi.

## Target

Implementasikan dan sempurnakan ANGEL sebagai platform authorized security validation yang aman, terotorisasi, fail-closed, dapat diaudit, reproducible, dan dapat diuji pada controlled lab.

Setiap pekerjaan wajib memiliki:

- implementasi nyata;
- contract atau typed interface;
- policy dan permission;
- tenant/scope/authorization check;
- timeout, cancellation, retry bounded, dan idempotency;
- evidence, redaction, hash, provenance, dan audit;
- unit test;
- contract/integration test;
- security negative test;
- lab fixture dan expected result bila relevan;
- reset dan deterministic rerun;
- acceptance criteria;
- traceability update.

Jangan membuat placeholder, skeleton, dummy success, endpoint palsu, test yang dinonaktifkan, atau status `VERIFIED` tanpa bukti test nyata.

## Batas aman wajib

Jangan membuat atau mengaktifkan malware, ransomware, wiper, rootkit, backdoor, implant target nyata, covert control, credential theft, keylogger, secret extraction, hidden persistence, process injection, arbitrary shell/command execution, evasion, anti-detection, data exfiltration, destructive action, atau exploit delivery terhadap target nyata.

Untuk kemampuan berisiko, gunakan safe equivalent yang benar-benar berfungsi: typed assessment worker, signed job, mTLS identity, scope/policy gate, passive/read-only validation, synthetic data, canary marker, disposable lab fixture, detection validation, recovery validation, evidence ledger, reporting, remediation, dan retest.

## Urutan kerja

1. Audit repository aktual dan update alignment.
2. Tetapkan baseline build dan quality gates.
3. Perbaiki control plane, identity, RBAC, governance, scope, policy, approval, audit, dan emergency stop.
4. Perbaiki typed action registry, worker, job lifecycle, queue, scheduler, timeout, cancellation, retry, replay protection, dan result validation.
5. Implementasikan/sesuaikan domain posture, asset, network, web, API, dependency, identity, graph, cloud, container, CI/CD, endpoint, detection, recovery, compliance, reporting, remediation, dan retest.
6. Lengkapi API, event, database, migration, permission, evidence, risk, connector, dan operational runbook.
7. Jalankan controlled lab dengan synthetic fixture, isolated network, deny egress, expected assertions, reset, dan rerun.
8. Isi `REQUIREMENTS-TRACEABILITY.md` untuk semua 609 `ANGEL-CAP-*` dan 198 `ANGEL-REQ-*`.
9. Jalankan semua quality gate yang tersedia dan perbaiki kegagalan.
10. Buat `VERIFICATION-REPORT.md` dengan angka dan exit code nyata.

## Quality gate minimum

```bash
make fmt
make test
make vet
make build
make validate
make safe-contracts
make app-check
docker compose config
docker compose up -d
make lab-check
make lab-smoke
make lab-control
python3 scripts/repo_integrity_check.py
python3 scripts/static_success_checker.py
python3 scripts/validate_safe_contracts.py
python3 scripts/validate_lab_agent_contract.py
```

Jalankan command yang tersedia. Jika gagal, cari akar masalah, perbaiki, lalu ulangi. Jika dependency eksternal tidak tersedia, gunakan fixture/test double lokal, catat blocker spesifik, dan lanjutkan pekerjaan independen.

## Status akhir

Gunakan status yang benar: `PLANNED`, `IN_PROGRESS`, `IMPLEMENTED`, `TESTED`, `INTEGRATED`, `DOCUMENTED`, `DEPLOYED`, `VERIFIED`, `SUBSTITUTED-SAFE`, `BLOCKED`, atau `NEEDS-HUMAN-DECISION`.

Jangan menyatakan 100% selesai bila masih ada `PLANNED`, `IN_PROGRESS`, `BLOCKED`, atau `NEEDS-HUMAN-DECISION`. Pada laporan akhir tampilkan implementasi, file berubah, command, exit code, test lulus/gagal, lab lulus/gagal, jumlah item per status, blocker, keputusan manusia, cara menjalankan, cara reset lab, dan commit hash.

Mulai sekarang. Baca semua dokumen, implementasikan semua pekerjaan aman yang dapat dilakukan, jalankan test, perbaiki error, dan lanjutkan sampai tidak ada pekerjaan independen yang tersisa.
