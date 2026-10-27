package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/alecthomas/kingpin/v2"
	"github.com/buildkite/go-buildkite/v5"
)

var (
	apiToken = kingpin.Flag("token", "Buildkite API token with read_organizations scope").Envar("BUILDKITE_API_TOKEN").Required().String()
	org      = kingpin.Flag("org", "Organization slug or UUID").Required().String()
	from     = kingpin.Flag("from", "Inclusive start date in UTC (YYYY-MM-DD)").Required().String()
	to       = kingpin.Flag("to", "Inclusive end date in UTC (YYYY-MM-DD)").Required().String()
	interval = kingpin.Flag("interval", "Reporting interval: day or month").Default("day").String()
)

func main() {
	kingpin.Parse()

	client, err := buildkite.NewClient(buildkite.WithTokenAuth(*apiToken))
	if err != nil {
		log.Fatal(err)
	}
	usage, _, err := client.Organizations.GetHostedAgentsUsage(context.Background(), *org, &buildkite.HostedAgentsUsageOptions{
		From: *from, To: *to, Interval: *interval,
	})
	if err != nil {
		log.Fatal(err)
	}
	data, err := json.MarshalIndent(usage, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(data))
}
