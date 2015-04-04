DEPS = $(go list -f '{{range .TestImports}}{{.}} {{end}}' ./...)

all: deps test

deps:
	go get -d -v ./...
	echo $(DEPS) | xargs -n1 go get -d

test: deps
	go list ./... | xargs -n1 go test -timeout=3s -v

.PHONY: all deps test