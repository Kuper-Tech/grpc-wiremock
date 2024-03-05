#!/bin/bash

set -euo pipefail

CompileDaemon -color -log-prefix \
    -pattern "(.+\\.go|.+\\.c|.+\\.sh)$" \
    -build='make install -C cmd/grpc2http' \
    -command='bash ./dev/init_wiremock.sh'
