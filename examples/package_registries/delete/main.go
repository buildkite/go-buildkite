package main

import (
	"context"
	"log"

	"github.com/buildkite/go-buildkite/v4"

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

	_, err = client.PackageRegistriesService.Delete(context.Background(), *org, *registrySlug)
	if err != nil {
		log.Fatalf("Getting registry %s failed: %s", *registrySlug, err)
	}

	log.Printf("Registry %s deleted!\n", *registrySlug)
}
