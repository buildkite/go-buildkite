# buildkite-go [![GoDoc](https://img.shields.io/badge/godoc-Reference-brightgreen.svg?style=flat)](http://godoc.org/github.com/buildkite/go-buildkite)

A [Go](http://golang.org) library and client for the [Buildkite API](https://buildkite.com/docs/api). This project draws a lot of it's structure and testing methods from [go-github](https://github.com/google/go-github).

# Usage

To get the package, execute:

```
go get github.com/buildkite/go-buildkite/buildkite
```

Simple shortened example for listing all pipelines is provided below, see examples for more.

```go
import (
    "github.com/buildkite/go-buildkite/buildkite"
)
...

config, err := buildkite.NewTokenConfig(*apiToken, false)

if err != nil {
	log.Fatalf("client config failed: %s", err)
}

client := buildkite.NewClient(config.Client())

pipelines, _, err := client.Pipelines.List(*org, nil)

```

Note: not everything in the API is present here just yet—if you need something please make an issue or submit a pull request.

# License

This library is distributed under the BSD-style license found in the LICENSE file.

