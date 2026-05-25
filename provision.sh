#!/usr/bin/env bash

set -euo pipefail

AWS_DIR="./vagrant/.aws"
CREDENTIALS_FILE="$AWS_DIR/credentials"

CONNECTOR_TEMPLATE="./vagrant/kafka/connect/s3.json.template"
CONNECTOR_OUTPUT="./vagrant/kafka/connect/s3.json"

mkdir -p "$AWS_DIR"

# Check if credentials already exist
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

# Ask for S3 bucket name 
read -rp "S3 Bucket Name: " S3_BUCKET_NAME 

# Generate connector config from template 
sed "s/<S3_BUCKET_NAME>/${S3_BUCKET_NAME}/g" \
  "$CONNECTOR_TEMPLATE" > "$CONNECTOR_OUTPUT" 

echo "Generated connector config:" 
echo "$CONNECTOR_OUTPUT"

cd vagrant

vagrant up

ansible-playbook -i hosts playbook.yml --timeout 20
