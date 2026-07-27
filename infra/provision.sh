#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

terraform -chdir="${SCRIPT_DIR}/00_base" apply -auto-approve

cd "${SCRIPT_DIR}/00_base/ansible"

python build_inventory.py
ansible-playbook -i inventory.yml playbook.yml

terraform -chdir="${SCRIPT_DIR}/01_ecs" apply -auto-approve
