package main

import (
	"fmt"
	"log"

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

	updatePipeline := buildkite.UpdatePipeline{
		Slug:        *buildkite.String("my-great-repo"),
		Description: *buildkite.String("This ia a great pipeline!"),
	}

	resp, err := client.Pipelines.Update(*org, &updatePipeline)

	if err != nil {
		log.Fatalf("Updating pipeline failed: %s", err)
	}

	fmt.Println(resp.StatusCode)
}