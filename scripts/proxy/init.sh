#!/bin/bash

set -euo pipefail

bash "${PROXY}/watch.sh" 2>&1 | logger -t "${PROXY_WATCH_HEADER}" &