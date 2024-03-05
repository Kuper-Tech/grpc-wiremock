#!/bin/bash

set -euo pipefail

ETC_SUPERVISORD_PATH="/etc/supervisord"

supervisord -d -c "${ETC_SUPERVISORD_PATH}/supervisord.conf"
