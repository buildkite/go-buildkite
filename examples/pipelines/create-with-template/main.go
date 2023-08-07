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

	pc := buildkite.CreatePipeline{
		Name:                 *buildkite.String("my-great-pipeline"),
		Repository:           *buildkite.String("git@github.com:my_great_org/my_great_repo.git"),
		DefaultBranch:        *buildkite.String("main"),
		PipelineTemplateUuid: *buildkite.String("fc18d3c1-ea62-4091-b84a-0cbf1c0252b5"),
	}

	pipeline, _, err := client.Pipelines.Create(*org, &pc)

	if err != nil {
		log.Fatalf("Creating pipeline with pipeline template failed: %s", err)
	}

	data, err := json.MarshalIndent(pipeline, "", "\t")

	if err != nil {
		log.Fatalf("json encode failed: %s", err)
	}

	fmt.Fprintf(os.Stdout, "%s", string(data))
}
