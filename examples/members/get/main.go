package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/buildkite/go-buildkite/v4"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	apiToken   = kingpin.Flag("token", "API token").Required().String()
	org        = kingpin.Flag("org", "Orginization slug").Required().String()
	memberUUID = kingpin.Flag("memberUUID", "Member UUID").Required().String()
)

func main() {
	kingpin.Parse()

	client, err := buildkite.NewOpts(buildkite.WithTokenAuth(*apiToken))
	if err != nil {
		log.Fatalf("creating buildkite API client failed: %v", err)
	}

	member, _, err := client.Members.Get(context.Background(), *org, *memberUUID)
	if err != nil {
		log.Fatalf("Getting member %s failed: %s", *memberUUID, err)
	}

	data, err := json.MarshalIndent(member, "", "\t")
	if err != nil {
		log.Fatalf("json encode failed: %s", err)
	}

	_, _ = fmt.Fprintf(os.Stdout, "%s", string(data))
}
