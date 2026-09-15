#!/usr/bin/env bash

set -Eeuo pipefail

if output=$(go test -tags=analytics_openapi -run '^TestAnalyticsResponseTypesMatchOpenAPI$' -count=1 . 2>&1); then
  exit 0
fi

printf '%s\n' "$output"
buildkite-agent annotate --style warning --context analytics-openapi-drift <<EOF
## Analytics OpenAPI drift detected

The Go response types no longer match [the Analytics OpenAPI specification](https://api.buildkite.com/v2/analytics/openapi.yaml).

~~~text
$output
~~~
EOF

exit 1
