package main

import (
	"context"
	"fmt"
	"log"

	"github.com/buildkite/go-buildkite/v3/buildkite"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	apiToken       = kingpin.Flag("token", "API token").Required().String()
	org            = kingpin.Flag("org", "Orginization slug").Required().String()
	clusterID      = kingpin.Flag("clusterID", "Cluster UUID").Required().String()
	queueID        = kingpin.Flag("queueID", "Cluster queue UUID").Required().String()
	newDescription = kingpin.Flag("description", "New description for the cluster queue").Required().String()
	debug          = kingpin.Flag("debug", "Enable debugging").Bool()
)

func main() {
	kingpin.Parse()

	client, err := buildkite.NewOpts(buildkite.WithTokenAuth(*apiToken))
	if err != nil {
		log.Fatalf("creating buildkite API client failed: %v", err)
	}

	clusterQueueUpdate := buildkite.ClusterQueueUpdate{Description: *newDescription}

	cq, _, err := client.ClusterQueues.Update(context.Background(), *org, *clusterID, *queueID, clusterQueueUpdate)
	if err != nil {
		log.Fatalf("Updating cluster queue failed: %s", err)
	}

	fmt.Printf("Updated cluster queue: %v\n", cq)
}
