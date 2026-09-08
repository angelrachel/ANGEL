# Dependency, Provenance, and License Register

## Purpose

This document records the current dependency and provenance baseline for ANGEL. It distinguishes ordinary use of an external dependency from copied or vendored source. It does not make an authorship or plagiarism determination.

## Current declared dependencies

| Manifest | Dependency | Current use observed | Follow-up |
|---|---|---|---|
| `requirements.txt` | `cryptography` | Imported by `src/c2/server/crypto/crypto.py` | Pin exact version and hash; record license. |
| `requirements.txt` | `websockets` | Declared for the Python websocket path | Confirm runtime caller and pin exact version/hash. |
| `requirements.txt` | `ruff`, `mypy`, `bandit`, `pytest`, `pytest-cov`, `pre-commit`, `detect-secrets` | Development/CI tooling | Pin versions and verify CI uses the same lock policy. |
| `go.mod` | `github.com/hirochachacha/go-smb2` | Declared; current direct use must be verified | Confirm whether it is still needed; preserve license notice if retained. |
| `go.mod` | `golang.org/x/sys` | Used by platform-specific code in parts of the tree | Pin and record upstream license. |
| `go.mod` | `github.com/geoffgarside/ber` | Declared indirect dependency | Confirm transitive path and necessity. |
| `go.mod` | `golang.org/x/crypto` | Declared indirect dependency | Confirm transitive path and necessity. |

## Provenance findings

The repository currently has no root `LICENSE`, `COPYING`, or `NOTICE` file and no consolidated third-party notice. This is a release and hiring-review blocker because dependency obligations cannot be demonstrated from the repository alone.

The audit did **not** establish evidence that the current source was copied from another project. Similarity to common security techniques, names of public tools, or use of standard libraries is not sufficient evidence of copying. The following items still require human provenance review before distribution:

- External executables and tools referenced by source, including Mimikatz, Certipy, Nginx, Terraform providers, Ansible collections, `flashrom`, QEMU, and EFI utilities.
- Historical or current technique names that may correspond to public research or projects.
- The fixed development key in `src/c2/server/crypto/shared_key.py`.
- The tracked `server` ELF artifact, whose source revision and build process are not currently recorded.

## Required evidence before claiming originality/compliance

1. Add an appropriate root license and project copyright/ownership statement.
2. Generate a third-party notice from the exact dependency lockfiles.
3. Pin Python dependencies with exact versions and hashes; retain `go.sum` and commit Terraform provider lockfiles.
4. Record upstream URLs, versions, checksums, and licenses for every external executable or collection required by deployment.
5. Remove or quarantine generated binaries from source control unless a release-artifact policy explicitly requires them.
6. Replace development/default secrets with documented runtime configuration and never commit operational keys.
7. Keep source provenance notes for any code imported, adapted, or generated from an external source.

## Automated check

Run the read-only integrity checker from the repository root:

```bash
python3 scripts/repo_integrity_check.py
```

The checker detects tracked generated files, tracked binaries, empty tracked files, broken local Python imports, unpinned Python dependencies, Terraform state files, and missing root provenance documents. It does not execute attack modules and does not modify the repository.
