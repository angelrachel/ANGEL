# ANGEL Database Schema

## Pendekatan migration

Repository ini menggunakan mekanisme migrasi yang berada di `src/assessment/server/database/`. Pilih:
- SQLite untuk local/lab (sesuai OPEN-DECISIONS #1).
- PostgreSQL untuk production multi-tenant (menunggu keputusan manusia).

## Entity utama

### Tenant / Organization
- id, name, created_at, updated_at, stopped

### Operator / User
- id, organization_id, identity_ref, role, status, created_at

### Engagement
- id, organization_id, name, authorized, starts_at, ends_at, policy_hash, stopped, created_at

### Scope
- id, engagement_id, kind (fixture/hostname/url/cidr), value, actions, excluded

### Authorization record
- id, engagement_id, operator_id, approval_id, approval_state, expires_at, created_at

### Job / Task
- id, engagement_id, agent_id, task_type, target_ref, action, mode, requested_by, state, version, updated_at, signature, nonce, created_at

### Evidence
- id, engagement_id, job_id, target, check_id, observed_at, engine_version, redacted, content_hash, parent_id, verification_status, created_at

### Finding
- id, engagement_id, job_id, evidence_ids, severity, title, description, confidence, reachability, business_criticality, blast_radius, status, reviewer_id, created_at

### Remediation
- id, finding_id, owner, plan, status, due_at, created_at

### Audit event
- id, engagement_id, event_type, actor_id, subject_id, payload_hash, chain_hash, created_at

## Sumber

- Lifecycle schema: `contracts/task/lifecycle.schema.json`
- OpenAPI: `contracts/api/openapi.yaml`
- Alignment: `docs/ANGEL_IMPLEMENTATION_ALIGNMENT.md`

Status: IMPLEMENTED secara concept; migrasi aktual berada di owner path database.
