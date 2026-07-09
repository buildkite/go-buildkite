package main

import (
	"context"
	"fmt"
	"log"

	"github.com/alecthomas/kingpin/v2"
	"github.com/buildkite/go-buildkite/v5"
)

var (
	apiToken = kingpin.Flag("token", "API token").Required().String()
	org      = kingpin.Flag("org", "Organization slug").Required().String()
	slug     = kingpin.Flag("slug", "Test suite slug").Required().String()
)

func main() {
	kingpin.Parse()

	client, err := buildkite.NewOpts(buildkite.WithTokenAuth(*apiToken))
	if err != nil {
		log.Fatalf("creating buildkite API client failed: %v", err)
	}

	resp, err := client.TestSuites.Delete(context.Background(), *org, *slug)
	if err != nil {
		log.Fatalf("Deleting test suite failed: %s", err)
	}

	fmt.Println(resp.StatusCode)
}
