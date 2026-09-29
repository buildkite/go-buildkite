package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/alecthomas/kingpin/v2"
	"github.com/buildkite/go-buildkite/v5"
)

var (
	apiToken     = kingpin.Flag("token", "API token").Required().String()
	org          = kingpin.Flag("org", "Organization slug").Required().String()
	connectionID = kingpin.Flag("connection", "Repository connection UUID").Required().String()
)

func main() {
	kingpin.Parse()

	client, err := buildkite.NewOpts(buildkite.WithTokenAuth(*apiToken))
	if err != nil {
		log.Fatalf("creating Buildkite API client failed: %v", err)
	}

	connection, _, err := client.RepositoryConnections.Get(context.Background(), *org, *connectionID)
	if err != nil {
		log.Fatalf("getting repository connection %s failed: %v", *connectionID, err)
	}

	data, err := json.MarshalIndent(connection, "", "\t")
	if err != nil {
		log.Fatalf("encoding repository connection failed: %v", err)
	}

	_, _ = fmt.Fprintln(os.Stdout, string(data))
}
