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
	org      = kingpin.Flag("org", "Orginization slug").Required().String()
)

func main() {
	kingpin.Parse()

	client, err := buildkite.NewOpts(buildkite.WithTokenAuth(*apiToken))
	if err != nil {
		log.Fatalf("creating buildkite API client failed: %v", err)
	}

	createPipeline := buildkite.CreatePipeline{
		Name:          "my-great-pipeline",
		Repository:    "git@github.com:my_great_org/my_great_repo2.git",
		Configuration: "env:\n \"FOO\": \"bar\"\nsteps:\n - command: \"script/release.sh\"\n   \"name\": \"Build 📦\"",
		Tags:          []string{"great", "pipeline"},
		Description:   "This ia a great pipeline!",
		ProviderSettings: &buildkite.GitHubSettings{
			TriggerMode: "code",
		},
	}

	pipeline, _, err := client.Pipelines.Create(context.Background(), *org, createPipeline)

	if err != nil {
		log.Fatalf("Updating pipeline failed: %s", err)
	}

	data, err := json.MarshalIndent(pipeline, "", "\t")

	if err != nil {
		log.Fatalf("json encode failed: %s", err)
	}

	_, _ = fmt.Fprintf(os.Stdout, "%s", string(data))
}
