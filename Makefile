all: test

fmt: 
	gofmt -w .

test:
	go test -timeout=3s -v ./...

.PHONY: all test
