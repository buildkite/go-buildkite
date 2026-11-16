package buildkite

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

const (
	testPoolPipelineID = "019a1dd3-3870-7e50-8a6d-f16efba96342"
	testPoolBuildID    = "019a1dd8-a49b-70c5-99b8-b8f9f0ed6824"
	testPoolID         = "019a1ddb-a210-7b03-b4a6-5fcd44828b31"
	testPoolPolicyJSON = `{"min_attempts":1,"max_attempts":5,"min_passed":2,"max_passed":3,"min_failed":null,"max_failed":4,"parallel_attempts":2,"initial_attempts":null}`
)

func TestTestPoolsService_List(t *testing.T) {
	t.Parallel()
	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)
	path := "/v2/organizations/acme/test-scheduler/pools"
	baseLink := "https://api.buildkite.com" + path + "?pipeline_id=" + testPoolPipelineID + "&build_id=" + testPoolBuildID + "&key=RSpec%3AMain&per_page=30"
	next := TestPoolsListLink(baseLink + "&after=opaque%2B%2F%3D")
	previous := TestPoolsListLink(baseLink + "&before=back%2B%2F%3D")
	server.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		want := values{"pipeline_id": testPoolPipelineID, "build_id": testPoolBuildID, "key": "RSpec:Main"}
		switch len(server.calls) {
		case 1:
			// Zero PerPage lets the API apply its default. No page-number option.
		case 2:
			want["after"], want["per_page"] = "opaque+/=", "30"
		case 3:
			want["before"], want["per_page"] = "back+/=", "30"
		default:
			t.Error("List fetched an unrequested page")
		}
		testFormValues(t, r, want)
		if len(server.calls) > 1 {
			_, _ = fmt.Fprint(w, `{"items":[],"links":{}}`)
			return
		}
		_, _ = fmt.Fprintf(w, `{"items":[{
			"id":%q,"organization_id":"org-id","suite_id":"suite-id","pipeline_id":%q,"build_id":%q,
			"key":"RSpec:Main","state":"planning","attempt_policy":%s,
			"lease":{"costs":{"custom":25.5},"max_attempts":100},
			"expires_at":"2026-11-16T23:00:00Z","created_at":"2026-11-16T22:00:00Z"
		}],"links":{"self":%q,"first":%q,"next":%q,"prev":%q}}`,
			testPoolID, testPoolPipelineID, testPoolBuildID, testPoolPolicyJSON, baseLink, baseLink, next, previous)
	})
	page, _, err := client.TestPools.List(context.Background(), "acme", &TestPoolsListOptions{
		PipelineID: testPoolPipelineID, BuildID: testPoolBuildID, Key: "RSpec:Main",
	})
	if err != nil {
		t.Fatal(err)
	}
	minPassed := 2
	want := TestPoolsList{
		Items: []TestPool{{
			ID: testPoolID, OrganizationID: "org-id", SuiteID: "suite-id", PipelineID: testPoolPipelineID, BuildID: testPoolBuildID,
			Key: "RSpec:Main", State: "planning",
			AttemptPolicy: TestPoolAttemptPolicies{Single: &TestPoolAttemptPolicy{
				MinAttempts: 1, MaxAttempts: 5, MinPassed: &minPassed, MaxPassed: 3,
				MaxFailed: 4, ParallelAttempts: 2,
			}},
			Lease:     TestPoolLeasePolicy{Costs: TestPoolCosts{"custom": 25.5}, MaxAttempts: 100},
			ExpiresAt: NewTimestamp(time.Date(2026, 11, 16, 23, 0, 0, 0, time.UTC)),
			CreatedAt: NewTimestamp(time.Date(2026, 11, 16, 22, 0, 0, 0, time.UTC)),
		}},
		Links: TestPoolsListLinks{Self: TestPoolsListLink(baseLink), First: TestPoolsListLink(baseLink), Next: next, Previous: previous},
	}
	if diff := cmp.Diff(want, page); diff != "" {
		t.Errorf("List (-want +got):\n%s", diff)
	}
	if len(server.calls) != 1 {
		t.Fatalf("List made %d requests, want one", len(server.calls))
	}
	policyJSON, err := json.Marshal(page.Items[0].AttemptPolicy)
	if err != nil {
		t.Fatal(err)
	}
	assertJSONEqual(t, string(policyJSON), testPoolPolicyJSON)
	for _, link := range []TestPoolsListLink{page.Links.Next, page.Links.Previous} {
		opt, err := link.ToOptions()
		if err != nil {
			t.Fatal(err)
		}
		empty, _, err := client.TestPools.List(context.Background(), "acme", opt)
		if err != nil {
			t.Fatal(err)
		}
		if empty.Items == nil || len(empty.Items) != 0 || empty.Links.Next != "" {
			t.Errorf("unexpected empty terminal page: %+v", empty)
		}
	}
}

