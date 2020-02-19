package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/buildkite/go-buildkite/v2/buildkite"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	apiToken     = kingpin.Flag("token", "API token").Required().String()
	org          = kingpin.Flag("org", "Orginization slug").Required().String()
	pipeline     = kingpin.Flag("pipeline", "Pipeline slug").Required().String()
	build        = kingpin.Flag("build", "Build Number slug").Required().String()
	artifactName = kingpin.Flag("artifact", "Artifact to download").String()
	debug        = kingpin.Flag("debug", "Enable debugging").Bool()
)

func main() {
	kingpin.Parse()

	config, err := buildkite.NewTokenConfig(*apiToken, *debug)

	if err != nil {
		log.Fatalf("client config failed: %s", err)
	}

	client := buildkite.NewClient(config.Client())

	artifacts, _, err := client.Artifacts.ListByBuild(*org, *pipeline, *build, nil)

	if err != nil {
		log.Fatalf("list artifacts failed: %s", err)
	}

	for _, artifact := range artifacts {
		if *artifactName == "" {
			data, err := json.MarshalIndent(artifact, "", "\t")

			if err != nil {
				log.Fatalf("json encode failed: %s", err)
			}

			fmt.Fprintf(os.Stdout, "%s\n", string(data))
		} else {
			if *artifactName == *artifact.Filename || *artifactName == *artifact.ID {
				_, err := client.Artifacts.DownloadArtifactByURL(*artifact.DownloadURL, os.Stdout)
				if err != nil {
					log.Fatalf("DownloadArtifactByURL failed: %s", err)
				}
				os.Exit(0)
			}
		}
	}
}
