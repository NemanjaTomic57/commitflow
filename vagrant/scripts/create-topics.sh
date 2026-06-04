#!/bin/bash

set -o errexit
set -o pipefail
set -o nounset

BOOTSTRAP_SERVER="192.168.56.10:9092"

delete_topic() {
  local topic="$1"

  echo "Deleting topic: $topic"

  kafka-topics.sh \
    --bootstrap-server "$BOOTSTRAP_SERVER" \
    --topic "$topic" \
    --delete
}

create_topic() {
  local topic="$1"

  echo "Creating topic: $topic"

  kafka-topics.sh \
    --bootstrap-server "$BOOTSTRAP_SERVER" \
    --create \
    --topic "$topic"
}

TOPICS=(
  "git.commits"
)

# for topic in "${TOPICS[@]}"; do
#   delete_topic "$topic"
# done

for topic in "${TOPICS[@]}"; do
  create_topic "$topic"
done
