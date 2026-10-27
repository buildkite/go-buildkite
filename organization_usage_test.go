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

func TestOrganizationsService_GetHostedAgentsUsage(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		org      string
		path     string
		opt      HostedAgentsUsageOptions
		query    values
		body     string
		scope    string
		interval string
		summary  HostedAgentsUsageSummary
		items    []HostedAgentsUsageItem
	}{
		{
			name:  "daily_preserves_period_rounding_and_int64",
			org:   "acme",
			path:  "/v2/organizations/acme/usage/hosted-agents",
			opt:   HostedAgentsUsageOptions{From: "2026-10-23", To: "2026-10-25"},
			query: values{"from": "2026-10-23", "to": "2026-10-25"},
			body: `{
				"description":"Hosted usage can change as late events arrive.",
				"reporting_period":{"from":"2026-10-23","to":"2026-10-25"},
				"generated_at":"2026-10-27T03:04:05Z",
				"scope":"organization_and_billing_children","interval":"day","unit":"vCPU-minutes",
				"definition":"CPU-weighted recognized shapes, combined by resource before flooring minutes.",
				"summary":{"vcpu_minutes":2147483660,"vcpu_minutes_by_resource":{
					"hosted_agents_linux_amd64_vcpu_minutes":2147483650,
					"hosted_agents_linux_arm64_vcpu_minutes":5,
					"hosted_agents_macos_m4_vcpu_minutes":5,
					"hosted_agents_windows_amd64_vcpu_minutes":0
				}},
				"items":[
					{"date":"2026-10-23","vcpu_minutes_by_resource":{
						"hosted_agents_linux_amd64_vcpu_minutes":1073741824,
						"hosted_agents_linux_arm64_vcpu_minutes":1,
						"hosted_agents_macos_m4_vcpu_minutes":2,
						"hosted_agents_windows_amd64_vcpu_minutes":0
					}},
					{"date":"2026-10-24","vcpu_minutes_by_resource":{
						"hosted_agents_linux_amd64_vcpu_minutes":0,
						"hosted_agents_linux_arm64_vcpu_minutes":0,
						"hosted_agents_macos_m4_vcpu_minutes":0,
						"hosted_agents_windows_amd64_vcpu_minutes":0
					}},
					{"date":"2026-10-25","vcpu_minutes_by_resource":{
						"hosted_agents_linux_amd64_vcpu_minutes":1073741825,
						"hosted_agents_linux_arm64_vcpu_minutes":3,
						"hosted_agents_macos_m4_vcpu_minutes":2,
						"hosted_agents_windows_amd64_vcpu_minutes":0
					}}
				]
			}`,
			scope:    "organization_and_billing_children",
			interval: "day",
			summary: HostedAgentsUsageSummary{VCPUMinutes: 2147483660, VCPUMinutesByResource: map[string]int64{
				"hosted_agents_linux_amd64_vcpu_minutes":   2147483650,
				"hosted_agents_linux_arm64_vcpu_minutes":   5,
				"hosted_agents_macos_m4_vcpu_minutes":      5,
				"hosted_agents_windows_amd64_vcpu_minutes": 0,
			}},
			items: []HostedAgentsUsageItem{
				{Date: "2026-10-23", VCPUMinutesByResource: map[string]int64{
					"hosted_agents_linux_amd64_vcpu_minutes":   1073741824,
					"hosted_agents_linux_arm64_vcpu_minutes":   1,
					"hosted_agents_macos_m4_vcpu_minutes":      2,
					"hosted_agents_windows_amd64_vcpu_minutes": 0,
				}},
				{Date: "2026-10-24", VCPUMinutesByResource: map[string]int64{
					"hosted_agents_linux_amd64_vcpu_minutes":   0,
					"hosted_agents_linux_arm64_vcpu_minutes":   0,
					"hosted_agents_macos_m4_vcpu_minutes":      0,
					"hosted_agents_windows_amd64_vcpu_minutes": 0,
				}},
				{Date: "2026-10-25", VCPUMinutesByResource: map[string]int64{
					"hosted_agents_linux_amd64_vcpu_minutes":   1073741825,
					"hosted_agents_linux_arm64_vcpu_minutes":   3,
					"hosted_agents_macos_m4_vcpu_minutes":      2,
					"hosted_agents_windows_amd64_vcpu_minutes": 0,
				}},
			},
		},
		{
			name:  "zero_monthly_uuid",
			org:   "7c202aaa-3165-4811-9813-173c4c285463",
			path:  "/v2/organizations/7c202aaa-3165-4811-9813-173c4c285463/usage/hosted-agents",
			opt:   HostedAgentsUsageOptions{From: "2026-10-15", To: "2026-10-25", Interval: "month"},
			query: values{"from": "2026-10-15", "to": "2026-10-25", "interval": "month"},
			body: `{
				"description":"Hosted usage can change as late events arrive.",
				"reporting_period":{"from":"2026-10-15","to":"2026-10-25"},
				"generated_at":"2026-10-27T03:04:05Z",
				"scope":"organization","interval":"month","unit":"vCPU-minutes",
				"definition":"CPU-weighted recognized shapes, combined by resource before flooring minutes.",
				"summary":{"vcpu_minutes":0,"vcpu_minutes_by_resource":{
					"hosted_agents_linux_amd64_vcpu_minutes":0,
					"hosted_agents_linux_arm64_vcpu_minutes":0,
					"hosted_agents_macos_m4_vcpu_minutes":0,
					"hosted_agents_windows_amd64_vcpu_minutes":0
				}},
				"items":[{"date":"2026-10-01","vcpu_minutes_by_resource":{
					"hosted_agents_linux_amd64_vcpu_minutes":0,
					"hosted_agents_linux_arm64_vcpu_minutes":0,
					"hosted_agents_macos_m4_vcpu_minutes":0,
					"hosted_agents_windows_amd64_vcpu_minutes":0
				}}]
			}`,
			scope:    "organization",
			interval: "month",
			summary: HostedAgentsUsageSummary{VCPUMinutesByResource: map[string]int64{
				"hosted_agents_linux_amd64_vcpu_minutes":   0,
				"hosted_agents_linux_arm64_vcpu_minutes":   0,
				"hosted_agents_macos_m4_vcpu_minutes":      0,
				"hosted_agents_windows_amd64_vcpu_minutes": 0,
			}},
			items: []HostedAgentsUsageItem{{Date: "2026-10-01", VCPUMinutesByResource: map[string]int64{
				"hosted_agents_linux_amd64_vcpu_minutes":   0,
				"hosted_agents_linux_arm64_vcpu_minutes":   0,
				"hosted_agents_macos_m4_vcpu_minutes":      0,
				"hosted_agents_windows_amd64_vcpu_minutes": 0,
			}}},
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
			got, resp, err := client.Organizations.GetHostedAgentsUsage(context.Background(), tt.org, &tt.opt)
			if err != nil {
				t.Fatal(err)
			}
			want := HostedAgentsUsage{
				Description:     "Hosted usage can change as late events arrive.",
				ReportingPeriod: UsageReportingPeriod{From: tt.opt.From, To: tt.opt.To},
				GeneratedAt:     NewTimestamp(time.Date(2026, 10, 27, 3, 4, 5, 0, time.UTC)),
				Scope:           tt.scope,
				Interval:        tt.interval,
				Unit:            "vCPU-minutes",
				Definition:      "CPU-weighted recognized shapes, combined by resource before flooring minutes.",
				Summary:         tt.summary,
				Items:           tt.items,
			}
			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("hosted-agent usage (-want +got):\n%s", diff)
			}
			if resp == nil || resp.StatusCode != http.StatusOK {
				t.Fatalf("response = %+v; want status 200", resp)
			}
			data, err := json.Marshal(got)
			if err != nil {
				t.Fatal(err)
			}
			assertJSONEqual(t, string(data), tt.body)
		})
	}
}

