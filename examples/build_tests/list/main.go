package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/buildkite/go-buildkite/v4"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	apiToken  = kingpin.Flag("token", "Buildkite API token").Envar("BUILDKITE_API_TOKEN").Required().String()
	org       = kingpin.Flag("org", "Organization slug").Required().String()
	buildUUID = kingpin.Flag("build-id", "Build UUID").Required().String()
	result    = kingpin.Flag("result", "Filter by execution result, for example failed, ^failed, passed, or ^passed.").String()
	state     = kingpin.Flag("state", "Filter by test state, for example enabled or muted.").String()
)

func main() {
	kingpin.Parse()

	client, err := buildkite.NewOpts(buildkite.WithTokenAuth(*apiToken))
	if err != nil {
		log.Fatalf("creating buildkite API client failed: %v", err)
	}

	opts := &buildkite.BuildTestsListOptions{
		Result:  *result,
		State:   *state,
		Include: "executions",
	}

	buildTests, _, err := client.BuildTests.List(context.Background(), *org, *buildUUID, opts)
	if err != nil {
		log.Fatalf("listing build tests for build %s failed: %s", *buildUUID, err)
	}

	data, err := json.MarshalIndent(buildTests, "", "\t")
	if err != nil {
		log.Fatalf("json encode failed: %s", err)
	}

	_, _ = fmt.Fprintf(os.Stdout, "%s", string(data))
}
