#!/bin/bash -xeu

DIR=$(pwd)

cd "${DIR}"/00_base

terraform apply -destroy -auto-approve
terraform apply -auto-approve

cd ansible

python build_inventory.py

ansible-playbook -i inventory.yml playbook.yml
