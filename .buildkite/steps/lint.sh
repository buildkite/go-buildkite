#!/usr/bin/env sh

set -euf

tmpdir="$(mktemp -d)"
trap 'rm -rf "$tmpdir"' EXIT

echo --- :go: Checking go mod tidiness...
git diff --binary --no-ext-diff > "$tmpdir/before-go-mod-tidy.diff"
go mod tidy
git diff --binary --no-ext-diff > "$tmpdir/after-go-mod-tidy.diff"
if ! cmp -s "$tmpdir/before-go-mod-tidy.diff" "$tmpdir/after-go-mod-tidy.diff"; then
  echo ^^^ +++
  echo "The go.mod or go.sum files are out of sync with the source code"
  echo "Please run \`go mod tidy\` locally, and commit the result."

  exit 1
fi

echo --- :go: Checking go formatting...
git diff --binary --no-ext-diff > "$tmpdir/before-go-fmt.diff"
go fmt ./...
git diff --binary --no-ext-diff > "$tmpdir/after-go-fmt.diff"
if ! cmp -s "$tmpdir/before-go-fmt.diff" "$tmpdir/after-go-fmt.diff"; then
  echo ^^^ +++
  echo "Files have not been formatted with gofmt."
  echo "Fix this by running \`go fmt ./...\` locally, and committing the result."

  exit 1
fi

echo --- :go: Running golangci-lint...
golangci-lint run

echo --- :go: Checking code generation...
git diff --binary --no-ext-diff > "$tmpdir/before-go-generate.diff"
go generate ./...
git diff --binary --no-ext-diff > "$tmpdir/after-go-generate.diff"
if ! cmp -s "$tmpdir/before-go-generate.diff" "$tmpdir/after-go-generate.diff"; then
  echo ^^^ +++
  echo "Generated code is out of date."
  echo "Please run \`go generate ./...\` locally, and commit the result."

  exit 1
fi

echo +++ Everything is clean and tidy! 🎉
