package buildkite

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBuildsService_GetFailureEvidence(t *testing.T) {
	t.Parallel()
	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)
	server.HandleFunc("/v2/organizations/org/pipelines/pipeline/builds/42/failure-evidence", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		require.Equal(t, "include=test_results&job_limit=3", r.URL.RawQuery)
		_, _ = fmt.Fprint(w, `{"schema_version":1,"build":{"id":"build-id","number":42,"state":"failed"},"job_limit":3,"jobs_truncated":true,"annotations":[],"collected_at":"2026-11-11T01:00:00.123456Z","jobs":[{"id":"job-id","exit_status":2,"step_key":null,"started_at":null,"annotations":[{"id":"annotation-id","job_id":"job-id","scope":"job","body":"Compiler failed","body_truncated":true}],"test_results":{"status":"found","counts":{"failed":1},"failures":[{"name":"invoice","failure_reason":"Expected 42, got 7","failed_attempts":2,"execution_id":"execution-id"}],"window":{"truncated":true}},"omissions":{"logs":"not_collected"}}]}`)
	})
	evidence, resp, err := client.Builds.GetFailureEvidence(context.Background(), "org", "pipeline", "42", &BuildFailureEvidenceOptions{
		JobLimit: 3,
		Include:  "test_results",
	})
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, 42, evidence.Build.Number)
	require.True(t, evidence.JobsTruncated)
	require.Len(t, evidence.Jobs, 1)
	job := evidence.Jobs[0]
	require.Nil(t, job.StepKey)
	require.Nil(t, job.StartedAt)
	require.Equal(t, 2, job.ExitStatus)
	require.True(t, job.Annotations[0].BodyTruncated)
	require.Equal(t, "Expected 42, got 7", job.TestResults.Failures[0].FailureReason)
	require.True(t, job.TestResults.Window.Truncated)
	require.Equal(t, "not_collected", job.Omissions["logs"])
	require.Equal(t, 123456000, evidence.CollectedAt.Nanosecond())
}

func TestBuildsService_GetFailureEvidenceForbidden(t *testing.T) {
	t.Parallel()
	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)
	server.HandleFunc("/v2/organizations/org/pipelines/pipeline/builds/42/failure-evidence", func(w http.ResponseWriter, r *http.Request) {
		require.Empty(t, r.URL.RawQuery)
		w.WriteHeader(http.StatusForbidden)
		_, _ = fmt.Fprint(w, `{"message":"Token requires read_suites"}`)
	})
	_, resp, err := client.Builds.GetFailureEvidence(context.Background(), "org", "pipeline", "42", nil)
	require.Error(t, err)
	require.Equal(t, http.StatusForbidden, resp.StatusCode)
}
