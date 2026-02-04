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
	apiToken = kingpin.Flag("token", "API token").Required().String()
	org      = kingpin.Flag("org", "Organization slug").Required().String()
	pipeline = kingpin.Flag("pipeline", "Pipeline slug").Required().String()
	build    = kingpin.Flag("build", "Build number").Required().String()
	job      = kingpin.Flag("job", "Job ID").Required().String()
	priority = kingpin.Flag("priority", "New priority for the job").Required().Int()
)

func main() {
	kingpin.Parse()

	client, err := buildkite.NewOpts(buildkite.WithTokenAuth(*apiToken))
	if err != nil {
		log.Fatalf("creating buildkite API client failed: %v", err)
	}

	updatedJob, _, err := client.Jobs.ReprioritizeJob(
		context.Background(),
		*org,
		*pipeline,
		*build,
		*job,
		&buildkite.JobReprioritizationOptions{Priority: *priority},
	)
	if err != nil {
		log.Fatalf("Reprioritizing job %s failed: %s", *job, err)
	}

	data, err := json.MarshalIndent(updatedJob, "", "\t")
	if err != nil {
		log.Fatalf("json encode failed: %s", err)
	}

	_, _ = fmt.Fprintf(os.Stdout, "%s", string(data))
}
