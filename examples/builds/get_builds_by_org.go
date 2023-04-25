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

	builds, _, err := client.Builds.ListByOrg(org, nil)

	if err != nil {
		log.Fatalf("Get build failed: %s", err)
	}

	data, err := json.MarshalIndent(builds, "", "\t")

	if err != nil {
		log.Fatalf("JSON encode failed: %s", err)
	}

	fmt.Fprintf(os.Stdout, "%s", string(data))
}
