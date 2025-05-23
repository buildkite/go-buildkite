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

	pipelineTemplateCreate := buildkite.PipelineTemplateCreateUpdate{
		Name:          "Production Pipeline uploader",
		Description:   "Production pipeline upload template",
		Configuration: "steps:\n  - label: \":pipeline:\"\n    command: \"buildkite-agent pipeline upload .buildkite/pipeline-production.yml\"\n",
		Available:     true,
	}

	pipelineTemplate, _, err := client.PipelineTemplates.Create(context.Background(), *org, pipelineTemplateCreate)
	if err != nil {
		log.Fatalf("Creating pipeline template failed: %s", err)
	}

	data, err := json.MarshalIndent(pipelineTemplate, "", "\t")
	if err != nil {
		log.Fatalf("json encode failed: %s", err)
	}

	_, _ = fmt.Fprintf(os.Stdout, "%s", string(data))
}
