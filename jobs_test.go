package buildkite

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestJobsService_ListByBuild(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-great-org/pipelines/sup-keith/builds/123/jobs", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = fmt.Fprint(w, `{
  "items": [
    {
      "id": "job-1",
      "graphql_id": "Sm9iLS0tam9iLTE=",
      "type": "script",
      "name": ":package: Build",
      "step_key": "build",
      "state": "passed",
      "build_url": "https://api.buildkite.com/v2/organizations/my-great-org/pipelines/sup-keith/builds/123",
      "web_url": "https://buildkite.com/my-great-org/sup-keith/builds/123#job-1",
      "log_url": "https://api.buildkite.com/v2/organizations/my-great-org/pipelines/sup-keith/builds/123/jobs/job-1/log",
      "raw_log_url": "https://api.buildkite.com/v2/organizations/my-great-org/pipelines/sup-keith/builds/123/jobs/job-1/log.txt",
      "command": "scripts/build.sh",
      "soft_failed": false,
      "exit_status": 0,
      "signal": 15,
      "signal_reason": "terminated",
      "artifact_paths": "logs/**/*",
      "agent_query_rules": ["queue=default"],
      "retried": false,
      "retries_count": 0,
      "priority": { "number": 0 },
      "cluster_id": "cluster-1",
      "cluster_url": "https://api.buildkite.com/v2/organizations/my-great-org/clusters/cluster-1",
      "cluster_queue_id": "queue-1",
      "cluster_queue_url": "https://api.buildkite.com/v2/organizations/my-great-org/clusters/cluster-1/queues/queue-1"
    }
  ],
  "links": {
    "self": "https://api.buildkite.com/v2/organizations/my-great-org/pipelines/sup-keith/builds/123/jobs?per_page=30",
    "next": "https://api.buildkite.com/v2/organizations/my-great-org/pipelines/sup-keith/builds/123/jobs?after=cursor-1&per_page=30"
  }
}`)
	})

	jobs, _, err := client.Jobs.ListByBuild(context.Background(), "my-great-org", "sup-keith", "123", nil)
	if err != nil {
		t.Errorf("ListByBuild returned error: %v", err)
	}

	exitStatus := 0
	signal := 15
	want := JobsList{
		Items: []Job{{
			ID:              "job-1",
			GraphQLID:       "Sm9iLS0tam9iLTE=",
			Type:            "script",
			Name:            ":package: Build",
			StepKey:         "build",
			State:           "passed",
			BuildURL:        "https://api.buildkite.com/v2/organizations/my-great-org/pipelines/sup-keith/builds/123",
			WebURL:          "https://buildkite.com/my-great-org/sup-keith/builds/123#job-1",
			LogURL:          "https://api.buildkite.com/v2/organizations/my-great-org/pipelines/sup-keith/builds/123/jobs/job-1/log",
			RawLogsURL:      "https://api.buildkite.com/v2/organizations/my-great-org/pipelines/sup-keith/builds/123/jobs/job-1/log.txt",
			Command:         "scripts/build.sh",
			SoftFailed:      false,
			ExitStatus:      &exitStatus,
			Signal:          &signal,
			SignalReason:    "terminated",
			ArtifactPaths:   "logs/**/*",
			AgentQueryRules: []string{"queue=default"},
			Retried:         false,
			RetriesCount:    0,
			Priority:        &JobPriority{Number: 0},
			ClusterID:       "cluster-1",
			ClusterURL:      "https://api.buildkite.com/v2/organizations/my-great-org/clusters/cluster-1",
			ClusterQueueID:  "queue-1",
			ClusterQueueURL: "https://api.buildkite.com/v2/organizations/my-great-org/clusters/cluster-1/queues/queue-1",
		}},
		Links: JobsListLinks{
			Self: "https://api.buildkite.com/v2/organizations/my-great-org/pipelines/sup-keith/builds/123/jobs?per_page=30",
			Next: "https://api.buildkite.com/v2/organizations/my-great-org/pipelines/sup-keith/builds/123/jobs?after=cursor-1&per_page=30",
		},
	}
	if diff := cmp.Diff(jobs, want); diff != "" {
		t.Errorf("ListByBuild diff: (-got +want)\n%s", diff)
	}
}

func TestJobsService_ListByBuild_WithOptions(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-great-org/pipelines/sup-keith/builds/123/jobs", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValuesList(t, r, valuesList{
			{"state[]", "passed"},
			{"state[]", "failed"},
			{"include_retried_jobs", "false"},
			{"per_page", "100"},
			{"after", "cursor-1"},
		})
		_, _ = fmt.Fprint(w, `{"items":[],"links":{}}`)
	})

	includeRetriedJobs := false
	opt := &JobsListOptions{
		State:              []string{"passed", "failed"},
		IncludeRetriedJobs: &includeRetriedJobs,
		PerPage:            100,
		After:              "cursor-1",
	}
	_, _, err := client.Jobs.ListByBuild(context.Background(), "my-great-org", "sup-keith", "123", opt)
	if err != nil {
		t.Errorf("ListByBuild returned error: %v", err)
	}
}

func TestJobsListLink_ToOptions(t *testing.T) {
	t.Parallel()

	link := JobsListLink("https://api.buildkite.com/v2/organizations/my-great-org/pipelines/sup-keith/builds/123/jobs?state[]=passed&state[]=failed&include_retried_jobs=false&after=cursor-1&per_page=100")
	opt, err := link.ToOptions()
	if err != nil {
		t.Fatalf("ToOptions returned error: %v", err)
	}

	includeRetriedJobs := false
	want := &JobsListOptions{
		State:              []string{"passed", "failed"},
		IncludeRetriedJobs: &includeRetriedJobs,
		PerPage:            100,
		After:              "cursor-1",
	}
	if diff := cmp.Diff(opt, want); diff != "" {
		t.Errorf("ToOptions diff: (-got +want)\n%s", diff)
	}
}

func TestJobsService_GetJob(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-great-org/pipelines/sup-keith/builds/123/jobs/job-1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = fmt.Fprint(w, `{
  "id": "job-1",
  "graphql_id": "Sm9iLS0tam9iLTE=",
  "type": "script",
  "name": ":package: Build",
  "step_key": "build",
  "state": "passed",
  "build_url": "https://api.buildkite.com/v2/organizations/my-great-org/pipelines/sup-keith/builds/123",
  "web_url": "https://buildkite.com/my-great-org/sup-keith/builds/123#job-1",
  "log_url": "https://api.buildkite.com/v2/organizations/my-great-org/pipelines/sup-keith/builds/123/jobs/job-1/log",
  "raw_log_url": "https://api.buildkite.com/v2/organizations/my-great-org/pipelines/sup-keith/builds/123/jobs/job-1/log.txt",
  "command": "scripts/build.sh",
  "soft_failed": false,
  "exit_status": 0,
  "signal": 15,
  "signal_reason": "terminated",
  "expired_at": "2026-06-03T04:15:41.618Z",
  "matrix": {
    "os": "linux",
    "go": "1.25"
  },
  "retried_by": {
    "id": "user-1",
    "name": "Keith Pitt",
    "email": "keith@buildkite.com"
  }
}`)
	})

	job, _, err := client.Jobs.GetJob(context.Background(), "my-great-org", "sup-keith", "123", "job-1")
	if err != nil {
		t.Errorf("GetJob returned error: %v", err)
	}

	exitStatus := 0
	signal := 15
	want := Job{
		ID:           "job-1",
		GraphQLID:    "Sm9iLS0tam9iLTE=",
		Type:         "script",
		Name:         ":package: Build",
		StepKey:      "build",
		State:        "passed",
		BuildURL:     "https://api.buildkite.com/v2/organizations/my-great-org/pipelines/sup-keith/builds/123",
		WebURL:       "https://buildkite.com/my-great-org/sup-keith/builds/123#job-1",
		LogURL:       "https://api.buildkite.com/v2/organizations/my-great-org/pipelines/sup-keith/builds/123/jobs/job-1/log",
		RawLogsURL:   "https://api.buildkite.com/v2/organizations/my-great-org/pipelines/sup-keith/builds/123/jobs/job-1/log.txt",
		Command:      "scripts/build.sh",
		SoftFailed:   false,
		ExitStatus:   &exitStatus,
		Signal:       &signal,
		SignalReason: "terminated",
		ExpiredAt:    NewTimestamp(must(time.Parse(BuildKiteDateFormat, "2026-06-03T04:15:41.618Z"))),
		Matrix:       map[string]any{"os": "linux", "go": "1.25"},
		RetriedBy:    &User{ID: "user-1", Name: "Keith Pitt", Email: "keith@buildkite.com"},
	}
	if diff := cmp.Diff(job, want); diff != "" {
		t.Errorf("GetJob diff: (-got +want)\n%s", diff)
	}
}

func TestJobsService_GetJobByOrg(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-great-org/jobs/job-1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = fmt.Fprint(w, `{
  "id": "job-1",
  "type": "script",
  "name": ":package: Build",
  "state": "passed",
  "build_url": "https://api.buildkite.com/v2/organizations/my-great-org/pipelines/sup-keith/builds/123",
  "web_url": "https://buildkite.com/my-great-org/sup-keith/builds/123#job-1",
  "log_url": "https://api.buildkite.com/v2/organizations/my-great-org/pipelines/sup-keith/builds/123/jobs/job-1/log",
  "command": "scripts/build.sh",
  "soft_failed": false,
  "exit_status": 0
}`)
	})

	job, _, err := client.Jobs.GetJobByOrg(context.Background(), "my-great-org", "job-1")
	if err != nil {
		t.Errorf("GetJobByOrg returned error: %v", err)
	}

	exitStatus := 0
	want := Job{
		ID:         "job-1",
		Type:       "script",
		Name:       ":package: Build",
		State:      "passed",
		BuildURL:   "https://api.buildkite.com/v2/organizations/my-great-org/pipelines/sup-keith/builds/123",
		WebURL:     "https://buildkite.com/my-great-org/sup-keith/builds/123#job-1",
		LogURL:     "https://api.buildkite.com/v2/organizations/my-great-org/pipelines/sup-keith/builds/123/jobs/job-1/log",
		Command:    "scripts/build.sh",
		SoftFailed: false,
		ExitStatus: &exitStatus,
	}
	if diff := cmp.Diff(job, want); diff != "" {
		t.Errorf("GetJobByOrg diff: (-got +want)\n%s", diff)
	}
}

func TestJobsService_UnblockJob(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-great-org/pipelines/sup-keith/builds/awesome-build/jobs/awesome-job-id/unblock", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		_, _ = fmt.Fprint(w, `{
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
		_, _ = fmt.Fprint(w, `{
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
		_, _ = fmt.Fprint(w, `{
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

func TestJobsService_JobLogExists(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-great-org/pipelines/sup-keith/builds/awesome-build/jobs/awesome-job-id/log", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodHead)
		w.Header().Set("Content-Length", "28")
		w.Header().Set("Accept-Ranges", "bytes")
	})

	exists, resp, err := client.Jobs.JobLogExists(context.Background(), "my-great-org", "sup-keith", "awesome-build", "awesome-job-id")
	if err != nil {
		t.Fatalf("JobLogExists returned error: %v", err)
	}
	if !exists {
		t.Error("JobLogExists returned false, want true")
	}
	if got := resp.Header.Get("Content-Length"); got != "28" {
		t.Errorf("Content-Length = %q, want %q", got, "28")
	}
	if got := resp.Header.Get("Accept-Ranges"); got != "bytes" {
		t.Errorf("Accept-Ranges = %q, want %q", got, "bytes")
	}
}

func TestJobsService_JobLogExists_NotFound(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-great-org/pipelines/sup-keith/builds/awesome-build/jobs/missing-job-id/log", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodHead)
		w.WriteHeader(http.StatusNotFound)
	})

	exists, resp, err := client.Jobs.JobLogExists(context.Background(), "my-great-org", "sup-keith", "awesome-build", "missing-job-id")
	if err != nil {
		t.Fatalf("JobLogExists returned error: %v", err)
	}
	if exists {
		t.Error("JobLogExists returned true, want false")
	}
	if resp == nil || resp.StatusCode != http.StatusNotFound {
		t.Fatalf("response = %#v, want status %d", resp, http.StatusNotFound)
	}
}

