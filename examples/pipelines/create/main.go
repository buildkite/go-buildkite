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

	createPipeline := buildkite.CreatePipeline{
		Name:          *buildkite.String("my-great-pipeline"),
		Repository:    *buildkite.String("git@github.com:my_great_org/my_great_repo2.git"),
		Configuration: *buildkite.String("env:\n \"FOO\": \"bar\"\nsteps:\n - command: \"script/release.sh\"\n   \"name\": \"Build 📦\""),
		Tags:          []string{"great", "pipeline"},
		Description:   *buildkite.String("This ia a great pipeline!"),
		ProviderSettings: &buildkite.GitHubSettings{
			TriggerMode: buildkite.String("code"),
		},
	}

	pipeline, _, err := client.Pipelines.Create(*org, &createPipeline)

	if err != nil {
		log.Fatalf("Updating pipeline failed: %s", err)
	}

	data, err := json.MarshalIndent(pipeline, "", "\t")

	if err != nil {
		log.Fatalf("json encode failed: %s", err)
	}

	fmt.Fprintf(os.Stdout, "%s", string(data))
}
