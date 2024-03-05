#!/bin/bash

set -euo pipefail

# User's host directories mounted as volume.
export GW_WIREMOCK_PATH="/home/mock"
export GW_CONTRACTS_PATH="/contracts"
export GW_CERTS_PATH="/etc/ssl/mock/share"

SCRIPTS=$(realpath "$(dirname "${BASH_SOURCE[0]}")")

# grpc-wiremock setup scripts variables.
export SCRIPTS="${SCRIPTS}"
export MOCKS="${SCRIPTS}/mocks"
export PROXY="${SCRIPTS}/proxy"
export MULTIAPI="${SCRIPTS}/multiapi"
export OTHER="${SCRIPTS}/other"
export ROUTING="${SCRIPTS}/routing"
export CERTS="${SCRIPTS}/routing/certs"

export WIREMOCK_ADDR="localhost:9000"

# Headers for rsyslog.
export ENTRYPOINT_HEADER="gw.entrypoint"
export WIREMOCK_RUN_HEADER="gw.wiremock.run"
export PROXY_GEN_HEADER="gw.proxy.gen"
export PROXY_WATCH_HEADER="gw.proxy.watch"
export MOCKS_GEN_HEADER="gw.mocks.gen"
export MOCKS_WATCH_HEADER="gw.mocks.watch"
export ROUTING_CERTS_GEN_HEADER="gw.routing.certs.gen"
export ROUTING_NGINX_GEN_HEADER="gw.routing.nginx.gen"
export ROUTING_NGINX_WATCH_HEADER="gw.routing.nginx.watch"
export ROUTING_NGINX_LOGS_HEADER="gw.routing.nginx.logs"
export MULTIAPI_LOGS_HEADER="gw.multiapi.supervisord.logs"
