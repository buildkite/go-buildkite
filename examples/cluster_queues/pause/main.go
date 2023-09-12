package main

import (
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

	config, err := buildkite.NewTokenConfig(*apiToken, *debug)

	if err != nil {
		log.Fatalf("client config failed: %s", err)
	}

	client := buildkite.NewClient(config.Client())

	clusterQueuePause := buildkite.ClusterQueuePause{
		Note: buildkite.String("Pausing dispatch over the weekend"),
	}

	resp, err := client.ClusterQueues.Pause(*org, *clusterID, *queueID, &clusterQueuePause)

	if err != nil {
		log.Fatalf("Pausing dispatch on cluster queue %s failed: %s", *queueID, err)
	}

	fmt.Println(resp.StatusCode)
}
