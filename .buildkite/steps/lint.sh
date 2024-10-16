#!/usr/bin/env sh

set -eufo

echo --- :go: Checking go mod tidiness...
go mod tidy
if ! git diff --no-ext-diff --exit-code; then
  echo ^^^ +++
  echo "The go.mod or go.sum files are out of sync with the source code"
  echo "Please run \`go mod tidy\` locally, and commit the result."

  exit 1
fi

echo --- :go: Checking go formatting...
go fmt ./...
if ! git diff --no-ext-diff --exit-code; then
  echo ^^^ +++
  echo "Files have not been formatted with gofmt."
  echo "Fix this by running \`go fmt ./...\` locally, and committing the result."

  exit 1
fi

echo --- :go: Running golangci-lint...
golangci-lint run

echo --- :go: Checking code generation...
go generate ./...
if ! git diff --no-ext-diff --exit-code; then
  echo ^^^ +++
  echo "Generated code is out of date."
  echo "Please run \`go generate ./...\` locally, and commit the result."

  exit 1
fi


echo +++ Everything is clean and tidy! 🎉
