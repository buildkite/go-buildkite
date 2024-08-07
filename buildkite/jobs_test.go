package buildkite

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestJobsService_UnblockJob(t *testing.T) {
	setup(t)
	t.Cleanup(teardown)

	mux.HandleFunc("/v2/organizations/my-great-org/pipelines/sup-keith/builds/awesome-build/jobs/awesome-job-id/unblock", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		fmt.Fprint(w, `{
  "id": "awesome-job-id",
  "state": "unblocked"
}`)
	})

	job, _, err := client.Jobs.UnblockJob("my-great-org", "sup-keith", "awesome-build", "awesome-job-id", nil)
	if err != nil {
		t.Errorf("UnblockJob returned error: %v", err)
	}

	want := &Job{ID: String("awesome-job-id"), State: String("unblocked")}
	if !reflect.DeepEqual(job, want) {
		t.Errorf("UnblockJob returned %+v, want %+v", job, want)
	}
}

func TestJobsService_RetryJob(t *testing.T) {
	setup(t)
	t.Cleanup(teardown)

	mux.HandleFunc("/v2/organizations/my-great-org/pipelines/sup-keith/builds/awesome-build/jobs/awesome-job-id/retry", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		fmt.Fprint(w, `{
  "id": "awesome-job-id",
  "state": "scheduled",
  "retries_count": 1,
  "retried": true
}`)
	})

	job, _, err := client.Jobs.RetryJob("my-great-org", "sup-keith", "awesome-build", "awesome-job-id")
	if err != nil {
		t.Errorf("RetryJob returned error: %v", err)
	}

	want := &Job{ID: String("awesome-job-id"), State: String("scheduled"), Retried: bool(true), RetriesCount: int(1)}
	if !reflect.DeepEqual(job, want) {
		t.Errorf("RetryJob returned %+v, want %+v", job, want)
	}
}

func TestJobsService_GetJobLog(t *testing.T) {
	setup(t)
	t.Cleanup(teardown)

	mux.HandleFunc("/v2/organizations/my-great-org/pipelines/sup-keith/builds/awesome-build/jobs/awesome-job-id/log", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
  "url": "https://api.buildkite.com/v2/organizations/my-great-org/pipelines/sub-keith/builds/awesome-build/jobs/awesome-job-id/log",
  "content": "This is the job's log output",
	"size": 28,
	"header_times": [1563337899810051000,1563337899811015000,1563337905336878000,1563337906589603000,156333791038291900]
}`)
	})

	job, _, err := client.Jobs.GetJobLog("my-great-org", "sup-keith", "awesome-build", "awesome-job-id")
	if err != nil {
		t.Errorf("GetJobLog returned error: %v", err)
	}

	want := &JobLog{
		URL:         String("https://api.buildkite.com/v2/organizations/my-great-org/pipelines/sub-keith/builds/awesome-build/jobs/awesome-job-id/log"),
		Content:     String("This is the job's log output"),
		Size:        Int(28),
		HeaderTimes: []int64{1563337899810051000, 1563337899811015000, 1563337905336878000, 1563337906589603000, 156333791038291900},
	}
	if !reflect.DeepEqual(job, want) {
		t.Errorf("GetJobLog returned %+v, want %+v", job, want)
	}
}

func TestJobsService_GetJobEnvironmentVariables(t *testing.T) {
	setup(t)
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
	mux.HandleFunc("/v2/organizations/my-great-org/pipelines/sup-keith/builds/15/jobs/awesome-job-id/env", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		body := map[string]map[string]string{
			"env": envVars,
		}
		jsonString, err := json.Marshal(body)
		if err != nil {
			t.Errorf("Cast map into string returned error: %v", err)
		}
		fmt.Fprintf(w, `%s`, jsonString)
	})

	jobEnvVars, _, err := client.Jobs.GetJobEnvironmentVariables("my-great-org", "sup-keith", "15", "awesome-job-id")
	if err != nil {
		t.Errorf("GetJobEnvironmentVariables returned error: %v", err)
	}

	want := &JobEnvs{
		EnvironmentVariables: &envVars,
	}
	if !reflect.DeepEqual(jobEnvVars, want) {
		t.Errorf("GetJobLog returned %+v, want %+v", jobEnvVars, want)
	}
}
