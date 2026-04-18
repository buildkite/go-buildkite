package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/buildkite/go-buildkite/v4"

	"github.com/alecthomas/kingpin/v2"
)

var (
	apiToken = kingpin.Flag("token", "API token").Required().String()
	org      = kingpin.Flag("org", "Organization slug").Required().String()
	slug     = kingpin.Flag("slug", "Pipeline slug").Required().String()
	number   = kingpin.Flag("number", "Build number").Required().String()
	jobID    = kingpin.Flag("job", "Job ID").Required().String()
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
		Body:    "An example job annotation!",
		Append:  false,
	}

	annotation, _, err := client.Annotations.CreateForJob(context.Background(), *org, *slug, *number, *jobID, annotationCreate)
	if err != nil {
		log.Fatalf("Creating annotation for job %s in build %s failed: %s", *jobID, *number, err)
	}

	data, err := json.MarshalIndent(annotation, "", "\t")
	if err != nil {
		log.Fatalf("json encode failed: %s", err)
	}

	_, _ = fmt.Fprintf(os.Stdout, "%s", string(data))
}
