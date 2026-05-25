#!/usr/bin/env bash

set -euo pipefail

AWS_DIR="./vagrant/.aws"
CREDENTIALS_FILE="$AWS_DIR/credentials"

ENV_FILE="./.env"

CONNECTOR_TEMPLATE="./vagrant/kafka/connect/s3.json.template"
CONNECTOR_OUTPUT="./vagrant/kafka/connect/s3.json"

mkdir -p "$AWS_DIR"

#################################################
# AWS CREDENTIALS
#################################################

if [[ -f "$CREDENTIALS_FILE" ]]; then
    echo "Using existing AWS credentials from:"
    echo "$CREDENTIALS_FILE"
else
    echo "AWS credentials not found."
    echo "Enter your AWS credentials"

    read -rp "AWS Access Key ID: " AWS_ACCESS_KEY_ID
    read -rsp "AWS Secret Access Key: " AWS_SECRET_ACCESS_KEY
    echo

    cat > "$CREDENTIALS_FILE" <<EOF
[default]
aws_access_key_id=$AWS_ACCESS_KEY_ID
aws_secret_access_key=$AWS_SECRET_ACCESS_KEY
EOF

    chmod 600 "$CREDENTIALS_FILE"

    echo "AWS credentials file created at:"
    echo "$CREDENTIALS_FILE"
fi

echo

#################################################
# GITLAB PERSONAL ACCESS TOKEN
#################################################

if [[ -f "$ENV_FILE" ]] && grep -q "^GITLAB_PAT=" "$ENV_FILE"; then
  echo "Using existing GITLAB_PAT from $ENV_FILE"
else
  echo "GitLab PAT not found in environment file."

  read -rsp "Enter GitLab Personal Access Token: " GITLAB_PAT
  echo

  echo "GITLAB_PAT=\"$GITLAB_PAT\"" >> "$ENV_FILE"
  
  chmod 600 "$ENV_FILE"

  echo ".env file created at:"
  echo "$ENV_FILE"
fi

echo

#################################################
# KAFKA CONNECT S3 BUCKET
#################################################

if test -f "$CONNECTOR_OUTPUT"; then
  echo "Using existing S3 sink connector config:"
  echo "$CONNECTOR_OUTPUT"
else
  read -rp "S3 Bucket Name: " S3_BUCKET_NAME 

  # Generate connector config from template 
  sed "s/<S3_BUCKET_NAME>/${S3_BUCKET_NAME}/g" \
    "$CONNECTOR_TEMPLATE" > "$CONNECTOR_OUTPUT" 

  echo "Generated connector config:" 
  echo "$CONNECTOR_OUTPUT"
fi

#################################################
# START ENVIRONMENT
#################################################

cd vagrant

vagrant up

ansible-playbook -i hosts playbook.yml --timeout 20
