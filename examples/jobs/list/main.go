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
	apiToken           = kingpin.Flag("token", "API token").Required().String()
	baseURL            = kingpin.Flag("base-url", "Buildkite API base URL").Default(buildkite.DefaultBaseURL).String()
	org                = kingpin.Flag("org", "Organization slug").Required().String()
	pipeline           = kingpin.Flag("pipeline", "Pipeline slug").Required().String()
	build              = kingpin.Flag("build", "Build number").Required().String()
	state              = kingpin.Flag("state", "Job state filter. Can be provided multiple times.").Strings()
	includeRetriedJobs = kingpin.Flag("include-retried-jobs", "Include retried jobs").Default("true").Bool()
	perPage            = kingpin.Flag("per-page", "Number of jobs to return per page").Int()
	after              = kingpin.Flag("after", "Cursor to list jobs after").String()
	before             = kingpin.Flag("before", "Cursor to list jobs before").String()
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

	opt := &buildkite.JobsListOptions{
		State:              *state,
		IncludeRetriedJobs: includeRetriedJobs,
		PerPage:            *perPage,
		After:              *after,
		Before:             *before,
	}

	jobs, _, err := client.Jobs.ListByBuild(context.Background(), *org, *pipeline, *build, opt)
	if err != nil {
		log.Fatalf("listing jobs for build %s in pipeline %s failed: %s", *build, *pipeline, err)
	}

	data, err := json.MarshalIndent(jobs, "", "\t")
	if err != nil {
		log.Fatalf("json encode failed: %s", err)
	}

	_, _ = fmt.Fprintf(os.Stdout, "%s", string(data))
}
