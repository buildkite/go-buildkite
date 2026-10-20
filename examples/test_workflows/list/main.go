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
	opt := &buildkite.TestWorkflowListOptions{ListOptions: buildkite.ListOptions{PerPage: 100}}
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "\t")
	for {
		workflows, resp, err := client.TestWorkflows.List(context.Background(), *org, *suite, opt)
		if err != nil {
			log.Fatal(err)
		}
		if err := encoder.Encode(workflows); err != nil {
			log.Fatal(err)
		}
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
}
