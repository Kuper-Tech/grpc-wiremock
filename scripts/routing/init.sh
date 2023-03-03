#!/bin/bash

set -euo pipefail

if ! certgen | logger -t "${ROUTING_CERTS_GEN_HEADER}"; then
  echo "Certificates generator exited with an error. Skip"
fi

sudo nginx

if ! confgen | logger -t "${ROUTING_NGINX_GEN_HEADER}"; then
  echo "Nginx configs generator exited with an error. Skip"
fi

