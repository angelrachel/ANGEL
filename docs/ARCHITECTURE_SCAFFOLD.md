# Architecture Scaffold Guide

The repository now follows the executable P0/P1 blueprint in `docs/IMPLEMENTATION_BLUEPRINT.md`. The lab remains restricted to registered assets and typed assessment operations.

| Boundary | Location | Current state | Next coding task |
|---|---|---|---|
| API and schemas | `contracts/` | Versioned task, result, lifecycle, and OpenAPI contracts | Add persistent API integration and compatibility tests |
| Go control-plane | `src/c2/` | Tested foundation with P0 lifecycle state machine | Connect lifecycle to durable task persistence and evidence |
| Angular UI | `apps/frontend-angular/` | Simulation workspace and typed boundary scaffold | Connect operator views to gateway APIs |
| .NET gateway | `apps/gateway-dotnet/` | Fail-closed task validation boundary and tests | Add persistent gateway-to-control-plane integration |
| LangGraph/MCP | `apps/orchestrator-langgraph/` | Fixture-only planner with tests | Add bounded multi-step lab workflow nodes |
| Replay and contracts | `simulation/`, `contracts/` | Strict fixture replay and contract validation | Add P0 acceptance fixtures and service checks |
| Agent metadata | `agents/` | Registry fixture | Add metadata validation and versioning |
| Lab services | `lab/services/fixture-http/` | Runnable owned-lab HTTP service with container build and tests | Add database, identity, telemetry, and recovery services |
| Infrastructure | `deploy/`, `docker-compose.yml` | Guarded deployment baseline plus fixture service composition | Run static validation and isolated lab deployment |

## Local commands

```bash
make fmt
make test
make vet
make build
make safe-contracts
make app-check
make validate
```

`make validate` does not deploy infrastructure. It checks source formatting, Go quality, repository hygiene, inventory generation, static findings, safe contracts, and shell syntax. `make app-check` builds the Angular scaffold and runs the .NET and LangGraph tests.

`deploy/scripts/deploy.sh` requires `ANGEL_ALLOW_APPLY=1` plus all three runtime secrets. `deploy/scripts/destroy.sh` requires `ANGEL_ALLOW_DESTROY=1`. Both guards are intentional and prevent an accidental apply/destroy from a copied shell command.

## Toolchain status

Go 1.27.1, Node 22.13.0/npm, Python 3.12, .NET SDK 10.0.401, Angular CLI 20.3.0, and LangGraph 1.2.11 are installed in the current development environment. Terraform remains unavailable because its binary download timed out. Node 22.13.0 is compatible with the pinned Angular CLI version. Safe Angular, .NET, and LangGraph scaffolds now exist; their routes remain limited to the approved observe/simulate contract boundary.

For a new shell in this development environment:

```bash
export DOTNET_ROOT="$HOME/.dotnet"
export PATH="$DOTNET_ROOT:$PATH"
export PATH="$HOME/.local/bin:$PATH"
source .venv-langgraph/bin/activate
```

The LangGraph pin is recorded in `apps/orchestrator-langgraph/requirements.txt`. The repository-local `.venv-langgraph` is ignored by Git and may be used for planner development.

All fixtures use `fixture://` references and the simulation contracts permit only `observe` or `simulate` modes.
