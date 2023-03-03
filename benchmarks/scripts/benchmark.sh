#!/bin/bash

set -euo pipefail

DEFAULT_PORT=8000
PORT=${DEFAULT_PORT}
OUT_DIR="/benchmarks/output/api_count_${WIREMOCK_API_COUNT}_$(date --iso=seconds)"

mkdir "${OUT_DIR}"

for i in $(seq 1 "${WIREMOCK_API_COUNT}"); do
  echo "===> Benchmark Wiremock API, port: ${PORT}"

  ab -c 10 -n 100 "http://wiremock:${PORT}/HealthCheck" > \
    "${OUT_DIR}/report_${PORT}.txt" 2>&1 &

  PORT=$((PORT+1))
done

sleep infinity