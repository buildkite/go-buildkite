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
	slug     = kingpin.Flag("slug", "Test suite slug").Required().String()
	period   = kingpin.Flag("period", "Aggregation period").Default("7days").String()
	branch   = kingpin.Flag("branch", "Branch filter expression").String()
	sortBy   = kingpin.Flag("sort-by", "Metric used to sort tests").Default("reliability").String()
	order    = kingpin.Flag("order", "Sort direction").Default("asc").String()
)

func main() {
	kingpin.Parse()

	client, err := buildkite.NewOpts(buildkite.WithTokenAuth(*apiToken))
	if err != nil {
		log.Fatalf("creating buildkite API client failed: %v", err)
	}

	tests, _, err := client.Tests.List(context.Background(), *org, *slug, &buildkite.TestsListOptions{
		Period: *period,
		Branch: *branch,
		SortBy: *sortBy,
		Order:  *order,
	})
	if err != nil {
		log.Fatalf("listing tests for suite %s failed: %s", *slug, err)
	}

	data, err := json.MarshalIndent(tests, "", "\t")
	if err != nil {
		log.Fatalf("json encode failed: %s", err)
	}

	_, _ = fmt.Fprintf(os.Stdout, "%s", string(data))
}
