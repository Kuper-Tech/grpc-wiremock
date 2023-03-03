#!/bin/bash

set -euo pipefail

PROXY_HOST="http://localhost"
PROXY_INPUT_PORT="3010"
PROXY_OUTPUT_PORT="80"
PROXY_BASE_URL="${PROXY_HOST}:${PROXY_OUTPUT_PORT}"
PROXY_OUTPUT_PATH="/var/proxy/proxy"

# Kill previous proxy instance if exists.
PID=$(sudo lsof -t -i:${PROXY_INPUT_PORT}) || true
if [ ! -z "${PID}" ]; then
  kill -9 "${PID}"
fi

# Generate new proxy.
rm -rf ${PROXY_OUTPUT_PATH}
mkdir -p ${PROXY_OUTPUT_PATH}

grpc2http \
  --input "${GW_CONTRACTS_PATH}" \
  --output "${PROXY_OUTPUT_PATH}" \
  --base-url "${PROXY_BASE_URL}"

# Install new proxy.
make -C "${PROXY_OUTPUT_PATH}" install
