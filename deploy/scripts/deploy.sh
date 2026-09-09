#!/usr/bin/env bash
set -euo pipefail

if [[ "${ANGEL_ALLOW_APPLY:-}" != "1" ]]; then
  echo "Refusing deployment: set ANGEL_ALLOW_APPLY=1 after reviewing the plan in an authorized lab." >&2
  exit 1
fi
: "${ANGEL_OPERATOR_KEY:?ANGEL_OPERATOR_KEY is required}"
: "${ANGEL_SHARED_KEY:?ANGEL_SHARED_KEY is required}"
: "${ANGEL_AUTH_SECRET:?ANGEL_AUTH_SECRET is required}"

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$ROOT/deploy/terraform"
terraform init -input=false
terraform plan -input=false
terraform apply -input=false -auto-approve

echo "[+] Waiting for lab control plane..."
sleep 30
cd ../ansible
ansible-playbook -i inventories/production/hosts.ini playbooks/site.yml
printf '%s\n' '[+] Deployment complete'
