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

func TestOrganizationsService_GetActiveUsersUsage(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		org      string
		path     string
		opt      ActiveUsersUsageOptions
		query    values
		body     string
		scope    string
		interval string
		summary  ActiveUsersUsageSummary
		items    []ActiveUsersUsageItem
	}{
		{
			name:  "default_daily",
			org:   "acme",
			path:  "/v2/organizations/acme/usage/active-users",
			opt:   ActiveUsersUsageOptions{From: "2026-10-22", To: "2026-10-25"},
			query: values{"from": "2026-10-22", "to": "2026-10-25"},
			body: `{
				"description":"Usage can change as late events arrive.",
				"reporting_period":{"from":"2026-10-22","to":"2026-10-25"},
				"generated_at":"2026-10-27T03:04:05Z",
				"scope":"organization_and_billing_children","interval":"day","unit":"users",
				"definition":"Distinct active users, deduplicated across billing children.",
				"summary":{"active_users":4,"average_monthly_active_users":4},
				"items":[
					{"date":"2026-10-22","active_users":3,"cumulative_active_users":3},
					{"date":"2026-10-23","active_users":0,"cumulative_active_users":3},
					{"date":"2026-10-24","active_users":2,"cumulative_active_users":4},
					{"date":"2026-10-25","active_users":0,"cumulative_active_users":4}
				]
			}`,
			scope:    "organization_and_billing_children",
			interval: "day",
			summary:  ActiveUsersUsageSummary{ActiveUsers: 4, AverageMonthlyActiveUsers: 4},
			items: []ActiveUsersUsageItem{
				{Date: "2026-10-22", ActiveUsers: 3, CumulativeActiveUsers: 3},
				{Date: "2026-10-23", ActiveUsers: 0, CumulativeActiveUsers: 3},
				{Date: "2026-10-24", ActiveUsers: 2, CumulativeActiveUsers: 4},
				{Date: "2026-10-25", ActiveUsers: 0, CumulativeActiveUsers: 4},
			},
		},
		{
			name:  "monthly_partial_range",
			org:   "acme",
			path:  "/v2/organizations/acme/usage/active-users",
			opt:   ActiveUsersUsageOptions{From: "2026-08-15", To: "2026-10-10", Interval: "month"},
			query: values{"from": "2026-08-15", "to": "2026-10-10", "interval": "month"},
			body: `{
				"description":"Usage can change as late events arrive.",
				"reporting_period":{"from":"2026-08-15","to":"2026-10-10"},
				"generated_at":"2026-10-27T03:04:05Z",
				"scope":"organization_and_billing_children","interval":"month","unit":"users",
				"definition":"Distinct active users, deduplicated across billing children.",
				"summary":{"active_users":7,"average_monthly_active_users":5},
				"items":[
					{"date":"2026-08-01","active_users":4,"cumulative_active_users":4},
					{"date":"2026-09-01","active_users":0,"cumulative_active_users":4},
					{"date":"2026-10-01","active_users":5,"cumulative_active_users":7}
				]
			}`,
			scope:    "organization_and_billing_children",
			interval: "month",
			summary:  ActiveUsersUsageSummary{ActiveUsers: 7, AverageMonthlyActiveUsers: 5},
			items: []ActiveUsersUsageItem{
				{Date: "2026-08-01", ActiveUsers: 4, CumulativeActiveUsers: 4},
				{Date: "2026-09-01", ActiveUsers: 0, CumulativeActiveUsers: 4},
				{Date: "2026-10-01", ActiveUsers: 5, CumulativeActiveUsers: 7},
			},
		},
		{
			name:  "zero_explicit_daily",
			org:   "7c202aaa-3165-4811-9813-173c4c285463",
			path:  "/v2/organizations/7c202aaa-3165-4811-9813-173c4c285463/usage/active-users",
			opt:   ActiveUsersUsageOptions{From: "2026-10-25", To: "2026-10-25", Interval: "day"},
			query: values{"from": "2026-10-25", "to": "2026-10-25", "interval": "day"},
			body: `{
				"description":"Usage can change as late events arrive.",
				"reporting_period":{"from":"2026-10-25","to":"2026-10-25"},
				"generated_at":"2026-10-27T03:04:05Z",
				"scope":"organization","interval":"day","unit":"users",
				"definition":"Distinct active users, deduplicated across billing children.",
				"summary":{"active_users":0,"average_monthly_active_users":0},
				"items":[{"date":"2026-10-25","active_users":0,"cumulative_active_users":0}]
			}`,
			scope:    "organization",
			interval: "day",
			summary:  ActiveUsersUsageSummary{},
			items:    []ActiveUsersUsageItem{{Date: "2026-10-25"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server, client, teardown := newMockServerAndClient(t)
			t.Cleanup(teardown)
			server.HandleFunc(tt.path, func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "GET")
				testFormValues(t, r, tt.query)
				_, _ = fmt.Fprint(w, tt.body)
			})

			got, resp, err := client.Organizations.GetActiveUsersUsage(context.Background(), tt.org, &tt.opt)
			if err != nil {
				t.Fatal(err)
			}
			want := ActiveUsersUsage{
				Description:     "Usage can change as late events arrive.",
				ReportingPeriod: UsageReportingPeriod{From: tt.opt.From, To: tt.opt.To},
				GeneratedAt:     NewTimestamp(time.Date(2026, 10, 27, 3, 4, 5, 0, time.UTC)),
				Scope:           tt.scope,
				Interval:        tt.interval,
				Unit:            "users",
				Definition:      "Distinct active users, deduplicated across billing children.",
				Summary:         tt.summary,
				Items:           tt.items,
			}
			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("active-user usage (-want +got):\n%s", diff)
			}
			if resp == nil || resp.StatusCode != http.StatusOK {
				t.Fatalf("response = %+v; want status 200", resp)
			}

			// Consumers such as MCP serialize reports again. Required zero counts
			// must survive that serialization, including both summary fields.
			data, err := json.Marshal(got)
			if err != nil {
				t.Fatal(err)
			}
			assertJSONEqual(t, string(data), tt.body)
		})
	}
}

