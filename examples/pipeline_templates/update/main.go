package main

import (
	"context"
	"fmt"
	"log"

	"github.com/buildkite/go-buildkite/v3/buildkite"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	apiToken     = kingpin.Flag("token", "API token").Required().String()
	org          = kingpin.Flag("org", "Orginization slug").Required().String()
	templateUUID = kingpin.Flag("templateUUID", "Cluster UUID").Required().String()
	debug        = kingpin.Flag("debug", "Enable debugging").Bool()
)

func main() {
	kingpin.Parse()

	client, err := buildkite.NewOpts(buildkite.WithTokenAuth(*apiToken))
	if err != nil {
		log.Fatalf("creating buildkite API client failed: %v", err)
	}

	pipelineTemplateUpdate := buildkite.PipelineTemplateCreateUpdate{
		Description: "Production pipeline template uploader",
	}

	resp, err := client.PipelineTemplates.Update(context.Background(), *org, *templateUUID, pipelineTemplateUpdate)
	if err != nil {
		log.Fatalf("Updating cluster %s failed: %s", *templateUUID, err)
	}

	fmt.Println(resp.StatusCode)
}
