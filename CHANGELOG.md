# Changelog

## [Unreleased]
### Added
- Repository initialized
- Project structure
- Hardened evidence chain verification for sequence, parent links, timestamps, and recomputed hashes.
- Added unit coverage for evidence tamper detection and parent-child relation validation.

### Changed
- Removed the source-code default operator API key; `ANGEL_OPERATOR_KEY` is injected by the runtime and authentication fails closed when absent.
- Added bounded JSON request decoding, strict JSON fields, and method checks to control-plane endpoints.
- Fixed IPv6 endpoint formatting in the TCP/UDP inventory scanner.
- Added a fail-closed engagement scope and approval gate with unit coverage.
- Applied the scope gate to the agent router; accepted routes are simulation acknowledgements and do not execute target actions.
- Hardened API middleware and route registry with fail-closed auth, constant-time token comparison, concurrency safety, and JSON health headers.
- Added report validation for identity, chronology, required finding fields, and supported severity values, plus Markdown escaping tests.
- Hardened the control-plane store with RWMutex snapshots and generation-aware expiry timers.
- Hardened task/result stores with defensive copies, nil rejection, and read locks.
- Fixed AES-GCM to propagate nonce RNG failures and reject truncated ciphertext.
- Removed tracked duplicate ELF build artifacts (`c2` and `server`).
- Removed stale Python manifests and CI jobs because the repository contains no Python source.
- Replaced the stale Python Dockerfile with a pinned Go multi-stage build and aligned Compose to port 8001.
- Made Compose operator/shared/auth secrets mandatory instead of silently using development defaults.
