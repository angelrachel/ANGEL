# ANGEL — Authorized Security Assessment Control Plane

ANGEL is a control-plane foundation for **authorized, non-destructive security assessments**. It provides scope enforcement, operator authentication, rate limiting, evidence integrity, reporting, passive reconnaissance, and deterministic simulation for training and validation.

## Current status

The repository contains a tested defensive foundation. Offensive execution is intentionally not implemented.

### Implemented foundation

- Fail-closed scope policy with host/path allowlists, expiry, private-address blocking, and rate limiting.
- Operator authentication, RBAC, replay protection, encrypted sessions, and event streaming.
- Evidence normalization and hash-chain support for chain of custody.
- Passive inventory and reporting primitives.
- Rules of Engagement validation with explicit authorization windows and emergency contacts.
- Explicit safety capability allowlist and request rejection guardrails.
- CI quality gates for linting, typing, security scanning, tests, dependency auditing, and container builds.

## Safety boundary

ANGEL must not be used to deploy implants, exploit third-party systems, steal credentials, evade endpoint controls, establish persistence, move laterally, exfiltrate data, destroy data, or clear forensic records. The repository deliberately does not implement those capabilities.

Use only with written authorization and a defined engagement scope. Prefer a lab or staging environment, synthetic data, and read-only checks. Every active check must be separately approved and auditable.

## Blueprint coverage

| Blueprint area | Repository implementation | Status |
|---|---|---|
| Architecture and control plane | `src/c2`, `src/auth.py`, `src/rbac.py`, `src/orchestrator` | Implemented foundation |
| Scope and rules of engagement | `src/scope.py`, `src/engagement.py` | Implemented, fail-closed |
| Auditable assessment planning | `src/assessment.py` | Passive/simulation checks only |
| Passive reconnaissance and inventory | `src/osint`, `src/api_intel.py`, `src/graphql_intel.py` | Implemented, non-destructive |
| Passive observation findings | `src/passive_findings.py` | Implemented, conservative informational severity |
| Evidence and chain of custody | `src/evidence` | Implemented with redaction and hash chain |
| Evidence export verification | `src/evidence/manifest.py` | Implemented with deterministic manifest hash |
| Reporting | `src/reporting.py`, `src/report/export.py` | Implemented in Markdown and JSON |
| Approval, replay, cancellation, and bounded workflows | `src/orchestrator`, `src/cancellation.py` | Implemented |
| Deployment separation and hardening | `deploy/`, `Dockerfile`, `docker-compose.yml` | Lab-ready baseline |
| Implant execution, credential theft, persistence, evasion, lateral movement, rootkits, destructive impact | Intentionally absent | Owner-only / not implemented |

The final row is deliberate. These capabilities are not required to validate the defensive control-plane foundation and would create unacceptable risk without a separately reviewed lab, authorization package, and safety design.

Operational procedures are documented in [`docs/OPERATIONS.md`](docs/OPERATIONS.md).
The detailed section-by-section status is tracked in [`docs/BLUEPRINT_STATUS.md`](docs/BLUEPRINT_STATUS.md).

## Development

```bash
python -m pip install -r requirements.txt
pytest -q
ruff check src tests
mypy src
bandit -q -r src
```

## Remaining work for the owner

1. Replace placeholder contact and deployment values with organization-approved values.
2. Define the written rules of engagement, asset inventory, retention period, and incident escalation contacts.
3. Configure production secrets through a secret manager; never commit `.env` files or keys.
4. Review Terraform/Ansible plans with the infrastructure owner before applying them.
5. Decide which passive data sources are legally and contractually permitted.
6. Perform a human security review and authorize any future active testing modules.

## License and authorization

No production deployment or external assessment is authorized by this repository alone. Obtain explicit, written permission from the system owner before running any active test.
