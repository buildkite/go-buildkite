#!/usr/bin/env bash

set -Eeuo pipefail

tmpdir="$(mktemp -d)"
trap 'rm -rf "$tmpdir"' EXIT

if ! go build -o "$tmpdir/analytics-openapi-drift" ./.buildkite/steps/analytics-openapi-drift; then
  echo "Analytics OpenAPI checker failed to build (exit code 2)." >&2
  exit 2
fi

if output=$("$tmpdir/analytics-openapi-drift" 2>&1); then
  exit 0
else
  status=$?
fi

printf '%s\n' "$output"
if [[ "$status" -ne 1 ]]; then
  echo "Analytics OpenAPI checker failed (exit code $status)." >&2
  exit "$status"
fi

buildkite-agent annotate --style warning --context analytics-openapi-drift <<EOF
## Analytics OpenAPI drift detected

The Go response types no longer match [the Analytics OpenAPI specification](https://api.buildkite.com/v2/analytics/openapi.yaml).

~~~text
$output
~~~
EOF

exit 1
