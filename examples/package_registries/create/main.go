package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/buildkite/go-buildkite/v4"

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

	reg, _, err := client.PackageRegistriesService.Create(context.Background(), *org, buildkite.CreatePackageRegistryInput{
		Name:        *registryName,
		Ecosystem:   *registryEcosystem,
		Description: *registryDescription,

		OIDCPolicy: buildkite.PackageRegistryOIDCPolicy{
			buildkite.OIDCPolicyStatement{
				Issuer: "https://agent.buildkite.com",
				Scopes: []string{"read_packages"},
				Claims: map[string]buildkite.ClaimRule{
					"pipeline_slug": {
						Equals: "my-pipeline",
					},
				},
			},
		},
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
