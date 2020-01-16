all: test

test:
	go test -timeout=3s -v ./...

.PHONY: all test
