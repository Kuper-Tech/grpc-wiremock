#!/bin/bash

set -euo pipefail

log_header "Init Multi-API" "${ENTRYPOINT_HEADER}"
bash "${MULTIAPI}/init.sh"

log_header "Install grpc-to-http proxy" "${ENTRYPOINT_HEADER}"
bash "${PROXY}/init.sh"

# After mockgen task.
#log_header "Setup mocks" "${ENTRYPOINT_HEADER}"
#bash "${MOCKS}/init.sh"

log_header "Install certificates. Setup nginx" "${ENTRYPOINT_HEADER}"
bash "${ROUTING}/init.sh"

log_header "Initialization is done" "${ENTRYPOINT_HEADER}"
bash "${ROUTING}/logs.sh"
