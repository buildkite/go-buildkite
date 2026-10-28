package buildkite

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestSpanAggregatesService_ListScopesBySuite(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name  string
		opt   *SpanAggregateScopesListOptions
		query string
		body  string
		want  []SpanAggregateScope
	}{
		{
			name:  "period and totals",
			opt:   &SpanAggregateScopesListOptions{Period: "7days"},
			query: "period=7days",
			body:  `[{"name":"pg","span_count":3,"duration_ns":9000000001},{"name":"","span_count":1,"duration_ns":42}]`,
			want: []SpanAggregateScope{
				{Name: "pg", SpanCount: 3, Duration: 9*time.Second + time.Nanosecond},
				{Name: "", SpanCount: 1, Duration: 42 * time.Nanosecond},
			},
		},
		{
			name: "default period and empty result",
			body: `[]`,
			want: []SpanAggregateScope{},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server, client, teardown := newMockServerAndClient(t)
			t.Cleanup(teardown)

			server.HandleFunc("/v2/analytics/organizations/my-great-org/suites/suite-example/span-aggregates/scopes", func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "GET")
				if got := r.URL.RawQuery; got != tt.query {
					t.Errorf("query = %q, want %q", got, tt.query)
				}
				_, _ = fmt.Fprint(w, tt.body)
			})

			got, resp, err := client.SpanAggregates.ListScopesBySuite(context.Background(), "my-great-org", "suite-example", tt.opt)
			if err != nil {
				t.Fatalf("ListScopesBySuite returned error: %v", err)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("ListScopesBySuite mismatch (-want +got):\n%s", diff)
			}
			if resp.StatusCode != http.StatusOK {
				t.Errorf("response status = %d, want 200", resp.StatusCode)
			}
		})
	}
}

func TestSpanAggregatesService_ListNamesByBuild(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name      string
		scopeName string
		opt       *SpanAggregateNamesListOptions
		query     string
		body      string
		want      SpanAggregateNames
	}{
		{
			name:      "encoded scope and fractional statistics",
			scopeName: "pg /&+?",
			opt:       &SpanAggregateNamesListOptions{Sort: "-p95", Limit: 1},
			query:     "limit=1&scope_name=pg+%2F%26%2B%3F&sort=-p95",
			body: `{"scope_name":"pg /&+?","span_names":[{
				"name":"SELECT users","span_count":3,"duration_ns":9000000001,
				"avg_duration_ns":3000000000.3333335,"p50_duration_ns":2500000000.5,
				"p95_duration_ns":4500000000.75,"max_duration_ns":5000000001
			}],"span_names_truncated":true}`,
			want: SpanAggregateNames{
				ScopeName: "pg /&+?",
				SpanNames: []SpanAggregateName{{
					Name: "SELECT users", SpanCount: 3, Duration: 9*time.Second + time.Nanosecond,
					AverageDuration: 3000000000.3333335, P50Duration: 2500000000.5,
					P95Duration: 4500000000.75, MaximumDuration: 5*time.Second + time.Nanosecond,
				}},
				SpanNamesTruncated: true,
			},
		},
		{
			name:  "unnamed scope and default options",
			query: "scope_name=",
			body:  `{"scope_name":"","span_names":[],"span_names_truncated":false}`,
			want:  SpanAggregateNames{SpanNames: []SpanAggregateName{}},
		},
		{
			name:      "unknown scope",
			scopeName: "unknown",
			opt:       &SpanAggregateNamesListOptions{},
			query:     "scope_name=unknown",
			body:      `{"scope_name":"unknown","span_names":[],"span_names_truncated":false}`,
			want:      SpanAggregateNames{ScopeName: "unknown", SpanNames: []SpanAggregateName{}},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server, client, teardown := newMockServerAndClient(t)
			t.Cleanup(teardown)

			server.HandleFunc(fmt.Sprintf("/v2/analytics/organizations/my-great-org/builds/%s/span-aggregates/span-names", testBuildUUID), func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "GET")
				if got := r.URL.RawQuery; got != tt.query {
					t.Errorf("query = %q, want %q", got, tt.query)
				}
				_, _ = fmt.Fprint(w, tt.body)
			})

			got, resp, err := client.SpanAggregates.ListNamesByBuild(context.Background(), "my-great-org", testBuildUUID, tt.scopeName, tt.opt)
			if err != nil {
				t.Fatalf("ListNamesByBuild returned error: %v", err)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("ListNamesByBuild mismatch (-want +got):\n%s", diff)
			}
			if resp.StatusCode != http.StatusOK {
				t.Errorf("response status = %d, want 200", resp.StatusCode)
			}
		})
	}
}

func TestSpanAggregatesService_ListScopesByBuild(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name string
		body string
		want []SpanAggregateScope
	}{
		{
			name: "scope totals",
			body: `[{"name":"pg","span_count":3,"duration_ns":9000000001},{"name":"","span_count":1,"duration_ns":42}]`,
			want: []SpanAggregateScope{
				{Name: "pg", SpanCount: 3, Duration: 9*time.Second + time.Nanosecond},
				{Name: "", SpanCount: 1, Duration: 42 * time.Nanosecond},
			},
		},
		{name: "empty result", body: `[]`, want: []SpanAggregateScope{}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server, client, teardown := newMockServerAndClient(t)
			t.Cleanup(teardown)
			server.HandleFunc(fmt.Sprintf("/v2/analytics/organizations/my-great-org/builds/%s/span-aggregates/scopes", testBuildUUID), func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "GET")
				if r.URL.RawQuery != "" {
					t.Errorf("unexpected query: %q", r.URL.RawQuery)
				}
				_, _ = fmt.Fprint(w, tt.body)
			})

			got, resp, err := client.SpanAggregates.ListScopesByBuild(context.Background(), "my-great-org", testBuildUUID)
			if err != nil {
				t.Fatalf("ListScopesByBuild returned error: %v", err)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("ListScopesByBuild mismatch (-want +got):\n%s", diff)
			}
			if resp.StatusCode != http.StatusOK {
				t.Errorf("response status = %d, want 200", resp.StatusCode)
			}
		})
	}
}

