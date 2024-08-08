package main

import (
	"context"
	"fmt"
	"log"

	"github.com/buildkite/go-buildkite/v3/buildkite"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	apiToken  = kingpin.Flag("token", "API token").Required().String()
	org       = kingpin.Flag("org", "Orginization slug").Required().String()
	clusterID = kingpin.Flag("clusterID", "Cluster UUID").Required().String()
	queueID   = kingpin.Flag("queueID", "Cluster queue UUID").Required().String()
	debug     = kingpin.Flag("debug", "Enable debugging").Bool()
)

func main() {
	kingpin.Parse()

	client, err := buildkite.NewOpts(buildkite.WithTokenAuth(*apiToken))
	if err != nil {
		log.Fatalf("creating buildkite API client failed: %v", err)
	}

	clusterQueuePause := buildkite.ClusterQueuePause{
		Note: buildkite.String("Pausing dispatch over the weekend"),
	}

	resp, err := client.ClusterQueues.Pause(context.Background(), *org, *clusterID, *queueID, &clusterQueuePause)

	if err != nil {
		log.Fatalf("Pausing dispatch on cluster queue %s failed: %s", *queueID, err)
	}

	fmt.Println(resp.StatusCode)
}