func TestJobsService_JobLogExists_ServerError(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-great-org/pipelines/sup-keith/builds/awesome-build/jobs/awesome-job-id/log", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodHead)
		w.WriteHeader(http.StatusInternalServerError)
	})

	exists, resp, err := client.Jobs.JobLogExists(context.Background(), "my-great-org", "sup-keith", "awesome-build", "awesome-job-id")
	if err == nil {
		t.Fatal("JobLogExists returned nil error, want an error")
	}
	if exists {
		t.Error("JobLogExists returned true, want false")
	}
	if resp == nil || resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("response = %#v, want status %d", resp, http.StatusInternalServerError)
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

		_, _ = fmt.Fprint(w, string(bytes))
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

func TestJobsService_ReprioritizeJob(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	input := JobReprioritizationOptions{Priority: 10}

	server.HandleFunc("/v2/organizations/my-great-org/pipelines/sup-keith/builds/awesome-build/jobs/awesome-job-id/reprioritize", func(w http.ResponseWriter, r *http.Request) {
		var v JobReprioritizationOptions
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("Error parsing json body: %v", err)
		}

		testMethod(t, r, "PUT")

		if diff := cmp.Diff(v, input); diff != "" {
			t.Errorf("Request body diff: (-got +want)\n%s", diff)
		}

		_, _ = fmt.Fprint(w, `{
  "id": "awesome-job-id",
  "state": "scheduled",
  "priority": {"number": 10}
}`)
	})

	job, _, err := client.Jobs.ReprioritizeJob(context.Background(), "my-great-org", "sup-keith", "awesome-build", "awesome-job-id", &input)
	if err != nil {
		t.Errorf("ReprioritizeJob returned error: %v", err)
	}

	want := Job{ID: "awesome-job-id", State: "scheduled", Priority: &JobPriority{Number: 10}}
	if diff := cmp.Diff(job, want); diff != "" {
		t.Errorf("ReprioritizeJob diff: (-got +want)\n%s", diff)
	}
}

