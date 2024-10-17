package buildkite

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestJobsService_UnblockJob(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-great-org/pipelines/sup-keith/builds/awesome-build/jobs/awesome-job-id/unblock", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		fmt.Fprint(w, `{
  "id": "awesome-job-id",
  "state": "unblocked"
}`)
	})

	job, _, err := client.Jobs.UnblockJob(context.Background(), "my-great-org", "sup-keith", "awesome-build", "awesome-job-id", nil)
	if err != nil {
		t.Errorf("UnblockJob returned error: %v", err)
	}

	want := Job{ID: "awesome-job-id", State: "unblocked"}
	if diff := cmp.Diff(job, want); diff != "" {
		t.Errorf("UnblockJob diff: (-got +want)\n%s", diff)
	}
}

func TestJobsService_RetryJob(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-great-org/pipelines/sup-keith/builds/awesome-build/jobs/awesome-job-id/retry", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		fmt.Fprint(w, `{
  "id": "awesome-job-id",
  "state": "scheduled",
  "retries_count": 1,
  "retried": true
}`)
	})

	job, _, err := client.Jobs.RetryJob(context.Background(), "my-great-org", "sup-keith", "awesome-build", "awesome-job-id")
	if err != nil {
		t.Errorf("RetryJob returned error: %v", err)
	}

	want := Job{ID: "awesome-job-id", State: "scheduled", Retried: true, RetriesCount: 1}
	if diff := cmp.Diff(job, want); diff != "" {
		t.Errorf("RetryJob diff: (-got +want)\n%s", diff)
	}
}

func TestJobsService_GetJobLog(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-great-org/pipelines/sup-keith/builds/awesome-build/jobs/awesome-job-id/log", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
  "url": "https://api.buildkite.com/v2/organizations/my-great-org/pipelines/sub-keith/builds/awesome-build/jobs/awesome-job-id/log",
  "content": "This is the job's log output",
	"size": 28,
	"header_times": [1563337899810051000,1563337899811015000,1563337905336878000,1563337906589603000,156333791038291900]
}`)
	})

	job, _, err := client.Jobs.GetJobLog(context.Background(), "my-great-org", "sup-keith", "awesome-build", "awesome-job-id")
	if err != nil {
		t.Errorf("GetJobLog returned error: %v", err)
	}

	want := JobLog{
		URL:         "https://api.buildkite.com/v2/organizations/my-great-org/pipelines/sub-keith/builds/awesome-build/jobs/awesome-job-id/log",
		Content:     "This is the job's log output",
		Size:        28,
		HeaderTimes: []int64{1563337899810051000, 1563337899811015000, 1563337905336878000, 1563337906589603000, 156333791038291900},
	}
	if diff := cmp.Diff(job, want); diff != "" {
		t.Errorf("GetJobLog diff: (-got +want)\n%s", diff)
	}
}

func TestJobsService_GetJobEnvironmentVariables(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	envVars := map[string]string{
		"CI":                              "true",
		"BUILDKITE":                       "true",
		"BUILDKITE_TAG":                   "",
		"BUILDKITE_REPO":                  "git@github.com:my-great-org/my-repo.git",
		"BUILDKITE_BRANCH":                "master",
		"BUILDKITE_COMMIT":                "a65572555600c07c7ee79a2bd909220e1ca5485b",
		"BUILDKITE_JOB_ID":                "bde076a8-bc2c-4fda-9652-10220a56d638",
		"BUILDKITE_COMMAND":               "buildkite-agent pipeline upload",
		"BUILDKITE_MESSAGE":               ":llama:",
		"BUILDKITE_BUILD_ID":              "c4e312cb-e734-4f0a-a5bd-1cac2535c57e",
		"BUILDKITE_BUILD_URL":             "https://buildkite.com/my-great-org/my-pipeline/builds/15",
		"BUILDKITE_AGENT_NAME":            "ci-1",
		"BUILDKITE_BUILD_NUMBER":          "15",
		"BUILDKITE_ORGANIZATION_SLUG":     "my-great-org",
		"BUILDKITE_PIPELINE_SLUG":         "sup-keith",
		"BUILDKITE_PULL_REQUEST":          "false",
		"BUILDKITE_BUILD_CREATOR":         "Keith Pitt",
		"BUILDKITE_REPO_SSH_HOST":         "github.com",
		"BUILDKITE_ARTIFACT_PATHS":        "",
		"BUILDKITE_PIPELINE_PROVIDER":     "github",
		"BUILDKITE_BUILD_CREATOR_EMAIL":   "keith@buildkite.com",
		"BUILDKITE_AGENT_META_DATA_LOCAL": "true",
	}

	server.HandleFunc("/v2/organizations/my-great-org/pipelines/sup-keith/builds/15/jobs/awesome-job-id/env", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		body := map[string]map[string]string{"env": envVars}
		bytes, err := json.Marshal(body)
		if err != nil {
			t.Errorf("json.Marshal(body) returned error: %v", err)
		}

		fmt.Fprint(w, string(bytes))
	})

	jobEnvVars, _, err := client.Jobs.GetJobEnvironmentVariables(context.Background(), "my-great-org", "sup-keith", "15", "awesome-job-id")
	if err != nil {
		t.Errorf("GetJobEnvironmentVariables returned error: %v", err)
	}

	want := JobEnvs{EnvironmentVariables: envVars}
	if diff := cmp.Diff(jobEnvVars, want); diff != "" {
		t.Errorf("GetJobLog diff: (-got +want)\n%s", diff)
	}
}
