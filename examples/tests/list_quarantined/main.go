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
	state    = kingpin.Flag("state", "Quarantine state to list").Required().Enum("muted", "skipped")
	page     = kingpin.Flag("page", "Page of results to retrieve").Int()
	perPage  = kingpin.Flag("per-page", "Number of tests to include per page").Int()
)

func main() {
	kingpin.Parse()

	client, err := buildkite.NewOpts(buildkite.WithTokenAuth(*apiToken))
	if err != nil {
		log.Fatalf("creating buildkite API client failed: %v", err)
	}

	ctx := context.Background()
	opt := &buildkite.QuarantinedTestsListOptions{
		ListOptions: buildkite.ListOptions{Page: *page, PerPage: *perPage},
	}

	var (
		tests []buildkite.Test
		resp  *buildkite.Response
	)
	switch *state {
	case "muted":
		tests, resp, err = client.Tests.ListMuted(ctx, *org, *slug, opt)
	case "skipped":
		tests, resp, err = client.Tests.ListSkipped(ctx, *org, *slug, opt)
	}
	if err != nil {
		log.Fatalf("listing %s tests failed: %s", *state, err)
	}

	data, err := json.MarshalIndent(tests, "", "\t")
	if err != nil {
		log.Fatalf("json encode failed: %s", err)
	}

	_, _ = fmt.Fprintf(os.Stdout, "%s", string(data))

	// These endpoints paginate, so this is one page of the quarantine list.
	if resp.NextPage != 0 {
		_, _ = fmt.Fprintf(os.Stderr, "\nmore results: next page %d, last page %d\n", resp.NextPage, resp.LastPage)
	}
}