func TestJobsService_Delete(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-great-org/pipelines/my-great-pipeline-slug/builds/69/jobs/420/log", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	_, err := client.Jobs.DeleteJobLog(context.Background(), "my-great-org", "my-great-pipeline-slug", "69", "420")
	if err != nil {
		t.Errorf("Jobs.DeleteJobLog returned error: %v", err)
	}
}

func TestJob_PromisedExitStatusJSON(t *testing.T) {
	t.Parallel()

	const payload = `{
  "id": "awesome-job-id",
  "state": "running",
  "promised_exit_status": 1,
  "promised_exit_status_at": "2026-06-03T04:15:41.618Z"
}`

	var job Job
	if err := json.Unmarshal([]byte(payload), &job); err != nil {
		t.Fatalf("Unmarshal returned error: %v", err)
	}

	if job.PromisedExitStatus == nil || *job.PromisedExitStatus != 1 {
		t.Errorf("PromisedExitStatus = %v, want 1", job.PromisedExitStatus)
	}
	if job.PromisedExitStatusAt == nil {
		t.Fatal("PromisedExitStatusAt = nil, want a timestamp")
	}
	if got := job.PromisedExitStatusAt.UTC().Format("2006-01-02T15:04:05Z"); got != "2026-06-03T04:15:41Z" {
		t.Errorf("PromisedExitStatusAt = %s, want 2026-06-03T04:15:41Z", got)
	}

	// Absent fields must stay nil (omitempty round-trip).
	var empty Job
	if err := json.Unmarshal([]byte(`{"id":"x","state":"passed"}`), &empty); err != nil {
		t.Fatalf("Unmarshal returned error: %v", err)
	}
	if empty.PromisedExitStatus != nil || empty.PromisedExitStatusAt != nil {
		t.Errorf("expected nil promised fields, got %v / %v", empty.PromisedExitStatus, empty.PromisedExitStatusAt)
	}
}
