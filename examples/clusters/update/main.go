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
	debug     = kingpin.Flag("debug", "Enable debugging").Bool()
)

func main() {
	kingpin.Parse()

	client, err := buildkite.NewOpts(buildkite.WithTokenAuth(*apiToken))
	if err != nil {
		log.Fatalf("creating buildkite API client failed: %v", err)
	}

	clusterUpdate := buildkite.ClusterUpdate{
		Description: buildkite.String("Development cluster"),
	}

	resp, err := client.Clusters.Update(context.Background(), *org, *clusterID, &clusterUpdate)

	if err != nil {
		log.Fatalf("Updating cluster %s failed: %s", *clusterID, err)
	}

	fmt.Println(resp.StatusCode)
}
