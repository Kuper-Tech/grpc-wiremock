#!/bin/bash

set -euo pipefail

if [ $# -ne 0 ]; then
    echo "executing: '$@'"
    exec $@ && exit 0
fi

SCRIPTS=$(realpath "$(dirname "${0}")")

source "${SCRIPTS}/env.sh"

# Change owner of directories.
MOCK_CERTS_PATH="/etc/ssl/mock"
LOGS_SUPERVISORD_PATH="/var/log/supervisord"

USER_ID=1000

sudo mkdir -p "${LOGS_SUPERVISORD_PATH}"

sudo chown -R \
  "${USER_ID}:${USER_ID}" \
  "${GW_WIREMOCK_PATH}" \
  "${GW_CONTRACTS_PATH}" \
  "${LOGS_SUPERVISORD_PATH}" \
  "${MOCK_CERTS_PATH}"

# Setup logger.
LOG="/var/log/wiremock"

## remove 'imklog' because no need to monitor kernel events.
if [ ! -f "/var/run/rsyslogd.pid" ]; then
    sudo sed -i '/imklog/s/^/#/' /etc/rsyslog.conf && sudo rsyslogd
fi
sudo touch "${LOG}" && tail -f ${LOG} &

# Include utilities.
source "${SCRIPTS}/other/log.sh"

# Run grpc-wiremock entrypoint.
source "${SCRIPTS}/entrypoint.sh"