func TestOrganizationsService_GetHostedAgentsUsageErrors(t *testing.T) {
	t.Parallel()
	for _, status := range []int{http.StatusBadRequest, http.StatusServiceUnavailable} {
		t.Run(fmt.Sprint(status), func(t *testing.T) {
			t.Parallel()
			server, client, teardown := newMockServerAndClient(t)
			t.Cleanup(teardown)
			body := fmt.Sprintf(`{"message":%q}`, http.StatusText(status))
			server.HandleFunc("/v2/organizations/acme/usage/hosted-agents", func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "GET")
				testFormValues(t, r, values{})
				w.WriteHeader(status)
				_, _ = fmt.Fprint(w, body)
			})
			got, resp, err := client.Organizations.GetHostedAgentsUsage(context.Background(), "acme", nil)
			var apiErr *ErrorResponse
			if !errors.As(err, &apiErr) || resp == nil || resp.StatusCode != status {
				t.Fatalf("response=%+v error=%v; want API error %d", resp, err, status)
			}
			if apiErr.Message != http.StatusText(status) || string(apiErr.RawBody) != body {
				t.Errorf("error=%+v", apiErr)
			}
			if diff := cmp.Diff(HostedAgentsUsage{}, got); diff != "" {
				t.Errorf("report returned on error (-want +got):\n%s", diff)
			}
		})
	}
}

