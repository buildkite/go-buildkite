package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

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

	config, err := buildkite.NewTokenConfig(*apiToken, *debug)

	if err != nil {
		log.Fatalf("client config failed: %s", err)
	}

	client := buildkite.NewClient(config.Client())

	pipelineTemplateCreate := buildkite.PipelineTemplateCreateUpdate{
		Name:          buildkite.String("Production Pipeline uploader"),
		Description:   buildkite.String("Production pipeline upload template"),
		Configuration: buildkite.String("steps:\n  - label: \":pipeline:\"\n    command: \"buildkite-agent pipeline upload .buildkite/pipeline-production.yml\"\n"),
		Available:     buildkite.Bool(true),
	}

	pipelineTemplate, _, err := client.PipelineTemplates.Create(*org, &pipelineTemplateCreate)

	if err != nil {
		log.Fatalf("Creating pipeline template failed: %s", err)
	}

	data, err := json.MarshalIndent(pipelineTemplate, "", "\t")

	if err != nil {
		log.Fatalf("json encode failed: %s", err)
	}

	fmt.Fprintf(os.Stdout, "%s", string(data))
}
