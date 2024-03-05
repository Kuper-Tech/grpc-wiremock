#!/bin/bash

set -euo pipefail

bash "${PROXY}/install.sh" 2>&1 | logger -t "${PROXY_GEN_HEADER}"