func TestOrganizationsService_GetConcurrentAgentsUsage(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		org        string
		path       string
		opt        ConcurrentAgentsUsageOptions
		query      values
		body       string
		scope      string
		interval   string
		agentType  string
		percentile int
		summary    ConcurrentAgentsUsageSummary
		items      []ConcurrentAgentsUsageItem
	}{
		{
			name:  "defaults",
			org:   "acme",
			path:  "/v2/organizations/acme/usage/concurrent-agents",
			opt:   ConcurrentAgentsUsageOptions{From: "2026-10-23", To: "2026-10-25"},
			query: values{"from": "2026-10-23", "to": "2026-10-25"},
			body: `{
				"description":"Recorded usage can change as late events arrive.",
				"reporting_period":{"from":"2026-10-23","to":"2026-10-25"},
				"generated_at":"2026-10-27T03:04:05Z",
				"scope":"organization_and_billing_children","interval":"day","unit":"agents",
				"definition":"Discrete percentile across zero-filled daily sums of organization peaks.",
				"agent_type":"self_hosted","percentile":95,
				"summary":{"max_concurrent_agents":13,"percentile_concurrent_agents":13},
				"items":[
					{"date":"2026-10-23","max_concurrent_agents":9},
					{"date":"2026-10-24","max_concurrent_agents":0},
					{"date":"2026-10-25","max_concurrent_agents":13}
				]
			}`,
			scope:      "organization_and_billing_children",
			interval:   "day",
			agentType:  "self_hosted",
			percentile: 95,
			summary:    ConcurrentAgentsUsageSummary{MaxConcurrentAgents: 13, PercentileConcurrentAgents: 13},
			items: []ConcurrentAgentsUsageItem{
				{Date: "2026-10-23", MaxConcurrentAgents: 9},
				{Date: "2026-10-24", MaxConcurrentAgents: 0},
				{Date: "2026-10-25", MaxConcurrentAgents: 13},
			},
		},
		{
			name: "hosted_monthly_preserves_daily_percentile",
			org:  "acme",
			path: "/v2/organizations/acme/usage/concurrent-agents",
			opt: ConcurrentAgentsUsageOptions{
				From: "2026-09-21", To: "2026-10-11", Interval: "month", Percentile: 95, AgentType: "hosted",
			},
			query: values{"from": "2026-09-21", "to": "2026-10-11", "interval": "month", "percentile": "95", "filter[agent_type]": "hosted"},
			// Twenty zero days followed by one day with 17 agents have a discrete
			// p95 of zero. The p95 of the two monthly bars would instead be 17.
			body: `{
				"description":"Recorded usage can change as late events arrive.",
				"reporting_period":{"from":"2026-09-21","to":"2026-10-11"},
				"generated_at":"2026-10-27T03:04:05Z",
				"scope":"organization_and_billing_children","interval":"month","unit":"agents",
				"definition":"Discrete percentile across zero-filled daily sums of organization peaks.",
				"agent_type":"hosted","percentile":95,
				"summary":{"max_concurrent_agents":17,"percentile_concurrent_agents":0},
				"items":[
					{"date":"2026-09-01","max_concurrent_agents":0},
					{"date":"2026-10-01","max_concurrent_agents":17}
				]
			}`,
			scope:      "organization_and_billing_children",
			interval:   "month",
			agentType:  "hosted",
			percentile: 95,
			summary:    ConcurrentAgentsUsageSummary{MaxConcurrentAgents: 17, PercentileConcurrentAgents: 0},
			items: []ConcurrentAgentsUsageItem{
				{Date: "2026-09-01", MaxConcurrentAgents: 0},
				{Date: "2026-10-01", MaxConcurrentAgents: 17},
			},
		},
		{
			name: "zero_daily_uuid",
			org:  "7c202aaa-3165-4811-9813-173c4c285463",
			path: "/v2/organizations/7c202aaa-3165-4811-9813-173c4c285463/usage/concurrent-agents",
			opt: ConcurrentAgentsUsageOptions{
				From: "2026-10-25", To: "2026-10-25", Interval: "day", Percentile: 1, AgentType: "self_hosted",
			},
			query: values{"from": "2026-10-25", "to": "2026-10-25", "interval": "day", "percentile": "1", "filter[agent_type]": "self_hosted"},
			body: `{
				"description":"Recorded usage can change as late events arrive.",
				"reporting_period":{"from":"2026-10-25","to":"2026-10-25"},
				"generated_at":"2026-10-27T03:04:05Z",
				"scope":"organization","interval":"day","unit":"agents",
				"definition":"Discrete percentile across zero-filled daily sums of organization peaks.",
				"agent_type":"self_hosted","percentile":1,
				"summary":{"max_concurrent_agents":0,"percentile_concurrent_agents":0},
				"items":[{"date":"2026-10-25","max_concurrent_agents":0}]
			}`,
			scope:      "organization",
			interval:   "day",
			agentType:  "self_hosted",
			percentile: 1,
			summary:    ConcurrentAgentsUsageSummary{},
			items:      []ConcurrentAgentsUsageItem{{Date: "2026-10-25"}},
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

			got, resp, err := client.Organizations.GetConcurrentAgentsUsage(context.Background(), tt.org, &tt.opt)
			if err != nil {
				t.Fatal(err)
			}
			want := ConcurrentAgentsUsage{
				Description:     "Recorded usage can change as late events arrive.",
				ReportingPeriod: UsageReportingPeriod{From: tt.opt.From, To: tt.opt.To},
				GeneratedAt:     NewTimestamp(time.Date(2026, 10, 27, 3, 4, 5, 0, time.UTC)),
				Scope:           tt.scope,
				Interval:        tt.interval,
				Unit:            "agents",
				Definition:      "Discrete percentile across zero-filled daily sums of organization peaks.",
				AgentType:       tt.agentType,
				Percentile:      tt.percentile,
				Summary:         tt.summary,
				Items:           tt.items,
			}
			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("concurrent-agent usage (-want +got):\n%s", diff)
			}
			if resp == nil || resp.StatusCode != http.StatusOK {
				t.Fatalf("response = %+v; want status 200", resp)
			}
			data, err := json.Marshal(got)
			if err != nil {
				t.Fatal(err)
			}
			assertJSONEqual(t, string(data), tt.body)
		})
	}
}

