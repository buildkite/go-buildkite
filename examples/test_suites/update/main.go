package main

import (
	"fmt"
	"log"

	"github.com/buildkite/go-buildkite/v3/buildkite"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	apiToken = kingpin.Flag("token", "API token").Required().String()
	org      = kingpin.Flag("org", "Orginization slug").Required().String()
	debug    = kingpin.Flag("debug", "Enable debugging").Bool()
)

func main() {
	kingpin.Parse()

	client, err := buildkite.NewOpts(buildkite.WithTokenAuth(*apiToken))
	if err != nil {
		log.Fatalf("creating buildkite API client failed: %v", err)
	}

	suiteUpdate := buildkite.TestSuite{
		DefaultBranch: buildkite.String("test"),
		Slug:          buildkite.String("example-slug"),
	}

	resp, err := client.TestSuites.Update(*org, &suiteUpdate)

	if err != nil {
		log.Fatalf("Updating test suite failed: %s", err)
	}

	fmt.Println(resp.StatusCode)
}
