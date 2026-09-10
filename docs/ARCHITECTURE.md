# ANGEL Architecture

## Control flow

```text
Operator
  -> API gateway
  -> engagement + authorization validation
  -> scope matcher + capability gate
  -> signed assessment job
  -> typed worker / safe adapter
  -> evidence ledger
  -> finding + risk engine
  -> report, remediation, retest
```

Every transition is policy-gated. A target mismatch, expired engagement, invalid signature, replayed nonce, missing approval, emergency stop, exhausted budget, or permanently denied capability fails closed before target interaction.

## Domain boundaries

The control plane owns engagement, operator identity, RBAC, scope, job lifecycle, configuration references, audit events, and kill switch. Assessment adapters are typed and receive only normalized inputs that have already passed policy. Adapters produce observations; they do not assign final severity or bypass the evidence ledger.

Evidence is immutable from the application perspective. Each record carries engagement and job identity, target, check identity, UTC timestamp, engine version, redaction status, content hash, parent reference, and verification status. Reports are deterministic projections from the evidence and finding stores and include an integrity hash.

The risk engine combines exploitability, reachability, authentication requirement, privilege requirement, data sensitivity, integrity and availability impact, business criticality, blast radius, detection coverage, confidence, and remediation availability. P0/Critical and P1/High are policy outcomes, not target quotas.

## Safe capability boundary

ANGEL supports passive discovery, bounded network and web checks, API contract validation, authorization and tenant isolation checks, read-only identity/cloud/endpoint posture, detection validation, supply-chain inventory, synthetic canary proof, controlled callback markers, lab fixtures, evidence capture, scoring, reporting, remediation, retest, and cleanup.

The following capabilities are denied at the policy layer: credential collection or extraction, persistence, destructive write, log deletion, covert channel, process injection, evasion, data exfiltration, and arbitrary command execution. Production-dangerous proof is moved to a disposable lab with synthetic data and reset verification.

## Deployment boundary

Local deployment uses the ANGEL service, an isolated fixture, and a persistent data volume. Containers run as non-root, drop capabilities, enable `no-new-privileges`, and use read-only filesystems where possible. Staging and production require organization-owned infrastructure, external secret management, encrypted backup, restore verification, access review, TLS, retention, and an approved migration plan.

## Verification matrix

| Area | Required proof |
|---|---|
| Governance | Authorization, schedule, scope, approval, stop, audit |
| Policy | Domain/CIDR/URL/fixture matching, deny list, rate budget, time window |
| Execution | Typed task, signature, lifecycle, timeout, cancellation, replay defense |
| Evidence | Redaction, hash chain, manifest, integrity verification, replay |
| Risk | Reachability, privilege path, business impact, confidence, deterministic severity |
| Delivery | JSON/Markdown report, remediation lifecycle, retest and closure |
| Lab | Seed, execute, expected finding, evidence, reset, no residual state |
| Operations | Health/readiness/metrics, logs, backup/restore, cleanup manifest |
