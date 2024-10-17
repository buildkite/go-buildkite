package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/buildkite/go-buildkite/v3"

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

	clusterTokenCreate := buildkite.ClusterTokenCreateUpdate{Description: "Dev squad agent fleet token"}

	token, _, err := client.ClusterTokens.Create(context.Background(), *org, *clusterID, clusterTokenCreate)
	if err != nil {
		log.Fatalf("Creating cluster token failed: %s", err)
	}

	data, err := json.MarshalIndent(token, "", "\t")
	if err != nil {
		log.Fatalf("json encode failed: %s", err)
	}

	fmt.Fprintf(os.Stdout, "%s", string(data))
}
