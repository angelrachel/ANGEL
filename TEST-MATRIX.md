# ANGEL Test Matrix

| Gate | Command | Evidence | Blocking rule |
|---|---|---|---|
| Format | `make fmt` | clean diff | fail blocks |
| Unit | `make test` | test output | fail blocks |
| Static | `make static-check` | checker output | fail blocks |
| Safe contracts | `make safe-contracts` | validator output | fail blocks |
| Integration | `make app-check` | gateway/planner/lab output | fail blocks |
| Lab | `make lab-check && make lab-smoke && make lab-control` | fixture assertions | fail blocks |
| Traceability | `bash scripts/verify-traceability.sh` | 807 IDs and paths | fail blocks |
| Secret scan | repository secret scanner | clean result | fail blocks |
| Reset/rerun | `bash lab/scripts/reset.sh` or equivalent | same expected hashes | fail blocks |
