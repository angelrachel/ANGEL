# Angular Frontend Boundary

This directory contains the Angular 20.3.0 standalone scaffold for the blueprint frontend boundary. The source is intentionally limited to a local observe/simulate UI foundation and does not contain live-target, command-execution, implant, persistence, or exploit endpoints.

## Local commands

```bash
npm install
npm run build
```

Map future UI work to `contracts/api/openapi.yaml`, `contracts/task/`, `contracts/authorization/`, `contracts/evidence/`, and `contracts/reporting/`. Keep live execution outside this boundary unless a separately reviewed architecture decision changes the policy gate.
