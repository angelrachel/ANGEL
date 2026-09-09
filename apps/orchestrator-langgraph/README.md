# LangGraph/MCP Advisory Boundary

This directory contains a fixture-only planner scaffold for the blueprint orchestrator boundary. The planner validates `fixture://` targets and `observe` or `simulate` modes, then emits a structured plan that remains unauthorized until the Go policy gate approves it. It does not execute commands, alter scope, create implants, delete evidence, or select live targets.

## Local commands

```bash
python3 -m venv .venv-langgraph
.venv-langgraph/bin/pip install -r apps/orchestrator-langgraph/requirements.txt
.venv-langgraph/bin/python apps/orchestrator-langgraph/planner.py simulation/fixtures/tasks/recon-http.json
```

Inputs and outputs must follow the JSON contracts under `contracts/`. Human review and the Go policy gate remain mandatory before a simulation route is accepted.
