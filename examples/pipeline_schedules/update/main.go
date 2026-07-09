package main

import (
	"context"
	"fmt"
	"log"

	"github.com/alecthomas/kingpin/v2"
	"github.com/buildkite/go-buildkite/v5"
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

	scheduleUpdate := buildkite.UpdatePipelineSchedule{
		Label:   buildkite.Some("Nightly build (paused)"),
		Enabled: buildkite.Some(false),
	}

	schedule, _, err := client.PipelineSchedules.Update(context.Background(), *org, *pipeline, *scheduleID, scheduleUpdate)
	if err != nil {
		log.Fatalf("Updating pipeline schedule %s failed: %s", *scheduleID, err)
	}

	fmt.Printf("Updated schedule %s: label=%q enabled=%t\n", *scheduleID, schedule.Label, schedule.Enabled)
}
