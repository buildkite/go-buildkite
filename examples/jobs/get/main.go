package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/alecthomas/kingpin/v2"
	"github.com/buildkite/go-buildkite/v5"
)

var (
	apiToken = kingpin.Flag("token", "API token").Required().String()
	baseURL  = kingpin.Flag("base-url", "Buildkite API base URL").Default(buildkite.DefaultBaseURL).String()
	org      = kingpin.Flag("org", "Organization slug").Required().String()
	pipeline = kingpin.Flag("pipeline", "Pipeline slug").Required().String()
	build    = kingpin.Flag("build", "Build number").Required().String()
	job      = kingpin.Flag("job", "Job ID").Required().String()
)

func main() {
	kingpin.Parse()

	client, err := buildkite.NewOpts(
		buildkite.WithTokenAuth(*apiToken),
		buildkite.WithBaseURL(*baseURL),
	)
	if err != nil {
		log.Fatalf("creating buildkite API client failed: %v", err)
	}

	fetchedJob, _, err := client.Jobs.GetJob(context.Background(), *org, *pipeline, *build, *job)
	if err != nil {
		log.Fatalf("getting job %s failed: %s", *job, err)
	}

	data, err := json.MarshalIndent(fetchedJob, "", "\t")
	if err != nil {
		log.Fatalf("json encode failed: %s", err)
	}

	_, _ = fmt.Fprintf(os.Stdout, "%s", string(data))
}
