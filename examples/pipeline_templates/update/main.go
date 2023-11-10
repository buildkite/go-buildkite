package main

import (
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

	config, err := buildkite.NewTokenConfig(*apiToken, *debug)

	if err != nil {
		log.Fatalf("client config failed: %s", err)
	}

	client := buildkite.NewClient(config.Client())

	pipelineTemplateUpdate := buildkite.PipelineTemplateCreateUpdate{
		Description: buildkite.String("Production pipeline template uploader"),
	}

	resp, err := client.PipelineTemplates.Update(*org, *templateUUID, &pipelineTemplateUpdate)

	if err != nil {
		log.Fatalf("Updating cluster %s failed: %s", *templateUUID, err)
	}

	fmt.Println(resp.StatusCode)
}