func TestOrganizationsService_GetConcurrentAgentsUsageErrors(t *testing.T) {
	t.Parallel()
	for _, status := range []int{http.StatusBadRequest, http.StatusUnauthorized, http.StatusForbidden, http.StatusNotFound, http.StatusServiceUnavailable} {
		t.Run(fmt.Sprint(status), func(t *testing.T) {
			t.Parallel()
			server, client, teardown := newMockServerAndClient(t)
			t.Cleanup(teardown)
			body := fmt.Sprintf(`{"message":%q}`, http.StatusText(status))
			server.HandleFunc("/v2/organizations/acme/usage/concurrent-agents", func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "GET")
				testFormValues(t, r, values{"from": "2026-10", "to": "2026-10-25", "interval": "week", "percentile": "101", "filter[agent_type]": "hosted & self_hosted"})
				w.WriteHeader(status)
				_, _ = fmt.Fprint(w, body)
			})
			got, resp, err := client.Organizations.GetConcurrentAgentsUsage(context.Background(), "acme", &ConcurrentAgentsUsageOptions{
				From: "2026-10", To: "2026-10-25", Interval: "week", Percentile: 101, AgentType: "hosted & self_hosted",
			})
			var apiErr *ErrorResponse
			if !errors.As(err, &apiErr) {
				t.Fatalf("error = %v; want *ErrorResponse", err)
			}
			if resp == nil || resp.StatusCode != status || apiErr.Message != http.StatusText(status) || string(apiErr.RawBody) != body {
				t.Errorf("response=%+v error=%+v", resp, apiErr)
			}
			if diff := cmp.Diff(ConcurrentAgentsUsage{}, got); diff != "" {
				t.Errorf("report returned on error (-want +got):\n%s", diff)
			}
		})
	}

	t.Run("missing_options", func(t *testing.T) {
		server, client, teardown := newMockServerAndClient(t)
		t.Cleanup(teardown)
		server.HandleFunc("/v2/organizations/acme/usage/concurrent-agents", func(w http.ResponseWriter, r *http.Request) {
			testFormValues(t, r, values{})
			w.WriteHeader(http.StatusBadRequest)
			_, _ = fmt.Fprint(w, `{"message":"from and to are required"}`)
		})
		_, resp, err := client.Organizations.GetConcurrentAgentsUsage(context.Background(), "acme", nil)
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
		_, _, err := client.Organizations.GetConcurrentAgentsUsage(ctx, "acme", &ConcurrentAgentsUsageOptions{
			From: "2026-10-23", To: "2026-10-25",
		})
		if !errors.Is(err, context.Canceled) {
			t.Errorf("error = %v; want context.Canceled", err)
		}
	})
}

