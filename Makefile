GOPATH := ${PWD}/.gopath
STAGING_PATH = ${PWD}/.gopath/src/github.com/wolfeidau
DEPS = $(go list -f '{{range .TestImports}}{{.}} {{end}}' ./...)

all: deps test

deps:
	@mkdir -p ${STAGING_PATH}
	@ln -s ${PWD} ${STAGING_PATH} || true
	@cd ${STAGING_PATH}/go-buildkite
	go get -d -v ./...

test: deps
	@cd ${STAGING_PATH}/go-buildkite
	go test -timeout=3s -v ./...

clean:
	rm -rf ${PWD}/.gopath || true

.PHONY: all deps test clean