#!/bin/bash -xeu

DIR=$(pwd)

cd "${DIR}"/terraform/base

terraform apply -destroy -auto-approve
