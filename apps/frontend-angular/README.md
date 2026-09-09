# Angular Frontend Boundary

This directory reserves the frontend boundary from the blueprint. It is intentionally a contract-first scaffold; no Angular runtime is installed in the repository yet.

## Next local step

When the API contract is approved, initialize the application with the pinned Angular CLI in a separate change. Map the UI to `contracts/api/openapi.yaml`, `contracts/task/`, `contracts/authorization/`, `contracts/evidence/`, and `contracts/reporting/`.

The frontend must expose observe/simulate workflows only until an additional architecture decision is approved.
