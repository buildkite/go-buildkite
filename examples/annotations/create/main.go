package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/buildkite/go-buildkite/v3/buildkite"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	apiToken  = kingpin.Flag("token", "API token").Required().String()
	org       = kingpin.Flag("org", "Orginization slug").Required().String()
	buildID   = kingpin.Flag("buildID", "Build UUID").Required().String()
	debug     = kingpin.Flag("debug", "Enable debugging").Bool()
)

func main() {
	kingpin.Parse()

	config, err := buildkite.NewTokenConfig(*apiToken, *debug)

	if err != nil {
		log.Fatalf("client config failed: %s", err)
	}

	client := buildkite.NewClient(config.Client())

	annotationCreate := buildkite.AnnotationCreate{
		Style:       buildkite.String("info"),
		Context:   	 buildkite.String("default"),
		Body:        buildkite.String("An example annotation!"),
		Append: 	 buildkite.Bool(false),
	}

	annotation, _, err := client.Annotations.Create(*org, *buildID, &annotationCreate)

	if err != nil {
		log.Fatalf("Creating annotation failed: %s", err)
	}

	data, err := json.MarshalIndent(annotation, "", "\t")

	if err != nil {
		log.Fatalf("json encode failed: %s", err)
	}

	fmt.Fprintf(os.Stdout, "%s", string(data))
}
