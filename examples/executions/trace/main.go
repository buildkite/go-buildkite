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
	apiToken    = kingpin.Flag("token", "API token").Required().String()
	org         = kingpin.Flag("org", "Organization slug").Required().String()
	slug        = kingpin.Flag("slug", "Test suite slug").Required().String()
	executionID = kingpin.Flag("executionID", "Execution ID").Required().String()
	view        = kingpin.Flag("view", "Trace view: summary or full").Default("summary").Enum("summary", "full")
)

func main() {
	kingpin.Parse()

	client, err := buildkite.NewOpts(buildkite.WithTokenAuth(*apiToken))
	if err != nil {
		log.Fatalf("creating buildkite API client failed: %v", err)
	}

	trace, _, err := client.Executions.GetTrace(context.Background(), *org, *slug, *executionID, &buildkite.ExecutionTraceOptions{
		View: buildkite.TraceView(*view),
	})
	if err != nil {
		log.Fatalf("Getting trace for execution %s failed: %s", *executionID, err)
	}

	data, err := json.MarshalIndent(trace, "", "\t")
	if err != nil {
		log.Fatalf("json encode failed: %s", err)
	}

	_, _ = fmt.Fprintf(os.Stdout, "%s", string(data))
}
