#!/bin/bash
set -e

echo "[+] Deploying ANGEL infrastructure..."

cd deploy/terraform
terraform init
terraform plan
terraform apply -auto-approve

echo "[+] Waiting for VPS to be ready..."
sleep 30

echo "[+] Running Ansible playbooks..."
cd ../ansible
ansible-playbook -i inventories/production/hosts.ini playbooks/site.yml

echo "[+] Deployment complete!"
echo "[+] Redirector IP: $(terraform output -raw redirector_ip)"
echo "[+] Teamserver IP: $(terraform output -raw teamserver_ip)"
