package main

import (
	"context"
	"log"
	"os"

	"github.com/buildkite/go-buildkite/v3"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	apiToken = kingpin.Flag("token", "API token").Required().String()
	org      = kingpin.Flag("org", "Orginization slug").Required().String()
	registry = kingpin.Flag("registry", "Registry Slug").Required().String()
	filePath = kingpin.Flag("file-path", "File path").Required().String()
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

	file, err := os.Open(*filePath)
	if err != nil {
		log.Fatalf("opening file %s failed: %v", *filePath, err)
	}

	ppu, _, err := client.PackagesService.RequestPresignedUpload(context.Background(), *org, *registry)
	if err != nil {
		log.Fatalf("Creating package failed: %v", err)
	}

	url, err := ppu.Perform(context.Background(), file)
	if err != nil {
		log.Fatalf("Package upload to S3 failed: %v", err)
	}

	log.Println("Uploaded package to: " + url)

	pkg, _, err := ppu.Finalize(context.Background(), url)
	if err != nil {
		log.Fatalf("Finalizing package failed: %v", err)
	}

	log.Printf("Package uploaded: %s", pkg.Name)
}
