#!/bin/bash

set -euo pipefail

if ! certgen | logger -t "${ROUTING_CERTS_GEN_HEADER}"; then
  echo "Certificate initial generating failed. Skip"
fi

[ "$(ps | grep '[n]ginx')" == "" ] && sudo nginx

if ! confgen | logger -t "${ROUTING_NGINX_GEN_HEADER}"; then
  echo "Nginx configs generator failed. Skip"
fi

if ! certgen | logger -t "${ROUTING_CERTS_GEN_HEADER}"; then
  echo "Certificate generating for mock hosts failed. Skip"
fi
