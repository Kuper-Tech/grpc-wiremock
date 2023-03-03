#!/bin/bash

set -euo pipefail

bash /docker-entrypoint.sh 2>&1 | logger -t "${WIREMOCK_RUN_HEADER}" &

bash "${OTHER}/wait_for_it.sh" "${WIREMOCK_ADDR}" | logger -t "${WIREMOCK_RUN_HEADER}"