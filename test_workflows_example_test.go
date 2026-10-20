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
			detail, _, err := client.TestWorkflows.Get(ctx, org, suite, workflow.ID)
			if err != nil {
				log.Fatal(err)
			}
			log.Print(detail.WebURL)
		}
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
}
