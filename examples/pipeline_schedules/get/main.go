package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/buildkite/go-buildkite/v5"

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

	schedule, _, err := client.PipelineSchedules.Get(context.Background(), *org, *pipeline, *scheduleID)
	if err != nil {
		log.Fatalf("Getting pipeline schedule %s failed: %s", *scheduleID, err)
	}

	data, err := json.MarshalIndent(schedule, "", "\t")
	if err != nil {
		log.Fatalf("json encode failed: %s", err)
	}

	_, _ = fmt.Fprintf(os.Stdout, "%s", string(data))
}