func TestOrganizationsService_GetActiveUsersUsageErrors(t *testing.T) {
	t.Parallel()
	for _, status := range []int{http.StatusBadRequest, http.StatusUnauthorized, http.StatusForbidden, http.StatusNotFound, http.StatusServiceUnavailable} {
		t.Run(fmt.Sprint(status), func(t *testing.T) {
			t.Parallel()
			server, client, teardown := newMockServerAndClient(t)
			t.Cleanup(teardown)
			body := fmt.Sprintf(`{"message":%q}`, http.StatusText(status))
			server.HandleFunc("/v2/organizations/acme/usage/active-users", func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "GET")
				// Invalid and partial values are passed through for API validation,
				// not normalized into a different reporting period by the SDK.
				testFormValues(t, r, values{"from": "2026-10", "to": "2026-10-25", "interval": "week & month"})
				w.WriteHeader(status)
				_, _ = fmt.Fprint(w, body)
			})
			got, resp, err := client.Organizations.GetActiveUsersUsage(context.Background(), "acme", &ActiveUsersUsageOptions{
				From: "2026-10", To: "2026-10-25", Interval: "week & month",
			})
			var apiErr *ErrorResponse
			if !errors.As(err, &apiErr) {
				t.Fatalf("error = %v; want *ErrorResponse", err)
			}
			if resp == nil || resp.StatusCode != status || apiErr.Message != http.StatusText(status) || string(apiErr.RawBody) != body {
				t.Errorf("response=%+v error=%+v", resp, apiErr)
			}
			if diff := cmp.Diff(ActiveUsersUsage{}, got); diff != "" {
				t.Errorf("report returned on error (-want +got):\n%s", diff)
			}
		})
	}

	t.Run("missing_options", func(t *testing.T) {
		server, client, teardown := newMockServerAndClient(t)
		t.Cleanup(teardown)
		server.HandleFunc("/v2/organizations/acme/usage/active-users", func(w http.ResponseWriter, r *http.Request) {
			testFormValues(t, r, values{})
			w.WriteHeader(http.StatusBadRequest)
			_, _ = fmt.Fprint(w, `{"message":"from and to are required"}`)
		})
		_, resp, err := client.Organizations.GetActiveUsersUsage(context.Background(), "acme", nil)
		var apiErr *ErrorResponse
		if !errors.As(err, &apiErr) || resp == nil || resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("response=%+v error=%v; want API error 400", resp, err)
		}
	})

	t.Run("canceled", func(t *testing.T) {
		_, client, teardown := newMockServerAndClient(t)
		t.Cleanup(teardown)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		_, _, err := client.Organizations.GetActiveUsersUsage(ctx, "acme", &ActiveUsersUsageOptions{
			From: "2026-10-22", To: "2026-10-25",
		})
		if !errors.Is(err, context.Canceled) {
			t.Errorf("error = %v; want context.Canceled", err)
		}
	})
}