func TestSpanAggregatesService_ListNamesBySuite(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name      string
		scopeName string
		opt       *SpanAggregateSuiteNamesListOptions
		query     string
		body      string
		want      SpanAggregateNames
	}{
		{
			name:      "period, encoded scope and fractional statistics",
			scopeName: "pg /&+?",
			opt: &SpanAggregateSuiteNamesListOptions{
				Period:                        "7days",
				SpanAggregateNamesListOptions: SpanAggregateNamesListOptions{Sort: "-p95", Limit: 1},
			},
			query: "limit=1&period=7days&scope_name=pg+%2F%26%2B%3F&sort=-p95",
			body: `{"scope_name":"pg /&+?","span_names":[{
				"name":"SELECT users","span_count":3,"duration_ns":9000000001,
				"avg_duration_ns":3000000000.3333335,"p50_duration_ns":2500000000.5,
				"p95_duration_ns":4500000000.75,"max_duration_ns":5000000001
			}],"span_names_truncated":true}`,
			want: SpanAggregateNames{
				ScopeName: "pg /&+?",
				SpanNames: []SpanAggregateName{{
					Name: "SELECT users", SpanCount: 3, Duration: 9*time.Second + time.Nanosecond,
					AverageDuration: 3000000000.3333335, P50Duration: 2500000000.5,
					P95Duration: 4500000000.75, MaximumDuration: 5*time.Second + time.Nanosecond,
				}},
				SpanNamesTruncated: true,
			},
		},
		{
			name:  "unnamed scope and default options",
			query: "scope_name=",
			body:  `{"scope_name":"","span_names":[],"span_names_truncated":false}`,
			want:  SpanAggregateNames{SpanNames: []SpanAggregateName{}},
		},
		{
			name:      "unknown scope and empty options",
			scopeName: "unknown",
			opt:       &SpanAggregateSuiteNamesListOptions{},
			query:     "scope_name=unknown",
			body:      `{"scope_name":"unknown","span_names":[],"span_names_truncated":false}`,
			want:      SpanAggregateNames{ScopeName: "unknown", SpanNames: []SpanAggregateName{}},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server, client, teardown := newMockServerAndClient(t)
			t.Cleanup(teardown)
			server.HandleFunc("/v2/analytics/organizations/my-great-org/suites/suite-example/span-aggregates/span-names", func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "GET")
				if got := r.URL.RawQuery; got != tt.query {
					t.Errorf("query = %q, want %q", got, tt.query)
				}
				_, _ = fmt.Fprint(w, tt.body)
			})

			got, resp, err := client.SpanAggregates.ListNamesBySuite(context.Background(), "my-great-org", "suite-example", tt.scopeName, tt.opt)
			if err != nil {
				t.Fatalf("ListNamesBySuite returned error: %v", err)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("ListNamesBySuite mismatch (-want +got):\n%s", diff)
			}
			if resp.StatusCode != http.StatusOK {
				t.Errorf("response status = %d, want 200", resp.StatusCode)
			}
		})
	}
}

func TestSpanAggregatesService_APIError(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name string
		call func(*Client) (*Response, error)
	}{
		{
			name: "suite scopes",
			call: func(client *Client) (*Response, error) {
				_, resp, err := client.SpanAggregates.ListScopesBySuite(context.Background(), "my-great-org", "suite-example", &SpanAggregateScopesListOptions{Period: "invalid"})
				return resp, err
			},
		},
		{
			name: "build names",
			call: func(client *Client) (*Response, error) {
				_, resp, err := client.SpanAggregates.ListNamesByBuild(context.Background(), "my-great-org", testBuildUUID, "pg", &SpanAggregateNamesListOptions{Limit: 101})
				return resp, err
			},
		},
		{
			name: "build scopes",
			call: func(client *Client) (*Response, error) {
				_, resp, err := client.SpanAggregates.ListScopesByBuild(context.Background(), "my-great-org", testBuildUUID)
				return resp, err
			},
		},
		{
			name: "suite names",
			call: func(client *Client) (*Response, error) {
				_, resp, err := client.SpanAggregates.ListNamesBySuite(context.Background(), "my-great-org", "suite-example", "pg", &SpanAggregateSuiteNamesListOptions{Period: "invalid"})
				return resp, err
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server, client, teardown := newMockServerAndClient(t)
			t.Cleanup(teardown)
			server.HandleFunc("/v2/analytics/", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusUnprocessableEntity)
				_, _ = fmt.Fprint(w, `{"message":"invalid parameter"}`)
			})

			resp, err := tt.call(client)
			var apiErr *ErrorResponse
			if !errors.As(err, &apiErr) {
				t.Fatalf("error = %v, want *ErrorResponse", err)
			}
			if resp.StatusCode != http.StatusUnprocessableEntity {
				t.Errorf("response status = %d, want 422", resp.StatusCode)
			}
		})
	}
}
