package main

import (
	"context"
	"fmt"
	"log"

	"github.com/buildkite/go-buildkite/v4"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	apiToken     = kingpin.Flag("token", "API token").Required().String()
	org          = kingpin.Flag("org", "Orginization slug").Required().String()
	templateUUID = kingpin.Flag("templateUUID", "Cluster UUID").Required().String()
)

func main() {
	kingpin.Parse()

	client, err := buildkite.NewOpts(buildkite.WithTokenAuth(*apiToken))
	if err != nil {
		log.Fatalf("creating buildkite API client failed: %v", err)
	}

	resp, err := client.PipelineTemplates.Delete(context.Background(), *org, *templateUUID)

	if err != nil {
		log.Fatalf("Deleting pipeline template %s failed: %s", *templateUUID, err)
	}

	fmt.Println(resp.StatusCode)
}
