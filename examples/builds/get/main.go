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
	apiToken = kingpin.Flag("token", "API token").Required().String()
	org      = kingpin.Flag("org", "Organization slug").Required().String()
	slug     = kingpin.Flag("slug", "Pipeline slug").Required().String()
	number   = kingpin.Flag("number", "Build number").Required().String()
)

func main() {
	kingpin.Parse()

	client, err := buildkite.NewOpts(buildkite.WithTokenAuth(*apiToken))
	if err != nil {
		log.Fatalf("creating buildkite API client failed: %v", err)
	}

	build, _, err := client.Builds.Get(context.Background(), *org, *slug, *number, &buildkite.BuildGetOptions{
		BuildsListOptions: buildkite.BuildsListOptions{
			ExcludeJobs:     true,
			ExcludePipeline: true,
		},
	})
	if err != nil {
		log.Fatalf("Getting build %s of pipeline %s failed: %s", *number, *slug, err)
	}

	data, err := json.MarshalIndent(build, "", "\t")
	if err != nil {
		log.Fatalf("json encode failed: %s", err)
	}

	_, _ = fmt.Fprintf(os.Stdout, "%s", string(data))
}