func TestTestPoolsService_Get(t *testing.T) {
	t.Parallel()
	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)
	namedJSON := `{"default":` + testPoolPolicyJSON + `,"retry":{"min_attempts":2,"max_attempts":6,"min_passed":null,"max_passed":4,"min_failed":1,"max_failed":5,"parallel_attempts":3,"initial_attempts":2}}`
	server.HandleFunc("/v2/organizations/acme/test-scheduler/pools/"+testPoolID, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testFormValues(t, r, values{"build_id": testPoolBuildID})
		_, _ = fmt.Fprintf(w, `{"id":%q,"state":"errored","error":{"message":"Test pool failed."},
			"attempt_policy":%s,"lease":{"costs":{"duration_p90_ms":120000},"max_attempts":25},
			"muted_tests":[{"scope":null,"name":"works","path":"user_spec.rb:42"}],"parallelism":0}`, testPoolID, namedJSON)
	})
	pool, _, err := client.TestPools.Get(context.Background(), "acme", testPoolID, testPoolBuildID)
	if err != nil {
		t.Fatal(err)
	}
	minPassed, minFailed, initial, parallelism := 2, 1, 2, 0
	name, path := "works", "user_spec.rb:42"
	want := TestPool{
		ID: testPoolID, State: "errored", Error: &TestPoolError{Message: "Test pool failed."},
		AttemptPolicy: TestPoolAttemptPolicies{Named: map[string]TestPoolAttemptPolicy{
			"default": {MinAttempts: 1, MaxAttempts: 5, MinPassed: &minPassed, MaxPassed: 3, MaxFailed: 4, ParallelAttempts: 2},
			"retry":   {MinAttempts: 2, MaxAttempts: 6, MaxPassed: 4, MinFailed: &minFailed, MaxFailed: 5, ParallelAttempts: 3, InitialAttempts: &initial},
		}},
		Lease:       TestPoolLeasePolicy{Costs: TestPoolCosts{"duration_p90_ms": 120000}, MaxAttempts: 25},
		MutedTests:  []TestPoolMutedTest{{Name: &name, Path: &path}},
		Parallelism: &parallelism,
	}
	if diff := cmp.Diff(want, pool); diff != "" {
		t.Errorf("Get (-want +got):\n%s", diff)
	}
	policyJSON, err := json.Marshal(pool.AttemptPolicy)
	if err != nil {
		t.Fatal(err)
	}
	assertJSONEqual(t, string(policyJSON), namedJSON)
}

func TestTestPoolsService_GetMetrics(t *testing.T) {
	t.Parallel()
	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)
	server.HandleFunc("/v2/organizations/acme/test-scheduler/pools/"+testPoolID+"/metrics", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testFormValues(t, r, values{"build_id": testPoolBuildID})
		if len(server.calls) == 2 {
			_, _ = fmt.Fprint(w, `{
				"pool":{"state":"consuming","drained":false},
				"leases":{"active":{"count":2,"attempt_count":7}},
				"timing":{"oldest_waiting_attempt_created_at":null,
					"oldest_active_lease_started_at":"2026-11-16T22:02:00Z",
					"oldest_expired_lease_started_at":null,"next_active_lease_expires_at":"2026-11-16T22:04:00Z"}
			}`)
			return
		}
		_, _ = fmt.Fprintf(w, `{
			"observed_at":"2026-11-16T22:03:00.903Z",
			"pool":{"id":%q,"state":"errored","drained":true,"expires_at":"2026-11-16T23:00:00Z"},
			"entries":{"total":5},
			"attempts":{"total":9,"states":{
				"waiting":{"count":3,"costs":{"duration_p50_ms":7000}},
				"leased":{"count":0,"costs":{"duration_p50_ms":0}},
				"completed":{"count":5,"costs":{"duration_p50_ms":11000}},
				"canceled":{"count":1,"costs":{"duration_p50_ms":2500}}
			},"results":{"passed":2,"failed":1,"errored":2}},
			"leases":{"active":{"count":0,"attempt_count":0},"expired":{"count":2,"attempt_count":3}},
			"timing":{"oldest_waiting_attempt_created_at":"2026-11-16T22:00:00Z",
				"oldest_active_lease_started_at":null,"oldest_expired_lease_started_at":"2026-11-16T22:01:00Z",
				"next_active_lease_expires_at":null}
		}`, testPoolID)
	})
	got, _, err := client.TestPools.GetMetrics(context.Background(), "acme", testPoolID, testPoolBuildID)
	if err != nil {
		t.Fatal(err)
	}
	want := TestPoolMetrics{
		ObservedAt: NewTimestamp(time.Date(2026, 11, 16, 22, 3, 0, 903000000, time.UTC)),
		Pool: TestPoolMetricsPool{ID: testPoolID, State: "errored", Drained: true,
			ExpiresAt: NewTimestamp(time.Date(2026, 11, 16, 23, 0, 0, 0, time.UTC))},
		Entries: TestPoolMetricsEntries{Total: 5},
		Attempts: TestPoolMetricsAttempts{Total: 9,
			States: map[string]TestPoolMetricsAttemptState{
				"waiting":   {Count: 3, Costs: TestPoolCosts{"duration_p50_ms": 7000}},
				"leased":    {Count: 0, Costs: TestPoolCosts{"duration_p50_ms": 0}},
				"completed": {Count: 5, Costs: TestPoolCosts{"duration_p50_ms": 11000}},
				"canceled":  {Count: 1, Costs: TestPoolCosts{"duration_p50_ms": 2500}},
			},
			Results: map[string]int64{"passed": 2, "failed": 1, "errored": 2},
		},
		Leases: TestPoolMetricsLeases{Expired: TestPoolMetricsLease{Count: 2, AttemptCount: 3}},
		Timing: TestPoolMetricsTiming{
			OldestWaitingAttemptCreatedAt: NewTimestamp(time.Date(2026, 11, 16, 22, 0, 0, 0, time.UTC)),
			OldestExpiredLeaseStartedAt:   NewTimestamp(time.Date(2026, 11, 16, 22, 1, 0, 0, time.UTC)),
		},
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("GetMetrics (-want +got):\n%s", diff)
	}

	// The active snapshot exercises the other timing fields and nonzero active
	// leases, rather than letting omitted or mis-tagged fields pass as zero.
	active, _, err := client.TestPools.GetMetrics(context.Background(), "acme", testPoolID, testPoolBuildID)
	if err != nil {
		t.Fatal(err)
	}
	wantActive := TestPoolMetrics{
		Pool:   TestPoolMetricsPool{State: "consuming"},
		Leases: TestPoolMetricsLeases{Active: TestPoolMetricsLease{Count: 2, AttemptCount: 7}},
		Timing: TestPoolMetricsTiming{
			OldestActiveLeaseStartedAt: NewTimestamp(time.Date(2026, 11, 16, 22, 2, 0, 0, time.UTC)),
			NextActiveLeaseExpiresAt:   NewTimestamp(time.Date(2026, 11, 16, 22, 4, 0, 0, time.UTC)),
		},
	}
	if diff := cmp.Diff(wantActive, active); diff != "" {
		t.Errorf("active GetMetrics (-want +got):\n%s", diff)
	}
}

