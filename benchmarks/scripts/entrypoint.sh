#!/bin/bash

set -euo pipefail

WIREMOCK_API_COUNT="${WIREMOCK_API_COUNT-1}"

bash /scripts/run.sh "${WIREMOCK_API_COUNT}"