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
	apiToken   = kingpin.Flag("token", "Buildkite API token with read_organizations scope").Envar("BUILDKITE_API_TOKEN").Required().String()
	org        = kingpin.Flag("org", "Organization slug or UUID").Required().String()
	from       = kingpin.Flag("from", "Inclusive start date in UTC (YYYY-MM-DD)").Required().String()
	to         = kingpin.Flag("to", "Inclusive end date in UTC (YYYY-MM-DD)").Required().String()
	interval   = kingpin.Flag("interval", "Reporting interval: day or month").Default("day").String()
	percentile = kingpin.Flag("percentile", "Discrete percentile of daily sums (1-100)").Default("95").Int()
	agentType  = kingpin.Flag("agent-type", "Agent type: self_hosted or hosted").Default("self_hosted").String()
)

func main() {
	kingpin.Parse()

	client, err := buildkite.NewClient(buildkite.WithTokenAuth(*apiToken))
	if err != nil {
		log.Fatal(err)
	}
	usage, _, err := client.Organizations.GetConcurrentAgentsUsage(context.Background(), *org, &buildkite.ConcurrentAgentsUsageOptions{
		From: *from, To: *to, Interval: *interval, Percentile: *percentile, AgentType: *agentType,
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
