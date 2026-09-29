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
	apiToken = kingpin.Flag("token", "API token").Required().String()
	org      = kingpin.Flag("org", "Organization slug").Required().String()
)

func main() {
	kingpin.Parse()

	client, err := buildkite.NewOpts(buildkite.WithTokenAuth(*apiToken))
	if err != nil {
		log.Fatalf("creating Buildkite API client failed: %v", err)
	}

	connections, _, err := client.RepositoryConnections.List(context.Background(), *org)
	if err != nil {
		log.Fatalf("listing repository connections failed: %v", err)
	}

	data, err := json.MarshalIndent(connections, "", "\t")
	if err != nil {
		log.Fatalf("encoding repository connections failed: %v", err)
	}

	_, _ = fmt.Fprintln(os.Stdout, string(data))
}
