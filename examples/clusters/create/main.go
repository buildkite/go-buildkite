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

	clusterCreate := buildkite.ClusterCreate{
		Name:        "Development Cluster",
		Description: buildkite.String("A cluster for development work"),
		Emoji:       buildkite.String(":toolbox:"),
		Color:       buildkite.String("#A9CCE3"),
	}

	cluster, _, err := client.Clusters.Create(*org, &clusterCreate)

	if err != nil {
		log.Fatalf("Creating cluster failed: %s", err)
	}

	data, err := json.MarshalIndent(cluster, "", "\t")

	if err != nil {
		log.Fatalf("json encode failed: %s", err)
	}

	fmt.Fprintf(os.Stdout, "%s", string(data))
}
