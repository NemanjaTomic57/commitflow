#!/usr/bin/env bash

set -euo pipefail

AWS_REGION="eu-central-1"
ECR_REPOSITORY="commitflow"

# Pre-build
TIMESTAMP="$(date -u +'%Y%m%d-%H%M%S')"

ACCOUNT_ID="$(
  aws sts get-caller-identity \
    --query 'Account' \
    --output text
)"

ECR_REGISTRY="${ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com"
ECR_URI="${ECR_REGISTRY}/${ECR_REPOSITORY}"

aws ecr get-login-password --region "${AWS_REGION}" \
  | docker login \
      --username AWS \
      --password-stdin "${ECR_REGISTRY}"

cat /etc/os-release

# Test & Lint

golangci-lint run   
yamllint .
hadolint Dockerfile 
test -z "$(gofmt -l .)"
go vet ./...

# Build

docker build -t "${ECR_URI}:${TIMESTAMP}" .
docker tag "${ECR_URI}:${TIMESTAMP}" "${ECR_URI}:latest"

docker push "${ECR_URI}:${TIMESTAMP}"
docker push "${ECR_URI}:latest"
