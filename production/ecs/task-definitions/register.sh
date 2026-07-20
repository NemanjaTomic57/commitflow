#!/usr/bin/env bash

set -euo pipefail

DIR="./"

if ! command -v aws >/dev/null 2>&1; then
    echo "Error: AWS CLI is not installed."
    exit 1
fi

echo "Registering task definition files..."

for file in "$DIR"/*.json; do
    [[ -f "$file" ]] || continue

    filename=$(basename "$file")

    echo "Registering ${filename}..."
    aws ecs register-task-definition --cli-input-json file://"${filename}" --output off
done

echo "Done."
