#!/usr/bin/env bash

set -euo pipefail

ENV_FILE="./.env"

env_err() {
  echo "Resolution: Copy .env.example to .env and enter all required variables to continue"
  exit 1
}

#################################################
# LOAD ENV FILE
#################################################

if [[ -f "$ENV_FILE" ]]; then
  set -a
  source "$ENV_FILE"
  set +a
else
  echo "ERROR: .env not found"
  env_err
fi

#################################################
# CHECK ENVIRONMENT VARIABLES
#################################################

if [[ "${GITHUB_PAT:-}" == "" ]]; then
  echo "ERROR: GitHub Personal Access Token not set in .env"
  env_err
fi

if [[ "${GITLAB_PAT:-}" == "" ]]; then
  echo "ERROR: GitLab Personal Access Token not set in .env"
  env_err
fi

#################################################
# START ENVIRONMENT
#################################################

cd vagrant

vagrant up

ansible-playbook -i hosts playbook.yml --timeout 20
