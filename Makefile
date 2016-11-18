GOPATH := $(shell pwd)/.gopath
STAGING_PATH = $(shell pwd)/.gopath/src/github.com/buildkite
DEPS = $(go list -f '{{range .TestImports}}{{.}} {{end}}' ./...)

all: deps test

deps:
	mkdir -p ${STAGING_PATH}
	go get -d -v ./...

test: deps
	cd ${STAGING_PATH}/go-buildkite
	go test -timeout=3s -v ./...

clean:
	rm -rf $(shell pwd)/.gopath || true

.PHONY: all deps test clean
