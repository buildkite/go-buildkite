package main

import (
	"context"
	"fmt"
	"log"

	"github.com/buildkite/go-buildkite/v4"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	apiToken = kingpin.Flag("token", "API token").Required().String()
	org      = kingpin.Flag("org", "Orginization slug").Required().String()
)

func main() {
	kingpin.Parse()

	client, err := buildkite.NewOpts(buildkite.WithTokenAuth(*apiToken))
	if err != nil {
		log.Fatalf("creating buildkite API client failed: %v", err)
	}

	limits, _, err := client.RateLimit.Get(context.Background(), *org)
	if err != nil {
		log.Fatalf("getting rate limits failed: %s", err)
	}

	fmt.Printf("Current rate limits:\n %v", limits)
}
