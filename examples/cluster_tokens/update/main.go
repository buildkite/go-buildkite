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
	tokenID   = kingpin.Flag("tokenID", "Cluster token UUID").Required().String()
	debug     = kingpin.Flag("debug", "Enable debugging").Bool()
)

func main() {
	kingpin.Parse()

	client, err := buildkite.NewOpts(buildkite.WithTokenAuth(*apiToken))
	if err != nil {
		log.Fatalf("creating buildkite API client failed: %v", err)
	}

	clusterTokenUpdate := buildkite.ClusterTokenCreateUpdate{
		Description: buildkite.String("Dev squad agent token"),
	}

	resp, err := client.ClusterTokens.Update(*org, *clusterID, *tokenID, &clusterTokenUpdate)

	if err != nil {
		log.Fatalf("Updating cluster token failed: %s", err)
	}

	fmt.Println(resp.StatusCode)
}
