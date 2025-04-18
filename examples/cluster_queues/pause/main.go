package main

import (
	"context"
	"fmt"
	"log"

	"github.com/buildkite/go-buildkite/v4"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	apiToken  = kingpin.Flag("token", "API token").Required().String()
	org       = kingpin.Flag("org", "Orginization slug").Required().String()
	clusterID = kingpin.Flag("clusterID", "Cluster UUID").Required().String()
	queueID   = kingpin.Flag("queueID", "Cluster queue UUID").Required().String()
	pauseNote = kingpin.Flag("note", "Note to add to the pause").String()
)

func main() {
	kingpin.Parse()

	client, err := buildkite.NewOpts(buildkite.WithTokenAuth(*apiToken))
	if err != nil {
		log.Fatalf("creating buildkite API client failed: %v", err)
	}

	clusterQueuePause := buildkite.ClusterQueuePause{Note: *pauseNote}
	cq, _, err := client.ClusterQueues.Pause(context.Background(), *org, *clusterID, *queueID, clusterQueuePause)
	if err != nil {
		log.Fatalf("Pausing dispatch on cluster queue %s failed: %s", *queueID, err)
	}

	fmt.Printf("Paused dispatch on cluster queue: %s with reason %q\n", cq.ID, cq.DispatchPausedNote)
}
