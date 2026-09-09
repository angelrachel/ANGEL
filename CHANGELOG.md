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
