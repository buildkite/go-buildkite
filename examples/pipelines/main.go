package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/buildkite/go-buildkite/buildkite"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	apiToken = kingpin.Flag("token", "API token").Required().String()
	org      = kingpin.Flag("org", "Organization slug").Required().String()
	debug    = kingpin.Flag("debug", "Enable debugging").Bool()

	listCommand = kingpin.Command("list", "List pipelines")

	getCommand  = kingpin.Command("get", "Get a pipeline")
	getPipeline = getCommand.Arg("pipeline", "Pipeline slug").Required().String()

	createCommand = kingpin.Command("create", "Create a pipeline")

	editCommand  = kingpin.Command("edit", "Edit a pipeline")
	editPipeline = editCommand.Arg("pipeline", "Pipeline slug").Required().String()

	deleteCommand  = kingpin.Command("delete", "Delete a pipeline")
	deletePipeline = deleteCommand.Arg("pipeline", "Pipeline slug").Required().String()

	commands = map[string]func(*buildkite.Client) (interface{}, error){
		"list": func(client *buildkite.Client) (interface{}, error) {
			pipelines, _, err := client.Pipelines.List(*org, nil)
			return pipelines, err
		},
		"get": func(client *buildkite.Client) (interface{}, error) {
			pipeline, _, err := client.Pipelines.Get(*org, *getPipeline)
			return pipeline, err
		},
		"create": func(client *buildkite.Client) (interface{}, error) {
			pipeline, _, err := client.Pipelines.Create(*org, readPipeline())
			return pipeline, err
		},
		"edit": func(client *buildkite.Client) (interface{}, error) {
			pipeline, _, err := client.Pipelines.Edit(*org, *editPipeline, readPipeline())
			return pipeline, err
		},
		"delete": func(client *buildkite.Client) (interface{}, error) {
			_, err := client.Pipelines.Delete(*org, *deletePipeline)
			return nil, err
		},
	}
)

func main() {
	command := kingpin.Parse()

	config, err := buildkite.NewTokenConfig(*apiToken, *debug)

	if err != nil {
		log.Fatalf("client config failed: %s", err)
	}

	client := buildkite.NewClient(config.Client())

	result, err := commands[command](client)
	if err != nil {
		log.Fatalf("%s failed: %s", command, err)
	}

	output(result)
}

func output(result interface{}) {
	if result == nil {
		return
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")

	err := enc.Encode(result)

	if err != nil {
		log.Fatalf("json encode failed: %s", err)
	}
}

func readPipeline() *buildkite.Pipeline {
	pipeline := &buildkite.Pipeline{}

	dec := json.NewDecoder(os.Stdin)
	err := dec.Decode(pipeline)
	if err != nil {
		log.Fatalf("json decode failed: %s", err)
	}

	return pipeline
}
