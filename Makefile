all: test

fmt: 
	gofmt -w .

test:
	go test -timeout=3s -v ./...

analytics-openapi.yaml: FORCE
	curl -s -o $@ https://api.buildkite.com/v2/analytics/openapi.yaml

.PHONY: all test

FORCE:
