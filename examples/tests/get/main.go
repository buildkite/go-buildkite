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
	apiToken = kingpin.Flag("token", "API token").Required().String()
	org      = kingpin.Flag("org", "Orginization slug").Required().String()
	slug     = kingpin.Flag("slug", "Test suite slug").Required().String()
	testID   = kingpin.Flag("testID", "Test ID").Required().String()
	debug    = kingpin.Flag("debug", "Enable debugging").Bool()
)

func main() {
	kingpin.Parse()

	config, err := buildkite.NewTokenConfig(*apiToken, *debug)

	if err != nil {
		log.Fatalf("client config failed: %s", err)
	}

	client := buildkite.NewClient(config.Client())

	test, _, err := client.Tests.Get(*org, *slug, *testID)

	if err != nil {
		log.Fatalf("Getting test %s failed: %s", *testID, err)
	}

	data, err := json.MarshalIndent(test, "", "\t")

	if err != nil {
		log.Fatalf("json encode failed: %s", err)
	}

	fmt.Fprintf(os.Stdout, "%s", string(data))
}
