# ANGEL Safe-Equivalent Roadmap

Roadmap ini menerjemahkan struktur blueprint ke pekerjaan yang dapat dikirimkan secara aman. Prioritasnya adalah membuat platform dapat dipakai untuk authorized engagement tanpa mengubahnya menjadi malware framework.

## Track A — Hardening control plane

| Pekerjaan | Output | Kriteria selesai |
|---|---|---|
| Authorization lifecycle | Contract approval, reviewer, expiry, emergency stop | Semua job expired/stopped ditolak sebelum adapter dipanggil. |
| Capability catalog | Versioned allow/deny catalog | Default decision deny dan perubahan tercatat di audit. |
| Tenant isolation | Organization-scoped queries | Cross-tenant access memiliki negative test. |
| Replay protection | Nonce, signature, lifecycle validation | Replay dan duplicate nonce selalu gagal closed. |
| Operational observability | Health, readiness, metrics, structured audit | Operator dapat memverifikasi status tanpa membaca secret. |

## Track B — Safe assessment modules

| Blueprint area | Safe implementation |
|---|---|
| C2/tasking | Typed assessment jobs dengan queue, timeout, cancellation, dan result store. |
| Implant/agent | Disposable lab agent simulator yang hanya mengirim synthetic status dan marker. |
| Malleable profile | Explicit assessment profile untuk fixture behavior; tidak untuk menyamarkan traffic. |
| Decoy | Isolated deception fixture untuk menguji detection dan routing; tidak untuk operasi tersembunyi. |
| SQL/NoSQL/input testing | Fixture-backed validators, parameterized probes, dan non-destructive evidence. |
| Exploit validation | Canary marker, controlled callback, dan expected-result assertion di lab. |
| Orchestrator/brain | Planner, risk assessment, hypothesis tracking, dan human approval gate. |
| Infrastructure | Terraform/Ansible untuk isolated lab, network boundary, logging, backup, dan reset. |
| Evidence/report | Immutable ledger, redaction, manifest, deterministic report, remediation, retest. |

## Track C — Verification

Setiap module baru harus memiliki:

- contract atau typed input/output;
- scope, capability, authorization, budget, dan time-window check;
- timeout dan cancellation;
- positive test dan negative test;
- evidence hash serta redaction;
- replay atau deterministic verification;
- lab fixture jika proof berisiko;
- cleanup/reset verification;
- dokumentasi threat model dan limitation.

## Explicitly not planned

Fitur-fitur berikut tidak menjadi pekerjaan yang tersisa dan tidak boleh ditambahkan ke production path:

- implant tersembunyi atau beacon yang menetap;
- persistence pada endpoint target;
- credential theft, keylogging, browser secret extraction, atau token harvesting;
- process injection, rootkit, firmware modification, dan anti-analysis;
- covert C2, evasion, traffic morphing, proxy rotation untuk bypass detection, DNS/ICMP tunneling;
- arbitrary shell/command execution, shellcode delivery, atau payload encryption untuk menghindari AV/EDR;
- data exfiltration;
- destructive write, ransomware, wiper, log deletion, atau availability disruption.

## Release gates

Sebelum merge atau release:

```bash
make validate
make app-check
```

Jika toolchain lokal tidak lengkap, pipeline CI harus menjalankan gate tersebut dengan Go, .NET SDK, Node.js, dan Python environment yang dipin. Tidak boleh menandai pekerjaan sebagai selesai hanya berdasarkan keberadaan file atau folder.
