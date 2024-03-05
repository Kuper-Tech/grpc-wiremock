#!/bin/bash

set -euo pipefail

killall tail || true
ps | grep '[s]cript' | awk '{print $1}' | xargs -r -n1 kill || true


SCRIPTS=$(realpath "$(dirname "${0}")")/../scripts

$SCRIPTS/init.sh