func TestOrganizationsService_GetTestExecutionsUsage(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		org      string
		path     string
		opt      TestExecutionsUsageOptions
		query    values
		body     string
		scope    string
		interval string
		summary  TestExecutionsUsageSummary
		items    []TestExecutionsUsageItem
	}{
		{
			name:  "default_daily_int64",
			org:   "acme",
			path:  "/v2/organizations/acme/usage/test-executions",
			opt:   TestExecutionsUsageOptions{From: "2026-10-23", To: "2026-10-25"},
			query: values{"from": "2026-10-23", "to": "2026-10-25"},
			body: `{
				"description":"Recorded executions can change as late events arrive.",
				"reporting_period":{"from":"2026-10-23","to":"2026-10-25"},
				"generated_at":"2026-10-27T03:04:05Z",
				"scope":"organization_and_billing_children","interval":"day","unit":"executions",
				"definition":"Recorded test executions including merged partial and correction rows, not managed tests.",
				"summary":{"test_executions":2147483651,"active_suites":2,"average_test_executions":715827883},
				"items":[
					{"date":"2026-10-23","test_executions":2147483647,"cumulative_test_executions":2147483647},
					{"date":"2026-10-24","test_executions":0,"cumulative_test_executions":2147483647},
					{"date":"2026-10-25","test_executions":4,"cumulative_test_executions":2147483651}
				]
			}`,
			scope:    "organization_and_billing_children",
			interval: "day",
			summary:  TestExecutionsUsageSummary{TestExecutions: 2147483651, ActiveSuites: 2, AverageTestExecutions: 715827883},
			items: []TestExecutionsUsageItem{
				{Date: "2026-10-23", TestExecutions: 2147483647, CumulativeTestExecutions: 2147483647},
				{Date: "2026-10-24", TestExecutions: 0, CumulativeTestExecutions: 2147483647},
				{Date: "2026-10-25", TestExecutions: 4, CumulativeTestExecutions: 2147483651},
			},
		},
		{
			name:  "monthly_partial_range",
			org:   "acme",
			path:  "/v2/organizations/acme/usage/test-executions",
			opt:   TestExecutionsUsageOptions{From: "2026-08-15", To: "2026-10-10", Interval: "month"},
			query: values{"from": "2026-08-15", "to": "2026-10-10", "interval": "month"},
			body: `{
				"description":"Recorded executions can change as late events arrive.",
				"reporting_period":{"from":"2026-08-15","to":"2026-10-10"},
				"generated_at":"2026-10-27T03:04:05Z",
				"scope":"organization_and_billing_children","interval":"month","unit":"executions",
				"definition":"Recorded test executions including merged partial and correction rows, not managed tests.",
				"summary":{"test_executions":11,"active_suites":2,"average_test_executions":3},
				"items":[
					{"date":"2026-08-01","test_executions":4,"cumulative_test_executions":4},
					{"date":"2026-09-01","test_executions":0,"cumulative_test_executions":4},
					{"date":"2026-10-01","test_executions":7,"cumulative_test_executions":11}
				]
			}`,
			scope:    "organization_and_billing_children",
			interval: "month",
			summary:  TestExecutionsUsageSummary{TestExecutions: 11, ActiveSuites: 2, AverageTestExecutions: 3},
			items: []TestExecutionsUsageItem{
				{Date: "2026-08-01", TestExecutions: 4, CumulativeTestExecutions: 4},
				{Date: "2026-09-01", TestExecutions: 0, CumulativeTestExecutions: 4},
				{Date: "2026-10-01", TestExecutions: 7, CumulativeTestExecutions: 11},
			},
		},
		{
			name:  "zero_explicit_daily_uuid",
			org:   "7c202aaa-3165-4811-9813-173c4c285463",
			path:  "/v2/organizations/7c202aaa-3165-4811-9813-173c4c285463/usage/test-executions",
			opt:   TestExecutionsUsageOptions{From: "2026-10-25", To: "2026-10-25", Interval: "day"},
			query: values{"from": "2026-10-25", "to": "2026-10-25", "interval": "day"},
			body: `{
				"description":"Recorded executions can change as late events arrive.",
				"reporting_period":{"from":"2026-10-25","to":"2026-10-25"},
				"generated_at":"2026-10-27T03:04:05Z",
				"scope":"organization","interval":"day","unit":"executions",
				"definition":"Recorded test executions including merged partial and correction rows, not managed tests.",
				"summary":{"test_executions":0,"active_suites":0,"average_test_executions":0},
				"items":[{"date":"2026-10-25","test_executions":0,"cumulative_test_executions":0}]
			}`,
			scope:    "organization",
			interval: "day",
			summary:  TestExecutionsUsageSummary{},
			items:    []TestExecutionsUsageItem{{Date: "2026-10-25"}},
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
			got, resp, err := client.Organizations.GetTestExecutionsUsage(context.Background(), tt.org, &tt.opt)
			if err != nil {
				t.Fatal(err)
			}
			want := TestExecutionsUsage{
				Description:     "Recorded executions can change as late events arrive.",
				ReportingPeriod: UsageReportingPeriod{From: tt.opt.From, To: tt.opt.To},
				GeneratedAt:     NewTimestamp(time.Date(2026, 10, 27, 3, 4, 5, 0, time.UTC)),
				Scope:           tt.scope,
				Interval:        tt.interval,
				Unit:            "executions",
				Definition:      "Recorded test executions including merged partial and correction rows, not managed tests.",
				Summary:         tt.summary,
				Items:           tt.items,
			}
			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("test-execution usage (-want +got):\n%s", diff)
			}
			if resp == nil || resp.StatusCode != http.StatusOK {
				t.Fatalf("response = %+v; want status 200", resp)
			}
			data, err := json.Marshal(got)
			if err != nil {
				t.Fatal(err)
			}
			assertJSONEqual(t, string(data), tt.body)
		})
	}
}

func TestOrganizationsService_GetTestExecutionsUsageErrors(t *testing.T) {
	t.Parallel()
	for _, status := range []int{http.StatusBadRequest, http.StatusServiceUnavailable} {
		t.Run(fmt.Sprint(status), func(t *testing.T) {
			t.Parallel()
			server, client, teardown := newMockServerAndClient(t)
			t.Cleanup(teardown)
			body := fmt.Sprintf(`{"message":%q}`, http.StatusText(status))
			server.HandleFunc("/v2/organizations/acme/usage/test-executions", func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "GET")
				testFormValues(t, r, values{})
				w.WriteHeader(status)
				_, _ = fmt.Fprint(w, body)
			})
			got, resp, err := client.Organizations.GetTestExecutionsUsage(context.Background(), "acme", nil)
			var apiErr *ErrorResponse
			if !errors.As(err, &apiErr) || resp == nil || resp.StatusCode != status {
				t.Fatalf("response=%+v error=%v; want API error %d", resp, err, status)
			}
			if apiErr.Message != http.StatusText(status) || string(apiErr.RawBody) != body {
				t.Errorf("error=%+v", apiErr)
			}
			if diff := cmp.Diff(TestExecutionsUsage{}, got); diff != "" {
				t.Errorf("report returned on error (-want +got):\n%s", diff)
			}
		})
	}
}

