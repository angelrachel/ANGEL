# Architecture Scaffold Guide

The repository now contains a contract-first scaffold matching the blueprint boundaries without claiming that absent runtimes are installed or operational.

| Boundary | Location | Current state | Next coding task |
|---|---|---|---|
| API and schemas | `contracts/` | JSON/OpenAPI contracts | Add schema validation tests and versioning |
| Go control-plane | `src/c2/` | Existing tested foundation | Map Go structs to contracts |
| Angular UI | `apps/frontend-angular/` | README boundary only | Install/pin Angular CLI and generate app |
| .NET gateway | `apps/gateway-dotnet/` | README boundary only | Install/pin .NET 10 SDK and generate project |
| LangGraph/MCP | `apps/orchestrator-langgraph/` | Advisory boundary only | Add planner project and contract tests |
| Simulation | `simulation/` | Fixtures and reserved adapters | Add replay runner using fixture refs |
| Agent metadata | `agents/` | Registry fixture | Add metadata validation and versioning |
| Lab fixtures | `lab/` | Reserved directories | Add synthetic mock services only |
| Infrastructure | `deploy/` | Terraform/Ansible baseline | Run static validation; apply only in isolated lab |

## Local commands

```bash
make fmt
make test
make vet
make build
make validate
```

`make validate` does not deploy infrastructure. It checks source formatting, Go quality, repository hygiene, inventory generation, static findings, and shell syntax.

`deploy/scripts/deploy.sh` requires `ANGEL_ALLOW_APPLY=1` plus all three runtime secrets. `deploy/scripts/destroy.sh` requires `ANGEL_ALLOW_DESTROY=1`. Both guards are intentional and prevent an accidental apply/destroy from a copied shell command.

## Toolchain status

Go, Node/npm, Python, and Terraform are available in the current development environment. The .NET SDK is not installed, and Angular/LangGraph runtimes have not been initialized. This is intentional: the repository now has named boundaries and contracts without creating misleading, unbuildable application files.

All fixtures use `fixture://` references and the simulation contracts permit only `observe` or `simulate` modes.
