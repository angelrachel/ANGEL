# Angular Frontend Boundary

This directory reserves the frontend boundary from the blueprint. Angular CLI 20.3.0 is installed in the development environment and is compatible with Node 22.13.0, but the application source has not been generated yet.

## Next local step

When the API contract is approved, initialize the application with `ng new` using the pinned CLI. Map the UI to `contracts/api/openapi.yaml`, `contracts/task/`, `contracts/authorization/`, `contracts/evidence/`, and `contracts/reporting/`.

The frontend must expose observe/simulate workflows only until an additional architecture decision is approved.
