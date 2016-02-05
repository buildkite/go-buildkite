# buildkite-go [![GoDoc](https://img.shields.io/badge/godoc-Reference-brightgreen.svg?style=flat)](http://godoc.org/github.com/wolfeidau/go-buildkite) [![Build status](https://badge.buildkite.com/f7561b01d3f2886b819d0825464bf9a3c90cd0d0a1a96a517d.svg)](https://buildkite.com/mark-at-wolfe-dot-id-dot-au/go-buildkite)

A [golang](http://golang.org) client for the [buildkite](https://buildkite.com/) API. This project draws a lot of it's structure and testing methods from [go-github](https://github.com/google/go-github).

# Usage

To get the package, execute:

```
go get gopkg.in/wolfeidau/go-buildkite.v1
```

Simple shortened example for listing all projects is provided below, see examples for more.

```go
import (
    "gopkg.in/wolfeidau/go-buildkite.v1"
)
...

config, err := buildkite.NewTokenConfig(*apiToken)

if err != nil {
	log.Fatalf("client config failed: %s", err)
}

client := buildkite.NewClient(config.Client())

projects, _, err := client.Projects.List(*org, nil)

```

# Disclaimer

Currently not everything in the [buildkite API](https://buildkite.com/docs/api/) is present here YET. If
you find something you need is missing please make a pull request.

# License

This library is distributed under the BSD-style license found in the LICENSE file.
