#!/bin/bash

set -euo pipefail

log_header() {
  local message="${1}" header="${2}"

  echo "▛ ▞ ▟ ${message}" \
    | tr '[a-z]' '[A-Z]' \
    | logger -t "${header}"
}