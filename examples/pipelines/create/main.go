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

	pipelineC := buildkite.CreatePipeline{Name: *buildkite.String("go-bk-pipe"),
		Repository: *buildkite.String("my-great-repo"),
		Steps: []buildkite.Step{
			{
				Type:    buildkite.String("script"),
				Name:    buildkite.String("Build :package"),
				Command: buildkite.String("script/release.sh"),
				Plugins: buildkite.Plugins{
					"my-org/docker#v3.3.0": {
						"image":   "node",
						"workdir": "/app",
					},
				},
			},
		},
		DefaultBranch: *buildkite.String("main"),
	}

	pipeline, _, err := client.Pipelines.Create(*org, &pipelineC)

	if err != nil {
		log.Fatalf("Creating test suite failed: %s", err)
	}

	data, err := json.MarshalIndent(pipeline, "", "\t")

	if err != nil {
		log.Fatalf("json encode failed: %s", err)
	}

	fmt.Fprintf(os.Stdout, "%s", string(data))
}