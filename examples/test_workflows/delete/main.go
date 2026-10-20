package main

import (
	"context"
	"log"

	"github.com/alecthomas/kingpin/v2"
	"github.com/buildkite/go-buildkite/v5"
)

var (
	apiToken = kingpin.Flag("token", "API token").Required().String()
	org      = kingpin.Flag("org", "Organization slug").Required().String()
	suite    = kingpin.Flag("suite", "Test suite slug or UUID").Required().String()
	workflow = kingpin.Flag("workflow", "Workflow UUID").Required().String()
)

func main() {
	kingpin.Parse()
	client, err := buildkite.NewClient(buildkite.WithTokenAuth(*apiToken))
	if err != nil {
		log.Fatal(err)
	}
	if _, err := client.TestWorkflows.Delete(context.Background(), *org, *suite, *workflow); err != nil {
		log.Fatal(err)
	}
}
