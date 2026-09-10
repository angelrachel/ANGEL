#!/usr/bin/env bash
set -euo pipefail

if [[ "${ANGEL_ALLOW_INFRA_CLEANUP:-}" != "1" ]]; then
  echo "Refusing destroy: set ANGEL_ALLOW_INFRA_CLEANUP=1 after reviewing the plan in an authorized lab." >&2
  exit 1
fi
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$ROOT/deploy/terraform"
terraform init -input=false
terraform destroy -input=false -auto-approve
printf '%s\n' '[+] Infrastructure cleanup completed'
