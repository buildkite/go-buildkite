package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/buildkite/go-buildkite/v3/buildkite"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	apiToken     = kingpin.Flag("token", "API token").Required().String()
	org          = kingpin.Flag("org", "Orginization slug").Required().String()
	registrySlug = kingpin.Flag("registry", "Registry Slug").Required().String()
)

func main() {
	kingpin.Parse()

	client, err := buildkite.NewOpts(buildkite.WithTokenAuth(*apiToken))
	if err != nil {
		log.Fatalf("creating buildkite API client failed: %v", err)
	}

	reg, _, err := client.PackageRegistriesService.Get(context.Background(), *org, *registrySlug)
	if err != nil {
		log.Fatalf("Getting registry %s failed: %s", *registrySlug, err)
	}

	data, err := json.MarshalIndent(reg, "", "\t")
	if err != nil {
		log.Fatalf("json encode failed: %s", err)
	}

	log.Println(string(data))
}
