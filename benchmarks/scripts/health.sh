#!/bin/bash

set -euo pipefail

DEFAULT_PORT=8000
PORT=${DEFAULT_PORT}

for i in $(seq 1 "${WIREMOCK_API_COUNT}"); do
  echo "===> Check health, port: ${PORT}"
  curl --fail http://localhost:${PORT}/HealthCheck || exit 1

  PORT=$((PORT+1))
done
