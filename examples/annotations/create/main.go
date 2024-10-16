package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/buildkite/go-buildkite/v3/buildkite"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	apiToken = kingpin.Flag("token", "API token").Required().String()
	org      = kingpin.Flag("org", "Orginization slug").Required().String()
	slug     = kingpin.Flag("slug", "Pipeline slug").Required().String()
	number   = kingpin.Flag("number", "Build number").Required().String()
	debug    = kingpin.Flag("debug", "Enable debugging").Bool()
)

func main() {
	kingpin.Parse()

	client, err := buildkite.NewOpts(buildkite.WithTokenAuth(*apiToken))
	if err != nil {
		log.Fatalf("creating buildkite API client failed: %v", err)
	}

	annotationCreate := buildkite.AnnotationCreate{
		Style:   "info",
		Context: "default",
		Body:    "An example annotation!",
		Append:  false,
	}

	annotation, _, err := client.Annotations.Create(context.Background(), *org, *slug, *number, annotationCreate)
	if err != nil {
		log.Fatalf("Listing annotations for build %s in pipeline %s failed: %s", *number, *slug, err)
	}

	data, err := json.MarshalIndent(annotation, "", "\t")
	if err != nil {
		log.Fatalf("json encode failed: %s", err)
	}

	fmt.Fprintf(os.Stdout, "%s", string(data))
}
