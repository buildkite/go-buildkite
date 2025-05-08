package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/buildkite/go-buildkite/v4"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	apiToken = kingpin.Flag("token", "API token").Required().String()
	org      = kingpin.Flag("org", "Orginization slug").Required().String()
	registry = kingpin.Flag("registry", "Registry slug").Required().String()
)

func main() {
	kingpin.Parse()

	client, err := buildkite.NewOpts(
		buildkite.WithTokenAuth(*apiToken),
		buildkite.WithBaseURL("http://api.buildkite.localhost"),
	)
	if err != nil {
		log.Fatalf("creating buildkite API client failed: %v", err)
	}

	opts := &buildkite.RegistryPackagesOptions{PerPage: "2"}
	packages, _, err := client.PackageRegistriesService.ListPackages(context.Background(), *org, *registry, opts)
	if err != nil {
		log.Fatalf("Getting registry packages failed: %v", err)
	}

	data, err := json.MarshalIndent(packages, "", "\t")
	if err != nil {
		log.Fatalf("json encode failed: %s", err)
	}

	log.Println(string(data))

	nextOpts, err := packages.Links.Next.ToOptions()
	if err != nil {
		log.Fatalf("options encoding failed: %s", err)
	}

	nextPackages, _, err := client.PackageRegistriesService.ListPackages(context.Background(), *org, *registry, nextOpts)
	if err != nil {
		log.Fatalf("Getting registry packages failed: %v", err)
	}

	data, err = json.MarshalIndent(nextPackages, "", "\t")
	if err != nil {
		log.Fatalf("json encode failed: %s", err)
	}

	log.Println(string(data))
}
