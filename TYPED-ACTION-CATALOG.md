# ANGEL Typed Action Catalog

Daftar action/ability yang tersedia di ANGEL dengan jenis, permission, scope, dan batasan yang tercatat. Action adalah satuan kerja terkontrol yang melewati authorization, scope, capability gate, dan evidence.

## Kategori action

### 1. Passive discovery
- `surface-map` — pemetaan permukaan target yang diizinkan (hostname, URL, CIDR dalam scope).
- `dns-lookup` — resolusi DNS pasif untuk target dalam scope.
- `certificate-inspection` — inspeksi sertifikat TLS tujuan publik.

### 2. Safe bounded assessment
- `tls-assessment` — pemeriksaan konfigurasi TLS tujuan publik.
- `api-contract` — validasi kontrak API terhadap fixture/target berizin.
- `authorization-matrix` — pemeriksaan matriks otorisasi dan tenant isolation.

### 3. Evidence & validation
- `evidence-collection` — pengumpulan bukti non-destruktif.
- `dependency-inventory` — inventaris dependensi dan SBOM terbatas.
- `detection-validation` — validasi deteksi dengan marker terkontrol.
- `synthetic-canary` — uji proof dengan marker sintetis di lab.
- `lab-proof` — proof terkontrol di lab dengan fixture disposable.

## Action yang ditolak

Capability berikut ditolak di policy layer dan tidak boleh dijadikan action:
- `credential-collection`
- `credential-extraction`
- `persistence`
- `destructive-write`
- `log-deletion`
- `covert-channel`
- `process-injection`
- `evasion`
- `data-exfiltration`
- `arbitrary-command`

## Kontrak action

Setiap action harus memiliki:
1. Kontrak input/output (typed atau schema).
2. Permission, scope, dan capability gate.
3. Timeout, cancellation, dan retry bounded.
4. Redaction dan hash bila menghasilkan evidence.
5. Lab fixture dan negative test bila relevan.
6. Status di MODULE-TRACEABILITY.csv dan REQUIREMENTS-TRACEABILITY.md.

## Mapping ke canonical path

| Action | Owner path | Test path |
|---|---|---|
| surface-map, tls-assessment | `src/assessment/osint/` | `src/assessment/osint/*_test.go` |
| api-contract, authorization-matrix | `src/assessment/server/api/`, `contracts/` | `src/assessment/server/api/*_test.go` |
| evidence-collection | `src/assessment/evidence/` | `src/assessment/evidence/*_test.go` |
| dependency-inventory | `src/assessment/collectors/`, `src/assessment/normalization/` | `src/assessment/collectors/*_test.go` |
| detection-validation | `src/assessment/dispatch/`, `src/assessment/control/` | `src/assessment/dispatch/*_test.go` |
| synthetic-canary, lab-proof | `lab/`, `simulation/` | `lab/acceptance/`, `simulation/` |

## Status

| Status | Arti |
|---|---|
| IMPLEMENTED | Action memiliki contract, test, dan evidence path. |
| PLANNED | Action didokumentasikan tapi belum diimplementasi. |
| SUBSTITUTED-SAFE | Tujuan validasi direalisasi lewat safe equivalent/lab. |
| DENIED | Action dipromosikan tapi ditolak di policy. |

Catatan: daftar ini dioperasikan dari struktur repository aktual dan `docs/ANGEL_IMPLEMENTATION_ALIGNMENT.md`. Tidak ada duplikasi folder hanya untuk mencocokkan nama blueprint.
