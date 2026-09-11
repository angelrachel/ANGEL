# ANGEL Acceptance Criteria

ANGEL hanya dinyatakan siap pada area yang artifact, contract, test, evidence, dan lab assertion-nya tersedia.

- Semua service melewati health/readiness.
- Semua job melewati tenant, authorization, scope, capability, budget, timeout, cancellation, dan audit gate.
- Semua result memiliki provenance, redaction state, hash, dan trace ID.
- Semua finding memiliki evidence, confidence, impact, reachability, dan reviewer state.
- Report dan retest dapat dibuat ulang dari evidence.
- Lab setup, assertion, reset, dan deterministic rerun lulus.
- Tidak ada secret, arbitrary command, destructive action, hidden persistence, credential extraction, evasion, atau data exfiltration.
- Status `VERIFIED` hanya diberikan setelah command dan test benar-benar lulus.
