FROM golang:1.7.3

RUN mkdir -p /go/src/github.com/buildkite/go-buildkite
ADD . /go/src/github.com/buildkite/go-buildkite

WORKDIR /go/src/github.com/buildkite/go-buildkite