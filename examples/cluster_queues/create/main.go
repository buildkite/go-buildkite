package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/buildkite/go-buildkite/v3/buildkite"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	apiToken  = kingpin.Flag("token", "API token").Required().String()
	org       = kingpin.Flag("org", "Orginization slug").Required().String()
	clusterID = kingpin.Flag("clusterID", "Cluster UUID").Required().String()
	debug     = kingpin.Flag("debug", "Enable debugging").Bool()
)

func main() {
	kingpin.Parse()

	client, err := buildkite.NewOpts(buildkite.WithTokenAuth(*apiToken))
	if err != nil {
		log.Fatalf("creating buildkite API client failed: %v", err)
	}

	clusterQueueCreate := buildkite.ClusterQueueCreate{
		Key:         buildkite.String("dev1"),
		Description: buildkite.String("Development 1 Cluster queue"),
	}

	queue, _, err := client.ClusterQueues.Create(context.Background(), *org, *clusterID, &clusterQueueCreate)

	if err != nil {
		log.Fatalf("Creating cluster queue failed: %s", err)
	}

	data, err := json.MarshalIndent(queue, "", "\t")

	if err != nil {
		log.Fatalf("json encode failed: %s", err)
	}

	fmt.Fprintf(os.Stdout, "%s", string(data))
}
