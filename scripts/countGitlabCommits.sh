#!/usr/bin/env bash

set -euo pipefail

set -a
source "../.env"
set +a

GITLAB_PAT="${GITLAB_PAT:-}"

if [[ -z "$GITLAB_PAT" ]]; then
    echo "Error: GITLAB_PAT environment variable is not set."
    exit 1
fi

API_URL="https://gitlab.com/api/v4"

# Fetch all projects visible to the token
projects=$(curl -s \
    --header "PRIVATE-TOKEN: ${GITLAB_PAT}" \
    "${API_URL}/projects?membership=true&per_page=100" |
    jq -r '.[].id')

total_commits=0

echo "Project Commit Counts"
echo "====================="

for project_id in $projects; do
    project_name=$(curl -s \
        --header "PRIVATE-TOKEN: ${GITLAB_PAT}" \
        "${API_URL}/projects/${project_id}" |
        jq -r '.path_with_namespace')

    page=1
    commit_count=0

    while true; do
        commits=$(curl -s \
            --header "PRIVATE-TOKEN: ${GITLAB_PAT}" \
            "${API_URL}/projects/${project_id}/repository/commits?per_page=100&page=${page}")

        count=$(echo "$commits" | jq 'length')

        if [[ "$count" -eq 0 ]]; then
            break
        fi

        commit_count=$((commit_count + count))
        page=$((page + 1))
    done

    echo "${project_name}: ${commit_count}"
    total_commits=$((total_commits + commit_count))
done

echo
echo "Total commits across all projects: ${total_commits}"
