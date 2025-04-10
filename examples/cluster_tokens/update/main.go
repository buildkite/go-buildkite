package main

import (
	"context"
	"fmt"
	"log"

	"github.com/buildkite/go-buildkite/v4"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	apiToken       = kingpin.Flag("token", "API token").Required().String()
	org            = kingpin.Flag("org", "Orginization slug").Required().String()
	clusterID      = kingpin.Flag("clusterID", "Cluster UUID").Required().String()
	tokenID        = kingpin.Flag("tokenID", "Cluster token UUID").Required().String()
	newDescription = kingpin.Flag("description", "New description for the cluster token").Required().String()
)

func main() {
	kingpin.Parse()

	client, err := buildkite.NewOpts(buildkite.WithTokenAuth(*apiToken))
	if err != nil {
		log.Fatalf("creating buildkite API client failed: %v", err)
	}

	clusterTokenUpdate := buildkite.ClusterTokenCreateUpdate{Description: *newDescription}

	token, _, err := client.ClusterTokens.Update(context.Background(), *org, *clusterID, *tokenID, clusterTokenUpdate)
	if err != nil {
		log.Fatalf("Updating cluster token failed: %s", err)
	}

	fmt.Printf("Updated cluster token: %s\n", token.Description)
}
