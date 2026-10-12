package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/alecthomas/kingpin/v2"
	"github.com/buildkite/go-buildkite/v5"
)

var (
	apiToken = kingpin.Flag("token", "API token").Required().String()
	org      = kingpin.Flag("org", "Organization slug").Required().String()
	slug     = kingpin.Flag("slug", "Test suite slug").Required().String()
	testID   = kingpin.Flag("testID", "Test ID").Required().String()
	action   = kingpin.Flag("action", "State to move the test into").Required().Enum("skip", "mute", "enable")
)

func main() {
	kingpin.Parse()

	client, err := buildkite.NewOpts(buildkite.WithTokenAuth(*apiToken))
	if err != nil {
		log.Fatalf("creating buildkite API client failed: %v", err)
	}

	ctx := context.Background()

	var test buildkite.Test
	switch *action {
	case "skip":
		test, _, err = client.Tests.Skip(ctx, *org, *slug, *testID)
	case "mute":
		test, _, err = client.Tests.Mute(ctx, *org, *slug, *testID)
	case "enable":
		test, _, err = client.Tests.Enable(ctx, *org, *slug, *testID)
	}
	if err != nil {
		log.Fatalf("%s test %s failed: %s", *action, *testID, err)
	}

	data, err := json.MarshalIndent(test, "", "\t")
	if err != nil {
		log.Fatalf("json encode failed: %s", err)
	}

	_, _ = fmt.Fprintf(os.Stdout, "%s", string(data))
}
