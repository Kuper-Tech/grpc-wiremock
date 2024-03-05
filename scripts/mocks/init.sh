#!/bin/bash

set -euo pipefail

DOMAINS_PATH="${GW_CONTRACTS_PATH}/services"

if ! mockgen first-run \
  --domains-path="${DOMAINS_PATH}" \
  --wiremock-path="${GW_WIREMOCK_PATH}"; then

  echo "Mocks autogen exited with an error. Skip"
fi
