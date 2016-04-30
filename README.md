# buildkite-go [![GoDoc](https://img.shields.io/badge/godoc-Reference-brightgreen.svg?style=flat)](http://godoc.org/github.com/buildkite/go-buildkite) [![Build status](https://badge.buildkite.com/56c797ece642ce625a93af462fe648aa0488b68c2d8846cfd2.svg)](https://buildkite.com/buildkite/buildkite-golang-client)

A [Go](http://golang.org) library and client for the [Buildkite API](https://buildkite.com/docs/api). This project draws a lot of it's structure and testing methods from [go-github](https://github.com/google/go-github).

# Usage

To get the package, execute:

```
go get gopkg.in/buildkite/go-buildkite.v2
```

Simple shortened example for listing all pipelines is provided below, see examples for more.

```go
import (
    "gopkg.in/buildkite/go-buildkite.v2"
)
...

config, err := buildkite.NewTokenConfig(*apiToken)

if err != nil {
	log.Fatalf("client config failed: %s", err)
}

client := buildkite.NewClient(config.Client())

pipelines, _, err := client.Pipelines.List(*org, nil)

```

Note: not everything in the API is present here just yet—if you need something please make an issue or submit a pull request.

# License

This library is distributed under the BSD-style license found in the LICENSE file.
