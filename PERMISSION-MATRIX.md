# ANGEL Permission Matrix

Matriks permission antara peran, konteks engagement, scope, dan action.

## Peran

| Peran | Deskripsi |
|---|---|
| `platform-admin` | Manajemen konfigurasi platform, secret reference, dan retention. |
| `owner` | Owner engagement, 승인한 scope, dan approval record. |
| `reviewer` | Review finding, severity gate, dan retest. |
| `operator` | Menjalankan assessment dengan job yang sudah disetujui. |
| `lab-agent` | Agent terkontrol di lab yang hanya berinteraksi dengan fixture. |

## Permission dan scope

Tiap permission memerlukan:
- Engagement ID yang valid dan aktif.
- Scope yang mencangkut target action.
- Capability yang diizinkan di policy.
- Time window yang tercover.
- Approval bila diperlukan.

| Action | Peran yang diizinkan | Scope requirement |
|---|---|---|
| surface-map | operator, owner | hostname/url/cidr/fixture dalam scope |
| tls-assessment | operator | tujuan publik dalam scope, pasif |
| api-contract | operator | target url/fixture dalam scope |
| authorization-matrix | operator, reviewer | tenant dalam engagement |
| evidence-collection | operator, owner | target berizin |
| dependency-inventory | operator, owner | target berizin, read-only |
| detection-validation | operator, reviewer | lab atau target berizin dengan marker |
| synthetic-canary | operator, lab-agent | lab fixture |
| lab-proof | lab-agent, operator | lab fixture dengan reset verification |

## Capability denied

| Capability | Alasan |
|---|---|
| credential-collection | Dilarang — credential theft |
| credential-extraction | Dilarang — secret extraction |
| persistence | Dilarang — hidden persistence |
| destructive-write | Dilarang — destructive action |
| log-deletion | Dilarang — log deletion |
| covert-channel | Dilarang — covert control |
| process-injection | Dilarang — process injection |
| evasion | Dilarang — evasion/anti-detection |
| data-exfiltration | Dilarang — data exfiltration |
| arbitrary-command | Dilarang — arbitrary command execution |

## Tenant isolation

- Setiap query dibatasi pada organization/engagement yang dimiliki peran.
- Cross-tenant access harus memiliki negative test.
- Lab agent hanya boleh berinteraksi dengan fixture:// target.

## Sumber kebenaran

- Policy asli: `contracts/authorization/capabilities.json`
- Schema: `contracts/authorization/capabilities.schema.json`
- Alignment: `docs/ANGEL_IMPLEMENTATION_ALIGNMENT.md`

Status: IMPLEMENTED (policy dan schema ada; test dan negative test berada di paket terkait).
