package main

import (
	"context"
	"fmt"
	"log"

	"github.com/buildkite/go-buildkite/v4"

	"github.com/alecthomas/kingpin/v2"
)

var (
	apiToken   = kingpin.Flag("token", "API token").Required().String()
	org        = kingpin.Flag("org", "Organization slug").Required().String()
	pipeline   = kingpin.Flag("pipeline", "Pipeline slug").Required().String()
	scheduleID = kingpin.Flag("scheduleID", "Pipeline schedule UUID").Required().String()
)

func main() {
	kingpin.Parse()

	client, err := buildkite.NewOpts(buildkite.WithTokenAuth(*apiToken))
	if err != nil {
		log.Fatalf("creating buildkite API client failed: %v", err)
	}

	resp, err := client.PipelineSchedules.Delete(context.Background(), *org, *pipeline, *scheduleID)
	if err != nil {
		log.Fatalf("Deleting pipeline schedule %s failed: %s", *scheduleID, err)
	}

	fmt.Println(resp.StatusCode)
}
