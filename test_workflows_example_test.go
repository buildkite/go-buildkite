package buildkite_test

import (
	"context"
	"log"
	"os"

	"github.com/buildkite/go-buildkite/v5"
)

func ExampleTestWorkflowsService() {
	client, err := buildkite.NewClient(buildkite.WithTokenAuth(os.Getenv("BUILDKITE_API_TOKEN")))
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	const org, suite = "acme", "web"

	opt := &buildkite.TestWorkflowListOptions{ListOptions: buildkite.ListOptions{PerPage: 100}}
	for {
		workflows, resp, err := client.TestWorkflows.List(ctx, org, suite, opt)
		if err != nil {
			log.Fatal(err)
		}
		for _, workflow := range workflows {
			log.Printf("%s: %s", workflow.ID, workflow.Monitor.Type)
		}
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	window, alarm, recover, label := 100, 0.075, 0.0, "flaky"
	workflow, _, err := client.TestWorkflows.Create(ctx, org, suite, buildkite.TestWorkflowCreate{
		Monitor: buildkite.TestWorkflowMonitor{
			Type: "transition_count", WindowSize: &window,
			AlarmThreshold: &alarm, RecoverThreshold: &recover,
		},
		Filters:        []buildkite.TestWorkflowFilter{{Key: "scm.branch", Operator: "is", Value: "main"}},
		AlarmActions:   []buildkite.TestWorkflowAction{{Type: "add_label", Label: &label}},
		RecoverActions: []buildkite.TestWorkflowAction{{Type: "remove_label", Label: &label}},
	})
	if err != nil {
		log.Fatal(err)
	}

	threshold := 0.125
	_, _, err = client.TestWorkflows.Update(ctx, org, suite, workflow.ID, buildkite.TestWorkflowUpdate{
		Monitor: buildkite.Some(buildkite.TestWorkflowMonitor{AlarmThreshold: &threshold}),
		// Unset lists remain unchanged. An initialized empty list clears one.
		RecoverActions: buildkite.Some([]buildkite.TestWorkflowAction{}),
	})
	if err != nil {
		log.Fatal(err)
	}

	workflow, _, err = client.TestWorkflows.Get(ctx, org, suite, workflow.ID)
	if err != nil {
		log.Fatal(err)
	}
	log.Print(workflow.WebURL)

	if _, err := client.TestWorkflows.Delete(ctx, org, suite, workflow.ID); err != nil {
		log.Fatal(err)
	}
}
