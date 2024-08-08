package main

import (
	"encoding/json"
	"log"

	"github.com/buildkite/go-buildkite/v3/buildkite"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	apiToken            = kingpin.Flag("token", "API token").Required().String()
	org                 = kingpin.Flag("org", "Orginization slug").Required().String()
	registryName        = kingpin.Flag("registry-name", "Registry Name").Required().String()
	registryEcosystem   = kingpin.Flag("registry-ecosystem", "Registry Ecosystem").Required().String()
	registryDescription = kingpin.Flag("registry-description", "Registry Description").String()
)

func main() {
	kingpin.Parse()

	client, err := buildkite.NewOpts(buildkite.WithTokenAuth(*apiToken))
	if err != nil {
		log.Fatalf("creating buildkite API client failed: %v", err)
	}

	reg, _, err := client.PackageRegistriesService.Create(*org, buildkite.CreatePackageRegistryInput{
		Name:        *registryName,
		Ecosystem:   *registryEcosystem,
		Description: *registryDescription,
	})
	if err != nil {
		log.Fatalf("Creating registry failed: %v", err)
	}

	data, err := json.MarshalIndent(reg, "", "\t")
	if err != nil {
		log.Fatalf("json encode failed: %s", err)
	}

	log.Println(string(data))
}
