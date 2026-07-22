package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/alecthomas/kingpin/v2"
	"github.com/buildkite/go-buildkite/v5"
)

var (
	apiToken     = kingpin.Flag("token", "API token").Required().String()
	org          = kingpin.Flag("org", "Organization slug").Required().String()
	registrySlug = kingpin.Flag("registry", "Registry Slug").Required().String()
	description  = kingpin.Flag("description", "Token description").Required().String()
)

func main() {
	kingpin.Parse()

	client, err := buildkite.NewOpts(buildkite.WithTokenAuth(*apiToken))
	if err != nil {
		log.Fatalf("creating buildkite API client failed: %v", err)
	}

	token, _, err := client.PackageRegistryTokensService.Create(context.Background(), *org, *registrySlug, buildkite.CreatePackageRegistryTokenInput{
		Description: *description,
	})
	if err != nil {
		log.Fatalf("Creating registry token failed: %v", err)
	}

	data, err := json.MarshalIndent(token, "", "\t")
	if err != nil {
		log.Fatalf("json encode failed: %s", err)
	}

	log.Println(string(data))
}
