package main

import (
	"fmt"
	"log"

	"github.com/buildkite/go-buildkite/v3/buildkite"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	apiToken  = kingpin.Flag("token", "API token").Required().String()
	org       = kingpin.Flag("org", "Orginization slug").Required().String()
	templateUUID = kingpin.Flag("templateUUID", "Cluster UUID").Required().String()
	debug     = kingpin.Flag("debug", "Enable debugging").Bool()
)

func main() {
	kingpin.Parse()

	config, err := buildkite.NewTokenConfig(*apiToken, *debug)

	if err != nil {
		log.Fatalf("client config failed: %s", err)
	}

	client := buildkite.NewClient(config.Client())

	resp, err := client.PipelineTemplates.Delete(*org, *templateUUID)

	if err != nil {
		log.Fatalf("Deleting pipeline template %s failed: %s", *templateUUID, err)
	}

	fmt.Println(resp.StatusCode)
}
