#!/bin/bash

set -euo pipefail

sudo touch /var/log/nginx/{access,error,not_mocked}.log
tail -f /var/log/nginx/* | logger -t "${ROUTING_NGINX_LOGS_HEADER}"
