#!/bin/bash

set -euo pipefail

# Watch contracts and restart grpc-to-http proxy.
CompileDaemon \
  -run-dir="/var/proxy" \
  -directory="${GW_CONTRACTS_PATH}" \
  -build="${PROXY}/run.sh" \
  -command="grpc-to-http-proxy" \
  -pattern="(.+\\.)proto$" \
  -log-prefix=false \
  -graceful-kill -graceful-timeout=3
