package main

import (
	"context"
	"log"

	"github.com/buildkite/go-buildkite/v5"

	"github.com/alecthomas/kingpin/v2"
)

var (
	apiToken     = kingpin.Flag("token", "API token").Required().String()
	org          = kingpin.Flag("org", "Organization slug").Required().String()
	registrySlug = kingpin.Flag("registry", "Registry Slug").Required().String()
	tokenID      = kingpin.Flag("token-id", "Registry Token ID").Required().String()
)

func main() {
	kingpin.Parse()

	client, err := buildkite.NewOpts(buildkite.WithTokenAuth(*apiToken))
	if err != nil {
		log.Fatalf("creating buildkite API client failed: %v", err)
	}

	_, err = client.PackageRegistryTokensService.Delete(context.Background(), *org, *registrySlug, *tokenID)
	if err != nil {
		log.Fatalf("Deleting registry token %s failed: %s", *tokenID, err)
	}

	log.Printf("Registry token %s deleted!\n", *tokenID)
}
