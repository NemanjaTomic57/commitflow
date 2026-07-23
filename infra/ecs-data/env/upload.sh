#!/usr/bin/env bash

set -euo pipefail

BUCKET="s3://commitflow-761018874759/env"
ENV_DIR="./"

if ! command -v aws >/dev/null 2>&1; then
    echo "Error: AWS CLI is not installed."
    exit 1
fi

if [[ ! -d "$ENV_DIR" ]]; then
    echo "Error: Directory '$ENV_DIR' does not exist."
    exit 1
fi

echo "Uploading environment files to ${BUCKET}..."

for file in "$ENV_DIR"/*.env; do
    [[ -f "$file" ]] || continue

    filename=$(basename "$file")

    echo "Uploading ${filename}..."
    aws s3 cp "$file" "${BUCKET}/${filename}"
done

echo "Done."
