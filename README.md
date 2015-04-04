# buildkite-go [![GoDoc](https://img.shields.io/badge/godoc-Reference-brightgreen.svg?style=flat)](http://godoc.org/github.com/wolfeidau/go-buildkite)

A [golang](http://golang.org) client for the [buildkite](https://buildkite.com/) API. This project draws a lot of it's structure and testing methods from [go-github](https://github.com/google/go-github).

# Usage

Simple example for listing all projects is provided below, see examples for more.

```golang

config, err := buildkite.NewTokenConfig(*apiToken)

if err != nil {
	log.Fatalf("client config failed: %s", err)
}

client := buildkite.NewClient(config.Client())

projects, _, err := client.Projects.List(*org, nil)

```

# Disclaimer

This is currently very early release, not everything in the [buildkite API](https://buildkite.com/docs/api/) is present here YET.

# License

This library is distributed under the BSD-style license found in the LICENSE file.