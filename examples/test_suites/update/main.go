package main

import (
	"context"
	"fmt"
	"log"

	"github.com/buildkite/go-buildkite/v4"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	apiToken         = kingpin.Flag("token", "API token").Required().String()
	org              = kingpin.Flag("org", "Orginization slug").Required().String()
	suite            = kingpin.Flag("suite", "Test suite slug").Required().String()
	newDefaultBranch = kingpin.Flag("new-default-branch", "New default branch").String()
)

func main() {
	kingpin.Parse()

	client, err := buildkite.NewOpts(buildkite.WithTokenAuth(*apiToken))
	if err != nil {
		log.Fatalf("creating buildkite API client failed: %v", err)
	}

	suiteUpdate := buildkite.TestSuite{DefaultBranch: *newDefaultBranch}

	_, _, err = client.TestSuites.Update(context.Background(), *org, *suite, suiteUpdate)
	if err != nil {
		log.Fatalf("Updating test suite failed: %s", err)
	}

	fmt.Println("Updated test suite with new default branch: ", *newDefaultBranch)
}