func TestOrganizationsService_GetConcurrentAgentsUsageByOrganization(t *testing.T) {
	t.Parallel()
	next := "https://api.buildkite.com/v2/organizations/acme/usage/concurrent-agents/organizations?from=2026-09-21&to=2026-10-11&interval=month&filter%5Bagent_type%5D=hosted&percentile=95&per_page=2&after=next-cursor"
	tests := []struct {
		name       string
		org        string
		path       string
		opt        ConcurrentAgentsUsageByOrganizationOptions
		query      values
		body       string
		scope      string
		interval   string
		agentType  string
		items      []ConcurrentAgentsUsageOrganization
		links      UsageReportLinks
		linkHeader string
	}{
		{
			name: "monthly_cursor_page",
			org:  "acme",
			path: "/v2/organizations/acme/usage/concurrent-agents/organizations",
			opt: ConcurrentAgentsUsageByOrganizationOptions{
				ConcurrentAgentsUsageOptions: ConcurrentAgentsUsageOptions{
					From: "2026-09-21", To: "2026-10-11", Interval: "month", AgentType: "hosted", Percentile: 95,
				},
				PerPage: 2, After: "opaque+/= &?",
			},
			query: values{"from": "2026-09-21", "to": "2026-10-11", "interval": "month", "filter[agent_type]": "hosted", "percentile": "95", "per_page": "2", "after": "opaque+/= &?"},
			body: `{
				"description":"Organization usage can change between pages.",
				"reporting_period":{"from":"2026-09-21","to":"2026-10-11"},
				"generated_at":"2026-10-27T03:04:05Z",
				"scope":"organization_and_billing_children","interval":"month","unit":"agents",
				"definition":"Per-organization daily peaks; percentiles are not additive.",
				"agent_type":"hosted","percentile":95,
				"items":[{
					"organization_uuid":"org-1","organization_slug":"acme","organization_name":"Acme",
					"summary":{"max_concurrent_agents":17,"percentile_concurrent_agents":0},
					"items":[
						{"date":"2026-09-01","max_concurrent_agents":0,"percentile_concurrent_agents":0},
						{"date":"2026-10-01","max_concurrent_agents":17,"percentile_concurrent_agents":17}
					]
				},{
					"organization_uuid":"org-2","organization_slug":"acme-child","organization_name":"Child",
					"summary":{"max_concurrent_agents":0,"percentile_concurrent_agents":0},
					"items":[
						{"date":"2026-09-01","max_concurrent_agents":0,"percentile_concurrent_agents":0},
						{"date":"2026-10-01","max_concurrent_agents":0,"percentile_concurrent_agents":0}
					]
				}],
				"links":{
					"self":"https://api.buildkite.com/v2/organizations/acme/usage/concurrent-agents/organizations?from=2026-09-21&to=2026-10-11&interval=month&filter%5Bagent_type%5D=hosted&percentile=95&per_page=2&after=opaque%2B%2F%3D+%26%3F",
					"next":"https://api.buildkite.com/v2/organizations/acme/usage/concurrent-agents/organizations?from=2026-09-21&to=2026-10-11&interval=month&filter%5Bagent_type%5D=hosted&percentile=95&per_page=2&after=next-cursor"
				}
			}`,
			scope: "organization_and_billing_children", interval: "month", agentType: "hosted",
			items: []ConcurrentAgentsUsageOrganization{
				{
					OrganizationUUID: "org-1", OrganizationSlug: "acme", OrganizationName: "Acme",
					Summary: ConcurrentAgentsUsageSummary{MaxConcurrentAgents: 17, PercentileConcurrentAgents: 0},
					Items: []ConcurrentAgentsUsageOrganizationItem{
						{Date: "2026-09-01"},
						{Date: "2026-10-01", MaxConcurrentAgents: 17, PercentileConcurrentAgents: 17},
					},
				},
				{
					OrganizationUUID: "org-2", OrganizationSlug: "acme-child", OrganizationName: "Child",
					Summary: ConcurrentAgentsUsageSummary{},
					Items: []ConcurrentAgentsUsageOrganizationItem{
						{Date: "2026-09-01"}, {Date: "2026-10-01"},
					},
				},
			},
			links: UsageReportLinks{
				Self: "https://api.buildkite.com/v2/organizations/acme/usage/concurrent-agents/organizations?from=2026-09-21&to=2026-10-11&interval=month&filter%5Bagent_type%5D=hosted&percentile=95&per_page=2&after=opaque%2B%2F%3D+%26%3F",
				Next: &next,
			},
			linkHeader: "<" + next + `>; rel="next"`,
		},
		{
			name: "daily_final_uuid_defaults",
			org:  "7c202aaa-3165-4811-9813-173c4c285463",
			path: "/v2/organizations/7c202aaa-3165-4811-9813-173c4c285463/usage/concurrent-agents/organizations",
			opt: ConcurrentAgentsUsageByOrganizationOptions{ConcurrentAgentsUsageOptions: ConcurrentAgentsUsageOptions{
				From: "2026-10-25", To: "2026-10-25",
			}},
			query: values{"from": "2026-10-25", "to": "2026-10-25"},
			body: `{
				"description":"Organization usage can change between pages.",
				"reporting_period":{"from":"2026-10-25","to":"2026-10-25"},
				"generated_at":"2026-10-27T03:04:05Z",
				"scope":"organization","interval":"day","unit":"agents",
				"definition":"Per-organization daily peaks; percentiles are not additive.",
				"agent_type":"self_hosted","percentile":95,
				"items":[{
					"organization_uuid":"org-1","organization_slug":"acme","organization_name":"Acme",
					"summary":{"max_concurrent_agents":9,"percentile_concurrent_agents":9},
					"items":[{"date":"2026-10-25","max_concurrent_agents":9,"percentile_concurrent_agents":9}]
				}],
				"links":{"self":"https://api.buildkite.com/v2/organizations/7c202aaa-3165-4811-9813-173c4c285463/usage/concurrent-agents/organizations?from=2026-10-25&to=2026-10-25","next":null}
			}`,
			scope: "organization", interval: "day", agentType: "self_hosted",
			items: []ConcurrentAgentsUsageOrganization{{
				OrganizationUUID: "org-1", OrganizationSlug: "acme", OrganizationName: "Acme",
				Summary: ConcurrentAgentsUsageSummary{MaxConcurrentAgents: 9, PercentileConcurrentAgents: 9},
				Items:   []ConcurrentAgentsUsageOrganizationItem{{Date: "2026-10-25", MaxConcurrentAgents: 9, PercentileConcurrentAgents: 9}},
			}},
			links: UsageReportLinks{Self: "https://api.buildkite.com/v2/organizations/7c202aaa-3165-4811-9813-173c4c285463/usage/concurrent-agents/organizations?from=2026-10-25&to=2026-10-25"},
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
				if tt.linkHeader != "" {
					w.Header().Set("Link", tt.linkHeader)
				}
				_, _ = fmt.Fprint(w, tt.body)
			})
			got, resp, err := client.Organizations.GetConcurrentAgentsUsageByOrganization(context.Background(), tt.org, &tt.opt)
			if err != nil {
				t.Fatal(err)
			}
			want := ConcurrentAgentsUsageByOrganization{
				Description:     "Organization usage can change between pages.",
				ReportingPeriod: UsageReportingPeriod{From: tt.opt.From, To: tt.opt.To},
				GeneratedAt:     NewTimestamp(time.Date(2026, 10, 27, 3, 4, 5, 0, time.UTC)),
				Scope:           tt.scope,
				Interval:        tt.interval,
				Unit:            "agents",
				Definition:      "Per-organization daily peaks; percentiles are not additive.",
				AgentType:       tt.agentType,
				Percentile:      95,
				Items:           tt.items,
				Links:           tt.links,
			}
			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("concurrency by organization (-want +got):\n%s", diff)
			}
			if resp == nil || resp.StatusCode != http.StatusOK || resp.Header.Get("Link") != tt.linkHeader {
				t.Fatalf("response = %+v; want status 200 and Link %q", resp, tt.linkHeader)
			}
			data, err := json.Marshal(got)
			if err != nil {
				t.Fatal(err)
			}
			assertJSONEqual(t, string(data), tt.body)
		})
	}
}

