#!/usr/bin/env bash
# Runs the full pipeline experiment: for each format, fill the topic with a fixed
# backlog (constant lag), then drain it under the resource cap and record RPS.
#
# Every scenario is identical except the wire format:
#   - same message count N (same starting lag)
#   - same consumer CPU/memory cap
#   - same single-partition topic, fresh per format
#
# Usage:  ./run-pipeline.sh [format ...]
#   N=1000000 CONSUMER_CPUS=0.5 ./run-pipeline.sh json protobuf
set -euo pipefail

cd "$(dirname "$0")"
COMPOSE="docker compose -f pipeline/docker-compose.yml"

export N="${N:-2000000}"
export CONSUMER_CPUS="${CONSUMER_CPUS:-1.0}"
export CONSUMER_MEM="${CONSUMER_MEM:-512M}"

FORMATS=("$@")
if [ ${#FORMATS[@]} -eq 0 ]; then
  FORMATS=(json jsonshort msgpack cbor protobuf avro)
fi

echo "==> config: N=$N cap=${CONSUMER_CPUS}cpu/${CONSUMER_MEM}  formats=${FORMATS[*]}"
echo "==> starting kafka"
$COMPOSE up -d --wait kafka

results="pipeline-results.txt"
: > "$results"

for fmt in "${FORMATS[@]}"; do
  echo "============================================================"
  echo "==> [$fmt] fill: producing $N messages (sets constant lag)"
  FORMAT="$fmt" $COMPOSE run --rm producer
  echo "==> [$fmt] drain: consuming under ${CONSUMER_CPUS}cpu/${CONSUMER_MEM}"
  FORMAT="$fmt" $COMPOSE run --rm consumer 2>&1 | tee -a "$results"
done

echo "============================================================"
echo "==> RESULT lines:"
grep "RESULT" "$results" || true
echo "==> tearing down kafka"
$COMPOSE down -v
