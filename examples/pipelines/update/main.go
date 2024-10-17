package main

import (
	"context"
	"fmt"
	"log"

	"github.com/buildkite/go-buildkite/v3"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	apiToken       = kingpin.Flag("token", "API token").Required().String()
	org            = kingpin.Flag("org", "Organization slug").Required().String()
	pipelineSlug   = kingpin.Flag("pipeline", "Pipeline slug").Required().String()
	newDescription = kingpin.Flag("description", "New pipeline description").Required().String()
	debug          = kingpin.Flag("debug", "Enable debugging").Bool()
)

func main() {
	kingpin.Parse()

	client, err := buildkite.NewOpts(buildkite.WithTokenAuth(*apiToken))
	if err != nil {
		log.Fatalf("creating buildkite API client failed: %v", err)
	}

	_, _, err = client.Pipelines.Update(context.Background(), *org, *pipelineSlug, buildkite.UpdatePipeline{Description: *newDescription})
	if err != nil {
		log.Fatalf("Updating pipeline failed: %s", err)
	}

	fmt.Println("Changed pipeline description to " + *newDescription)
}
