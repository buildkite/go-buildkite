package buildkite

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestTestsService_Get(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/analytics/organizations/my-great-org/suites/suite-example/tests/b3abe2e9-35c5-4905-85e1-8c9f2da3240f", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = fmt.Fprint(w,
			`
			{
				"id": "b3abe2e9-35c5-4905-85e1-8c9f2da3240f",
				"url": "https://api.buildkite.com/v2/analytics/organizations/my-great-org/suite-example/tests/b3abe2e9-35c5-4905-85e1-8c9f2da3240f",
				"web_url": "https://buildkite.com/organizations/my-great-org/analytics/suite-example/tests/b3abe2e9-35c5-4905-85e1-8c9f2da3240f",
				"name": "TestExample1_Create",
				"scope": "User#email",
				"location": "./resources/test_example_test.go:123",
				"file_name": "./resources/test_example_test.go"
			}`)
	})

	got, _, err := client.Tests.Get(context.Background(), "my-great-org", "suite-example", "b3abe2e9-35c5-4905-85e1-8c9f2da3240f")
	if err != nil {
		t.Errorf("TestSuites.Get returned error: %v", err)
	}

	want := Test{
		ID:       "b3abe2e9-35c5-4905-85e1-8c9f2da3240f",
		URL:      "https://api.buildkite.com/v2/analytics/organizations/my-great-org/suite-example/tests/b3abe2e9-35c5-4905-85e1-8c9f2da3240f",
		WebURL:   "https://buildkite.com/organizations/my-great-org/analytics/suite-example/tests/b3abe2e9-35c5-4905-85e1-8c9f2da3240f",
		Name:     "TestExample1_Create",
		Scope:    "User#email",
		Location: "./resources/test_example_test.go:123",
		FileName: "./resources/test_example_test.go",
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("TestsService.Get diff: (-got +want)\n%s", diff)
	}
}
