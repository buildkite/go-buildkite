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

	clusterCreate := buildkite.ClusterCreate{
		Name:        "Development Cluster",
		Description: "A cluster for development work",
		Emoji:       ":toolbox:",
		Color:       "#A9CCE3",
	}

	cluster, _, err := client.Clusters.Create(context.Background(), *org, clusterCreate)
	if err != nil {
		log.Fatalf("Creating cluster failed: %s", err)
	}

	data, err := json.MarshalIndent(cluster, "", "\t")
	if err != nil {
		log.Fatalf("json encode failed: %s", err)
	}

	_, _ = fmt.Fprintf(os.Stdout, "%s", string(data))
}
