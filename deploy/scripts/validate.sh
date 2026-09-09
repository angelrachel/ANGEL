#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$ROOT"

command -v terraform >/dev/null || { echo "terraform is required" >&2; exit 1; }
terraform -chdir=deploy/terraform fmt -check -recursive
terraform -chdir=deploy/terraform init -backend=false -input=false
terraform -chdir=deploy/terraform validate
bash -n deploy/scripts/deploy.sh deploy/scripts/destroy.sh deploy/scripts/validate.sh
printf '%s\n' 'deployment static validation passed; no apply or destroy was executed'
