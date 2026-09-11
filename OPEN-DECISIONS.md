# ANGEL Open Decisions

Keputusan berikut wajib dikonfirmasi sebelum deployment production, tetapi tidak boleh menghambat implementasi lokal/lab:

1. Database production: SQLite untuk local/lab atau PostgreSQL untuk deployment multi-tenant.
2. Queue production: durable local queue atau managed broker.
3. Object storage evidence: local encrypted store atau managed object storage.
4. Identity provider dan MFA untuk operator.
5. Retention dan legal hold sesuai organisasi.
6. Connector cloud/SIEM/ticketing yang benar-benar diizinkan.
7. SLO, RTO, RPO, dan batas rate production.
8. Deployment target, certificate authority, dan key management.

Hermes harus menggunakan fixture lokal untuk keputusan yang belum diberikan dan menandai dependency sebagai `NEEDS-HUMAN-DECISION`.
