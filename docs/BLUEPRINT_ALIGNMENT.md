# ANGEL Blueprint Alignment

## Tujuan dokumen

Dokumen ini memetakan blueprint ANGEL 25 layer ke implementasi repository. Repository ini adalah **authorized security validation platform**. Semua aktivitas harus memiliki otorisasi tertulis, scope eksplisit, time window, reviewer, budget, dan emergency stop.

Blueprint awal memuat istilah offensive security yang lebih luas daripada capability yang diizinkan untuk production. Oleh karena itu, komponen yang berisiko tinggi dipetakan ke **safe equivalent** atau ditandai sebagai **out of scope**, bukan diklaim sudah diimplementasikan.

## Status implementasi

| Status | Arti |
|---|---|
| Implemented | Sudah tersedia dan diuji dalam repository. |
| Partial | Fondasi atau adapter tersedia, tetapi belum mencakup seluruh target blueprint. |
| Safe equivalent | Tujuan validasi tersedia dalam bentuk terbatas, terkontrol, atau lab-only. |
| Planned | Belum tersedia; dapat dikerjakan tanpa melewati capability boundary. |
| Denied | Sengaja tidak didukung karena dapat menjadi malware, credential theft, evasion, destructive action, atau arbitrary execution. |

## Matriks layer 1–25

| Layer blueprint | Padanan di repository | Status | Catatan |
|---:|---|---|---|
| 1 | Authorization & engagement control | Implemented | `src/assessment/governance`, `control`, dan schema authorization. |
| 2 | Identity, RBAC & operator security | Implemented | Gateway, audit boundary, operator authorization, dan tenant checks. |
| 3 | Scope enforcement & policy engine | Implemented | Deny-by-default matcher, capability gate, target/scope validation. |
| 4 | Orchestration & execution control | Implemented | Signed lifecycle, queue, timeout, cancellation, replay defense. |
| 5 | Configuration & secrets | Implemented | Environment validation, redaction, secret non-commit policy. |
| 6 | Asset inventory | Implemented | `src/assessment/assets`. |
| 7 | Passive reconnaissance | Partial | DNS, subdomain, technology, web, dan cloud collectors yang bounded. |
| 8 | Network exposure assessment | Partial | Port assessment terbatas dan rate-limited; bukan covert scanning. |
| 9 | Web exposure assessment | Partial | Safe HTTP checks dan fixture-backed validation. |
| 10 | Cloud exposure assessment | Partial | Read-only cloud posture adapters. |
| 11 | Authentication validation | Safe equivalent | Contract dan authorization checks; tidak ada brute force atau credential stuffing. |
| 12 | Authorization and tenant isolation | Implemented | Typed checks dan negative tests. |
| 13 | Input validation | Partial | Non-destructive validators dan API contract checks. |
| 14 | Server-side validation | Safe equivalent | Synthetic/canary proof di lab, bukan arbitrary RCE. |
| 15 | API logic validation | Implemented | OpenAPI contract dan authorization matrix checks. |
| 16 | Identity posture | Partial | Read-only posture checks. |
| 17 | Cloud IAM posture | Partial | Read-only policy/control validation. |
| 18 | Endpoint posture | Safe equivalent | Inventory dan control validation; tidak ada implant. |
| 19 | Detection validation | Implemented | Detection checks dan controlled callback markers. |
| 20 | Supply-chain posture | Partial | Dependency inventory dan risk normalization. |
| 21 | Controlled lab | Implemented | Disposable HTTP fixture, synthetic data, reset flow. |
| 22 | Safe proof and replay | Implemented | Simulation fixtures, replay verification, negative tests. |
| 23 | Evidence ledger | Implemented | Hash chain, redaction, manifest, parent-child references. |
| 24 | Attack path, risk and P0/P1 | Implemented | Deterministic scoring dari evidence, impact, reachability, dan confidence. |
| 25 | Reporting, remediation and retest | Implemented | JSON/Markdown delivery, remediation lifecycle, retest, cleanup. |

## Mapping arsitektur blueprint

| Istilah blueprint | Implementasi repository | Batasan |
|---|---|---|
| Frontend Angular | `apps/frontend-angular` | Operator dashboard untuk engagement, jobs, findings, readiness, dan reports. |
| API Gateway .NET 10 | `apps/gateway-dotnet` | Policy-gated API; bukan covert listener atau beacon endpoint. |
| LangGraph orchestrator | `apps/orchestrator-langgraph` dan `src/assessment/orchestrator` | Planner menghasilkan typed workflow yang harus melewati policy. |
| C2 framework | Controlled assessment task runner | Tidak ada implant, persistence, covert channel, atau arbitrary command. |
| Agent/implant | Lab agent simulator dan disposable fixture | Hanya synthetic target dan controlled callback marker. |
| Malleable profile | Assessment profile yang eksplisit dan diaudit | Tidak digunakan untuk menyamarkan malware traffic atau evasi deteksi. |
| Decoy/deception | Isolated deception lab fixture | Tidak menyembunyikan operasi pada target produksi. |
| Infrastructure | `deploy/terraform`, `deploy/ansible`, `deploy/nginx` | Deployment berfokus pada control plane, lab, logging, dan safe defaults. |
| Brain | `src/assessment/modules/brain` | Keputusan selalu policy-gated dan dapat dihentikan operator. |
| Console | Angular dashboard dan API client | Evidence dan audit menjadi sumber kebenaran. |

## Capability boundary

Capability berikut didukung jika scope dan authorization valid:

- passive discovery dan bounded network/web assessment;
- API contract, authentication/authorization, dan tenant-isolation validation;
- read-only identity, cloud, endpoint, detection, dan supply-chain posture;
- synthetic canary dan controlled callback marker;
- lab-only proof dengan disposable fixtures;
- evidence capture, redaction, hash-chain verification, risk scoring, reporting, remediation, retest, dan cleanup.

Capability berikut ditolak pada policy layer dan tidak boleh ditambahkan sebagai production feature:

- credential collection atau extraction;
- keylogging, browser secret capture, atau webcam/screen harvesting;
- hidden persistence, process injection, rootkit, atau firmware modification;
- covert C2, traffic evasion, proxy rotation untuk menghindari deteksi, atau DNS/ICMP tunneling;
- arbitrary command execution dan shellcode/payload delivery;
- data exfiltration;
- destructive write, ransomware, wiper, log deletion, atau availability attack.

## Definition of done

Sebuah layer hanya boleh diberi status **Implemented** apabila:

1. capability memiliki contract atau typed interface;
2. execution melewati authorization, scope, time window, budget, dan capability policy;
3. evidence memiliki timestamp, hash, redaction state, dan parent reference;
4. terdapat positive dan negative test;
5. output dapat direplay atau diverifikasi;
6. production-dangerous behavior diganti dengan synthetic fixture atau lab-only proof;
7. dokumentasi tidak mengklaim fitur yang belum tersedia.

## Posisi terhadap blueprint full-attack

Blueprint full-attack dapat digunakan sebagai referensi threat-model dan daftar kebutuhan red-team, tetapi **tidak boleh dianggap sebagai status implementasi repository ini**. Repository ini sengaja mengutamakan auditability, authorization, reproducibility, dan fail-closed behavior.
