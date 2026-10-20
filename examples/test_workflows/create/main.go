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
)

func main() {
	kingpin.Parse()
	client, err := buildkite.NewClient(buildkite.WithTokenAuth(*apiToken))
	if err != nil {
		log.Fatal(err)
	}
	workflow, _, err := client.TestWorkflows.Create(context.Background(), *org, *suite, buildkite.TestWorkflowCreate{
		Monitor: buildkite.TestWorkflowMonitorInput{
			Type:             buildkite.Some("transition_count"),
			WindowSize:       buildkite.Some(100),
			AlarmThreshold:   buildkite.Some(0.1),
			RecoverThreshold: buildkite.Some(0.0),
		},
		Filters: []buildkite.TestWorkflowFilter{
			{Key: "scm.branch", Operator: "is", Value: "main"},
		},
		AlarmActions: []buildkite.TestWorkflowActionInput{
			{Type: "add_label", Label: buildkite.Some("flaky")},
		},
		RecoverActions: []buildkite.TestWorkflowActionInput{
			{Type: "remove_label", Label: buildkite.Some("flaky")},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "\t")
	if err := encoder.Encode(workflow); err != nil {
		log.Fatal(err)
	}
}
