#!/usr/bin/env bash

set -euo pipefail

set -a
source "../.env"
set +a

GITHUB_PAT="${GITHUB_PAT:-}"

if [[ -z "$GITHUB_PAT" ]]; then
    echo "Error: GITHUB_PAT environment variable is not set."
    exit 1
fi

API_URL="https://api.github.com"

# Fetch all repositories visible to the token
repos=$(curl -s \
    -H "Authorization: Bearer ${GITHUB_PAT}" \
    -H "Accept: application/vnd.github+json" \
    "${API_URL}/user/repos?per_page=100" |
    jq -r '.[].full_name')

total_commits=0

echo "Repository Commit Counts"
echo "========================"

for repo in $repos; do
    page=1
    commit_count=0

    while true; do
        commits=$(curl -s \
            -H "Authorization: Bearer ${GITHUB_PAT}" \
            -H "Accept: application/vnd.github+json" \
            "${API_URL}/repos/${repo}/commits?per_page=100&page=${page}")

        count=$(echo "$commits" | jq 'length')

        if [[ "$count" -eq 0 ]]; then
            break
        fi

        commit_count=$((commit_count + count))
        page=$((page + 1))
    done

    echo "${repo}: ${commit_count}"
    total_commits=$((total_commits + commit_count))
done

echo
echo "Total commits across all repositories: ${total_commits}"
