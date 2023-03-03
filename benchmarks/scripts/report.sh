#!/bin/bash

set -euo pipefail

PATH_TO_REPORTS="${1-output}"

reports='[]'

for path in $(find "${PATH_TO_REPORTS}" -type d -name "api_*"); do
  api_count="$(echo "${path}" | grep -Eo '[0-9]+$')"

  report=$(echo $(find "${path}" -name '*.txt' -exec tail -1 {} \; | awk '{print $2}' | xargs) \
    | jq --arg api_count "${api_count}" -s "{ min_response_time_api_count_$api_count:min, max_response_time_api_count_$api_count:max, avg_response_time_api_count_$api_count: (add/length), median_response_time_api_count_$api_count: (sort|.[(length/2|floor)]) }")

  reports=$(echo "${reports}" | jq ". + [${report}]")
done

echo "${reports}" | jq > "${PATH_TO_REPORTS}/report.json"
