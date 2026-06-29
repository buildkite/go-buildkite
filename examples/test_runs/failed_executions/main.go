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

type failedExecutionsForRun struct {
	SuiteSlug        string                      `json:"suite_slug,omitempty"`
	RunID            string                      `json:"run_id,omitempty"`
	FailedExecutions []buildkite.FailedExecution `json:"failed_executions,omitempty"`
}

var (
	apiToken               = kingpin.Flag("token", "API token").Required().String()
	org                    = kingpin.Flag("org", "Organization slug").Required().String()
	pipeline               = kingpin.Flag("pipeline", "Pipeline slug").Required().String()
	buildNumber            = kingpin.Flag("build", "Build number").Required().String()
	suite                  = kingpin.Flag("suite", "Only fetch failed executions for this test suite slug").String()
	includeFailureExpanded = kingpin.Flag("include-failure-expanded", "Include expanded failure details").Bool()
	page                   = kingpin.Flag("page", "Page of failed executions to retrieve").Int()
	perPage                = kingpin.Flag("per-page", "Number of failed executions to include per page").Int()
)

func main() {
	kingpin.Parse()

	client, err := buildkite.NewOpts(buildkite.WithTokenAuth(*apiToken))
	if err != nil {
		log.Fatalf("creating buildkite API client failed: %v", err)
	}

	ctx := context.Background()
	build, _, err := client.Builds.Get(ctx, *org, *pipeline, *buildNumber, &buildkite.BuildGetOptions{
		IncludeTestEngine: true,
	})
	if err != nil {
		log.Fatalf("Getting build %s of pipeline %s failed: %s", *buildNumber, *pipeline, err)
	}

	if build.TestEngine == nil || len(build.TestEngine.Runs) == 0 {
		log.Fatalf("build %s of pipeline %s has no Test Analytics runs", *buildNumber, *pipeline)
	}

	var results []failedExecutionsForRun
	for _, run := range build.TestEngine.Runs {
		if *suite != "" && run.Suite.Slug != *suite {
			continue
		}

		failedExecutions, _, err := client.TestRuns.GetFailedExecutions(ctx, *org, run.Suite.Slug, run.ID, &buildkite.FailedExecutionsOptions{
			IncludeFailureExpanded: *includeFailureExpanded,
			Page:                   *page,
			PerPage:                *perPage,
		})
		if err != nil {
			log.Fatalf("Getting failed executions for test run %s of suite %s failed: %s", run.ID, run.Suite.Slug, err)
		}

		results = append(results, failedExecutionsForRun{
			SuiteSlug:        run.Suite.Slug,
			RunID:            run.ID,
			FailedExecutions: failedExecutions,
		})
	}

	if len(results) == 0 {
		log.Fatalf("build %s of pipeline %s has no Test Analytics runs for suite %s", *buildNumber, *pipeline, *suite)
	}

	data, err := json.MarshalIndent(results, "", "\t")
	if err != nil {
		log.Fatalf("json encode failed: %s", err)
	}

	_, _ = fmt.Fprintf(os.Stdout, "%s", string(data))
}
