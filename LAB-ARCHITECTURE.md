# ANGEL Lab Architecture

Lab adalah lingkungan isolated untuk proof berisiko dengan synthetic data dan disposable fixture.

## Prinsip

- Tidak ada route ke production atau target eksternal arbitrari.
- Semua target adalah fixture:// atau synthetic.
- Setelah scenario selesai, lab direset dan di-verifikasi tidak ada residual state.
- Semua execution melewati scope, capability, dan approval gate.

## Komponen

### Fixture
- `lab/services/fixture-http/` — HTTP fixture disposable.
- Target: web-app-01 dengan endpoint /healthz, /version, /security-posture, dan 404 path.

### Acceptance
- `lab/acceptance/p0_fixture_http.sh` — smoke test fixture HTTP.
- `lab/acceptance/p0_control_plane.sh` — smoke test control plane dengan task queue dan scope deny.

### Simulation
- `simulation/fixtures/tasks/recon-http.json` — task simulasi.
- `simulation/fixtures/reports/recon-http-result.json` — expected result.

## Scenario wajib

1. Scope deny — target di luar scope ditolak.
2. Policy deny — capability terlarang ditolak.
3. Invalid signature — job tidak valid ditolak.
4. Replay detection — nonce duplikat ditolak.
5. Timeout / cancel — execution berhenti sesuai batas.
6. Asset / posture finding — ditemukan di lab fixture.
7. Web / API validation — check berjalan di fixture.
8. Identity / trust graph — validasi terisolasi.
9. SBOM / cloud / container posture — nilai aman di lab.
10. Synthetic detection marker — marker terekam.
11. Evidence tamper — hash mismatch terdeteksi.
12. Recovery restore — reset berhasil dan state bersih.
13. Severity evidence gate — hanya evidence valid yang naik severity.
14. Report / retest — report dapat dibuat ulang dari evidence.
15. Emergency stop — engagement berhenti dan job ditolak.
16. Reset dan deterministic rerun — scenario dapat diulang dengan hasil sama.

## Reset

Reset dilakukan dengan:
- Menghentikan fixture.
- Membersihkan database/lab state.
- Memastikan tidak ada sisa evidence atau finding dari run sebelumnya.

## Sumber

- Fixture: `lab/services/fixture-http/`
- Acceptance: `lab/acceptance/`
- Simulation: `simulation/fixtures/`, `simulation/replay/`
- Alignment: `docs/ANGEL_IMPLEMENTATION_ALIGNMENT.md`
