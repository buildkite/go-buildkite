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

	pipelineTemplateUpdate := buildkite.PipelineTemplateCreateUpdate{Description: "Production pipeline template uploader"}

	cluster, _, err := client.PipelineTemplates.Update(context.Background(), *org, *templateUUID, pipelineTemplateUpdate)
	if err != nil {
		log.Fatalf("Updating cluster %s failed: %s", *templateUUID, err)
	}

	fmt.Printf("Updated cluster %s: new description: %s\n", *templateUUID, cluster.Description)
}
