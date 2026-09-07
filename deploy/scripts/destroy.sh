#!/bin/bash
set -e

echo "[+] Destroying ANGEL infrastructure..."

cd deploy/terraform
terraform destroy -auto-approve

echo "[+] Infrastructure destroyed!"
