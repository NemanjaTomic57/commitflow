#!/usr/bin/env bash

set -euo pipefail

ENV_FILE="./.env"

AWS_DIR="./vagrant/.aws"
AWS_CREDENTIALS_FILE="$AWS_DIR/credentials"

CONNECTOR_TEMPLATE="./vagrant/kafka/connect/s3.json.template"
CONNECTOR_OUTPUT="./vagrant/kafka/connect/s3.json"

mkdir -p "$AWS_DIR"

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
# AWS CREDENTIALS
#################################################

if [[ "${AWS_ACCESS_KEY_ID:-}" == "" || "${AWS_SECRET_ACCESS_KEY:-}" == "" ]]; then
  echo "ERROR: AWS credentials not set in .env"
  env_err
fi

#################################################
# GITHUB PERSONAL ACCESS TOKEN
#################################################

if [[ "${GITHUB_PAT:-}" == "" ]]; then
  echo "ERROR: GitHub Personal Access Token not set in .env"
  env_err
fi

#################################################
# GITLAB PERSONAL ACCESS TOKEN
#################################################

if [[ "${GITLAB_PAT:-}" == "" ]]; then
  echo "ERROR: GitLab Personal Access Token not set in .env"
  env_err
fi

#################################################
# KAFKA CONNECT S3 SINK
#################################################

if [[ "${AWS_S3_BUCKET:-}" == "" ]]; then
  echo "ERROR: S3 bucket name not set in .env"
  env_err
fi

#################################################
# START ENVIRONMENT
#################################################

cd vagrant

vagrant up

ansible-playbook -i hosts playbook.yml --timeout 20