func TestOrganizationsService_GetConcurrentAgentsUsageByOrganizationErrors(t *testing.T) {
	t.Parallel()
	for _, status := range []int{http.StatusBadRequest, http.StatusServiceUnavailable} {
		t.Run(fmt.Sprint(status), func(t *testing.T) {
			t.Parallel()
			server, client, teardown := newMockServerAndClient(t)
			t.Cleanup(teardown)
			body := fmt.Sprintf(`{"message":%q}`, http.StatusText(status))
			server.HandleFunc("/v2/organizations/acme/usage/concurrent-agents/organizations", func(w http.ResponseWriter, r *http.Request) {
				testFormValues(t, r, values{})
				w.WriteHeader(status)
				_, _ = fmt.Fprint(w, body)
			})
			got, resp, err := client.Organizations.GetConcurrentAgentsUsageByOrganization(context.Background(), "acme", nil)
			var apiErr *ErrorResponse
			if !errors.As(err, &apiErr) || resp == nil || resp.StatusCode != status {
				t.Fatalf("response=%+v error=%v; want API error %d", resp, err, status)
			}
			if apiErr.Message != http.StatusText(status) || string(apiErr.RawBody) != body {
				t.Errorf("error=%+v", apiErr)
			}
			if diff := cmp.Diff(ConcurrentAgentsUsageByOrganization{}, got); diff != "" {
				t.Errorf("report returned on error (-want +got):\n%s", diff)
			}
		})
	}
}

