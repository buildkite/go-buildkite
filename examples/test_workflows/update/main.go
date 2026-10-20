package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

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
	// For a score-based monitor: change only the alarm threshold, leaving
	// filters, alarm actions, and other monitor fields unchanged.
	input := buildkite.TestWorkflowUpdate{
		Monitor: buildkite.Some(buildkite.TestWorkflowMonitorInput{
			AlarmThreshold: buildkite.Some(0.2),
		}),
	}
	// To also clear recovery actions, set an empty slice instead of leaving
	// RecoverActions unset:
	// input.RecoverActions = buildkite.Some([]buildkite.TestWorkflowActionInput{})
	result, _, err := client.TestWorkflows.Update(context.Background(), *org, *suite, *workflow, input)
	if err != nil {
		log.Fatal(err)
	}
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "\t")
	if err := encoder.Encode(result); err != nil {
		log.Fatal(err)
	}
}
