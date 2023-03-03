#!/bin/bash

set -euo pipefail

WIREMOCK_API_COUNT="${1}"
DEFAULT_PORT=8000
PORT=${DEFAULT_PORT}

for i in $(seq 1 "${WIREMOCK_API_COUNT}"); do
  echo "===> Run Wiremock API, port: ${PORT}"

  java \
    -cp "/var/wiremock/lib/*:/var/wiremock/extensions/*" \
    com.github.tomakehurst.wiremock.standalone.WireMockServerRunner --port ${PORT} &

  PORT=$((PORT+1))
done

sleep infinity