package main

import (
	"context"
	"fmt"
	"log"

	"github.com/buildkite/go-buildkite/v3"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	apiToken  = kingpin.Flag("token", "API token").Required().String()
	org       = kingpin.Flag("org", "Orginization slug").Required().String()
	clusterID = kingpin.Flag("clusterID", "Cluster UUID").Required().String()
	tokenID   = kingpin.Flag("tokenID", "Cluster token UUID").Required().String()
	debug     = kingpin.Flag("debug", "Enable debugging").Bool()
)

func main() {
	kingpin.Parse()

	client, err := buildkite.NewOpts(buildkite.WithTokenAuth(*apiToken))
	if err != nil {
		log.Fatalf("creating buildkite API client failed: %v", err)
	}

	resp, err := client.ClusterTokens.Delete(context.Background(), *org, *clusterID, *tokenID)

	if err != nil {
		log.Fatalf("Deleting cluster token %s failed: %s", *tokenID, err)
	}

	fmt.Println(resp.StatusCode)
}