func TestTestPoolsService_APIErrors(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name, suffix, message string
		status                int
	}{
		{"list invalid cursors", "", "Provide only one of `after` or `before`", http.StatusBadRequest},
		{"list missing scope", "", "Your access token doesn't have the read_test_pool scope", http.StatusForbidden},
		{"get unavailable", "/" + testPoolID, "No test pool found", http.StatusNotFound},
		{"metrics rate limited", "/" + testPoolID + "/metrics", "Rate limit exceeded", http.StatusTooManyRequests},
		{"metrics service unavailable", "/" + testPoolID + "/metrics", "Test Scheduler is temporarily unavailable", http.StatusServiceUnavailable},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			server, client, teardown := newMockServerAndClient(t)
			t.Cleanup(teardown)
			// Check service error propagation, not the shared client's retry policy.
			client = must(NewClient(WithBaseURL(client.BaseURL.String()), WithTokenAuth("fixture-token"), WithMaxRetries(0)))
			body := fmt.Sprintf(`{"message":%q}`, tc.message)
			server.HandleFunc("/v2/organizations/acme/test-scheduler/pools"+tc.suffix, func(w http.ResponseWriter, r *http.Request) {
				if got := r.Header.Get("Authorization"); got != "Bearer fixture-token" {
					t.Errorf("Authorization = %q", got)
				}
				w.WriteHeader(tc.status)
				_, _ = fmt.Fprint(w, body)
			})
			var resp *Response
			var err error
			switch tc.suffix {
			case "":
				_, resp, err = client.TestPools.List(context.Background(), "acme", &TestPoolsListOptions{
					PipelineID: testPoolPipelineID, BuildID: testPoolBuildID, After: "a", Before: "b",
				})
			case "/" + testPoolID:
				_, resp, err = client.TestPools.Get(context.Background(), "acme", testPoolID, testPoolBuildID)
			default:
				_, resp, err = client.TestPools.GetMetrics(context.Background(), "acme", testPoolID, testPoolBuildID)
			}
			var apiErr *ErrorResponse
			if !errors.As(err, &apiErr) {
				t.Fatalf("error = %v, want *ErrorResponse", err)
			}
			if resp == nil || resp.StatusCode != tc.status || apiErr.Response.StatusCode != tc.status || apiErr.Message != tc.message || string(apiErr.RawBody) != body {
				t.Errorf("response = %+v, error = %+v; want status %d and body %s", resp, apiErr, tc.status, body)
			}
		})
	}
}
