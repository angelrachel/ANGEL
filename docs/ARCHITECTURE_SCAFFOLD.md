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

Go 1.27.1, Node 22.13.0/npm, and Python 3.12 are available in the current development environment. Terraform, the .NET SDK, Angular CLI, and the LangGraph virtual environment are not verified as available in every shell. Node 22.13.0 is compatible with the pinned Angular CLI version. The Angular and .NET application sources remain ungenerated boundaries until their contracts are approved.

For a new shell in this development environment:

```bash
export DOTNET_ROOT="$HOME/.dotnet"
export PATH="$DOTNET_ROOT:$PATH"
export PATH="$HOME/.local/bin:$PATH"
source "$HOME/.venvs/angel-langgraph/bin/activate"
```

The LangGraph pin is recorded in `apps/orchestrator-langgraph/requirements.txt`; install it into a local virtual environment before using that advisory boundary.

All fixtures use `fixture://` references and the simulation contracts permit only `observe` or `simulate` modes.
