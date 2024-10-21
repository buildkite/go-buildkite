package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/buildkite/go-buildkite/v3"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	apiToken     = kingpin.Flag("token", "API token").Required().String()
	org          = kingpin.Flag("org", "Orginization slug").Required().String()
	registrySlug = kingpin.Flag("registry", "Registry Slug").Required().String()
	filePath     = kingpin.Flag("file-path", "File path").Required().String()
)

func main() {
	kingpin.Parse()

	client, err := buildkite.NewOpts(buildkite.WithTokenAuth(*apiToken))
	if err != nil {
		log.Fatalf("creating buildkite API client failed: %v", err)
	}

	file, err := os.Open(*filePath)
	if err != nil {
		log.Fatalf("opening file %s failed: %v", *filePath, err)
	}

	pkg, _, err := client.PackagesService.Create(context.Background(), *org, *registrySlug, buildkite.CreatePackageInput{Package: file})
	if err != nil {
		log.Fatalf("Creating package failed: %v", err)
	}

	data, err := json.MarshalIndent(pkg, "", "\t")
	if err != nil {
		log.Fatalf("json encode failed: %s", err)
	}

	log.Println(string(data))
}
