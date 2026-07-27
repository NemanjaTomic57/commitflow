#!/bin/bash -xeu

DIR=$(pwd)

cd "${DIR}"/terraform/base

terraform apply -destroy -auto-approve
terraform apply -auto-approve

cd "$DIR"/ansible

python build_inventory.py

ansible-playbook -i inventory.yml playbook.yml
