# Security Policy

## Scope

ANGEL is intended for authorized, non-destructive security assessment and defensive validation. Do not use it against systems without explicit written authorization. Do not commit credentials, target data, private keys, or engagement evidence to this repository.

## Reporting a vulnerability

Please open a **private GitHub security advisory** for this repository, or contact the repository owner through the organization-approved security channel. Do not publish exploit details, secrets, or proof-of-concept payloads in a public issue.

When reporting, include the affected commit or version, impact, reproduction steps using synthetic data, and a suggested mitigation. Redact tokens, personal data, hostnames, and customer information.

## Supported versions

| Version | Supported |
|---|---|
| `main` | Yes |
| development snapshots | Best effort |

## Secret handling

Use environment variables or a managed secret store for `ANGEL_OPERATOR_KEY` and `ANGEL_SHARED_KEY`. Rotate any secret that has been exposed, and treat the local SQLite database as sensitive because it may contain evidence metadata.

## Disclosure

The repository owner should coordinate remediation and disclosure timelines with affected stakeholders. This policy does not grant permission to test third-party systems.