func TestOrganizationsService_GetHostedAgentsUsageByPipeline(t *testing.T) {
	t.Parallel()
	pipelineUUID, pipelineSlug := "pipeline-1", "build"
	next := "https://api.buildkite.com/v2/organizations/acme/usage/hosted-agents/pipelines?from=2026-09-21&to=2026-10-11&interval=month&per_page=2&after=next-cursor"
	tests := []struct {
		name       string
		org        string
		opt        HostedAgentsUsageByPipelineOptions
		query      values
		body       string
		items      []HostedAgentsUsagePipeline
		linkHeader string
	}{
		{
			name: "monthly_sparse_cursor_page",
			org:  "acme",
			opt: HostedAgentsUsageByPipelineOptions{
				HostedAgentsUsageOptions: HostedAgentsUsageOptions{From: "2026-09-21", To: "2026-10-11", Interval: "month"},
				PerPage:                  2, After: "opaque+/= &?",
			},
			query: values{"from": "2026-09-21", "to": "2026-10-11", "interval": "month", "per_page": "2", "after": "opaque+/= &?"},
			body: `{
				"description":"CSV rows can change as late events arrive.",
				"reporting_period":{"from":"2026-09-21","to":"2026-10-11"},
				"generated_at":"2026-10-27T03:04:05Z",
				"scope":"organization_and_billing_children","interval":"month","unit":"vCPU-minutes",
				"definition":"CPU-weighted minutes floored per CSV row, not overview totals.",
				"items":[
					{"date":"2026-09-01","organization_uuid":"org-1","pipeline_id":4294967301,"pipeline_uuid":"pipeline-1","pipeline_slug":"build","instance_shape":"linux-amd64-2x4","duration_microseconds":59999999,"job_count":3,"vcpu_minutes":1},
					{"date":"2026-10-01","organization_uuid":"org-2","pipeline_id":4294967302,"pipeline_uuid":null,"pipeline_slug":null,"instance_shape":"future-shape","duration_microseconds":9007199254740993,"job_count":4294967303,"vcpu_minutes":0}
				],
				"links":{
					"self":"https://api.buildkite.com/v2/organizations/acme/usage/hosted-agents/pipelines?from=2026-09-21&to=2026-10-11&interval=month&per_page=2&after=opaque%2B%2F%3D+%26%3F",
					"next":"https://api.buildkite.com/v2/organizations/acme/usage/hosted-agents/pipelines?from=2026-09-21&to=2026-10-11&interval=month&per_page=2&after=next-cursor"
				}
			}`,
			items: []HostedAgentsUsagePipeline{
				{Date: "2026-09-01", OrganizationUUID: "org-1", PipelineID: 4294967301, PipelineUUID: &pipelineUUID, PipelineSlug: &pipelineSlug, InstanceShape: "linux-amd64-2x4", DurationMicroseconds: 59999999, JobCount: 3, VCPUMinutes: 1},
				{Date: "2026-10-01", OrganizationUUID: "org-2", PipelineID: 4294967302, InstanceShape: "future-shape", DurationMicroseconds: 9007199254740993, JobCount: 4294967303},
			},
			linkHeader: "<" + next + `>; rel="next"`,
		},
		{
			name: "empty_final_uuid_defaults",
			org:  "7c202aaa-3165-4811-9813-173c4c285463",
			opt: HostedAgentsUsageByPipelineOptions{HostedAgentsUsageOptions: HostedAgentsUsageOptions{
				From: "2026-10-25", To: "2026-10-25",
			}},
			query: values{"from": "2026-10-25", "to": "2026-10-25"},
			body: `{
				"description":"CSV rows can change as late events arrive.",
				"reporting_period":{"from":"2026-10-25","to":"2026-10-25"},
				"generated_at":"2026-10-27T03:04:05Z",
				"scope":"organization","interval":"day","unit":"vCPU-minutes",
				"definition":"CPU-weighted minutes floored per CSV row, not overview totals.",
				"items":[],
				"links":{"self":"https://api.buildkite.com/v2/organizations/7c202aaa-3165-4811-9813-173c4c285463/usage/hosted-agents/pipelines?from=2026-10-25&to=2026-10-25","next":null}
			}`,
			items: []HostedAgentsUsagePipeline{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server, client, teardown := newMockServerAndClient(t)
			t.Cleanup(teardown)
			server.HandleFunc("/v2/organizations/"+tt.org+"/usage/hosted-agents/pipelines", func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "GET")
				testFormValues(t, r, tt.query)
				if tt.linkHeader != "" {
					w.Header().Set("Link", tt.linkHeader)
				}
				_, _ = fmt.Fprint(w, tt.body)
			})
			got, resp, err := client.Organizations.GetHostedAgentsUsageByPipeline(context.Background(), tt.org, &tt.opt)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(tt.items, got.Items); diff != "" {
				t.Errorf("pipeline rows (-want +got):\n%s", diff)
			}
			if resp == nil || resp.StatusCode != http.StatusOK || resp.Header.Get("Link") != tt.linkHeader {
				t.Fatalf("response = %+v; want status 200 and Link %q", resp, tt.linkHeader)
			}
			data, err := json.Marshal(got)
			if err != nil {
				t.Fatal(err)
			}
			assertJSONEqual(t, string(data), tt.body)
		})
	}
}

func TestOrganizationsService_GetHostedAgentsUsageByPipelineErrors(t *testing.T) {
	t.Parallel()
	for _, status := range []int{http.StatusBadRequest, http.StatusNotFound, http.StatusServiceUnavailable} {
		t.Run(fmt.Sprint(status), func(t *testing.T) {
			t.Parallel()
			server, client, teardown := newMockServerAndClient(t)
			t.Cleanup(teardown)
			body := fmt.Sprintf(`{"message":%q}`, http.StatusText(status))
			server.HandleFunc("/v2/organizations/acme/usage/hosted-agents/pipelines", func(w http.ResponseWriter, r *http.Request) {
				testFormValues(t, r, values{})
				w.WriteHeader(status)
				_, _ = fmt.Fprint(w, body)
			})
			got, resp, err := client.Organizations.GetHostedAgentsUsageByPipeline(context.Background(), "acme", nil)
			var apiErr *ErrorResponse
			if !errors.As(err, &apiErr) || resp == nil || resp.StatusCode != status {
				t.Fatalf("response=%+v error=%v; want API error %d", resp, err, status)
			}
			if apiErr.Message != http.StatusText(status) || string(apiErr.RawBody) != body {
				t.Errorf("error=%+v", apiErr)
			}
			if diff := cmp.Diff(HostedAgentsUsageByPipeline{}, got); diff != "" {
				t.Errorf("report returned on error (-want +got):\n%s", diff)
			}
		})
	}
}
