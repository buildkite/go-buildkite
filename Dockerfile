FROM golang:1.6

RUN mkdir -p /go/src/github.com/buildkite/go-buildkite
WORKDIR /go/src/github.com/buildkite/go-buildkite
