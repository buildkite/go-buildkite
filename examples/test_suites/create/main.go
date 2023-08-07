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
	debug    = kingpin.Flag("debug", "Enable debugging").Bool()
)

func main() {
	kingpin.Parse()

	config, err := buildkite.NewTokenConfig(*apiToken, *debug)

	if err != nil {
		log.Fatalf("client config failed: %s", err)
	}

	client := buildkite.NewClient(config.Client())

	suiteCreate := buildkite.TestSuiteCreate{
		Name:			"RSpec tests",   						
		DefaultBranch:	"main",        
		TeamUuids: 		[]string{"474de468-84d6-46dc-ba23-bac1add44a60"},             
	}

	suite, _, err := client.TestSuites.Create(*org, &suiteCreate)

	if err != nil {
		log.Fatalf("Creating test suite failed: %s", err)
	}

	data, err := json.MarshalIndent(suite, "", "\t")

	if err != nil {
		log.Fatalf("json encode failed: %s", err)
	}

	fmt.Fprintf(os.Stdout, "%s", string(data))
}