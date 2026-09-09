package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/alecthomas/kingpin/v2"
	"github.com/buildkite/go-buildkite/v5"
)

var (
	apiToken  = kingpin.Flag("token", "API token").Required().String()
	baseURL   = kingpin.Flag("base-url", "Buildkite API base URL").Default(buildkite.DefaultBaseURL).String()
	org       = kingpin.Flag("org", "Organization slug").Required().String()
	buildUUID = kingpin.Flag("build-id", "Build UUID").Required().String()
	limit     = kingpin.Flag("limit", "Maximum number of executions to return (API default: 20)").Int()
	withTrace = kingpin.Flag("with-trace", "Also fetch the trace of the slowest execution if it has one").Bool()
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

	ctx := context.Background()

	executions, _, err := client.Executions.ListSlowestByBuild(ctx, *org, *buildUUID, &buildkite.SlowestExecutionsOptions{Limit: *limit})
	if err != nil {
		log.Fatalf("listing slowest executions for build %s failed: %v", *buildUUID, err)
	}

	printJSON(executions)

	if !*withTrace || len(executions) == 0 {
		return
	}

	slowest := executions[0]
	if !slowest.HasTrace {
		log.Printf("execution %s has no trace within retention", slowest.ID)
		return
	}

	trace, _, err := client.Executions.GetTrace(ctx, *org, slowest.SuiteSlug, slowest.ID, nil)
	if err != nil {
		log.Fatalf("getting trace for execution %s failed: %v", slowest.ID, err)
	}

	printJSON(trace)
}

func printJSON(v any) {
	data, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		log.Fatalf("json encode failed: %v", err)
	}

	fmt.Println(string(data))
}
