package buildkite

import (
	"context"
	"fmt"
)

// UsageReportingPeriod is the inclusive UTC date range used for a usage report.
type UsageReportingPeriod struct {
	From string `json:"from"`
	To   string `json:"to"`
}

// ActiveUsersUsageOptions specifies the period and interval for
// [OrganizationsService.GetActiveUsersUsage]. The API validates these values.
type ActiveUsersUsageOptions struct {
	// From and To are required inclusive UTC dates in YYYY-MM-DD format.
	// From must be no earlier than three calendar months before today;
	// To must be no later than today and no earlier than From.
	From string `url:"from,omitempty"`
	To   string `url:"to,omitempty"`

	// Interval is "day" (the default) or "month".
	Interval string `url:"interval,omitempty"`
}

// ActiveUsersUsageSummary contains the Usage page's organization-wide counts,
// including billing children.
type ActiveUsersUsageSummary struct {
	// ActiveUsers deduplicates users across the entire requested period and
	// billing children. It is not the sum of the interval counts.
	ActiveUsers int64 `json:"active_users"`

	// AverageMonthlyActiveUsers is the Usage page's average monthly count,
	// which can differ from ActiveUsers for multi-month periods.
	AverageMonthlyActiveUsers int64 `json:"average_monthly_active_users"`
}

// ActiveUsersUsageItem contains distinct users for an interval and the
// cumulative distinct users within the requested period through that interval.
type ActiveUsersUsageItem struct {
	// Date is YYYY-MM-DD. Monthly buckets use the first of the calendar month,
	// even when From is later; counts still honor the exact requested range.
	Date                  string `json:"date"`
	ActiveUsers           int64  `json:"active_users"`
	CumulativeActiveUsers int64  `json:"cumulative_active_users"`
}

// ActiveUsersUsage reports distinct active users using the same calculations
// as the Usage page. Late events can revise these values.
type ActiveUsersUsage struct {
	Description     string               `json:"description"`
	ReportingPeriod UsageReportingPeriod `json:"reporting_period"`

	// GeneratedAt is the calculation time, not an ingestion watermark.
	// Cached reports retain this time; late events appear after cache refresh.
	GeneratedAt *Timestamp `json:"generated_at"`

	// Scope is "organization" when there are no billing children, otherwise
	// "organization_and_billing_children".
	Scope string `json:"scope"`

	Interval   string                  `json:"interval"`
	Unit       string                  `json:"unit"`
	Definition string                  `json:"definition"`
	Summary    ActiveUsersUsageSummary `json:"summary"`

	// Items is not paginated. Missing intervals have zero active users and
	// carry forward the cumulative count.
	Items []ActiveUsersUsageItem `json:"items"`
}

// GetActiveUsersUsage returns active-user usage for an organization and its
// billing children. The organization can be a slug or UUID. Requests require
// read_organizations and access to the requested organization.
// A disabled organization usage breakdowns feature
// returns 404; invalid options return 400 and unavailable usage data returns
// 503. These responses remain API errors, not empty reports.
func (os *OrganizationsService) GetActiveUsersUsage(ctx context.Context, org string, opt *ActiveUsersUsageOptions) (ActiveUsersUsage, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/usage/active-users", org)
	u, err := addOptions(u, opt)
	if err != nil {
		return ActiveUsersUsage{}, nil, err
	}

	req, err := os.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return ActiveUsersUsage{}, nil, err
	}

	var usage ActiveUsersUsage
	resp, err := os.client.Do(req, &usage)
	if err != nil {
		return ActiveUsersUsage{}, resp, err
	}

	return usage, resp, nil
}

// ConcurrentAgentsUsageOptions specifies the period, interval and filters for
// [OrganizationsService.GetConcurrentAgentsUsage]. The API validates values.
type ConcurrentAgentsUsageOptions struct {
	// From and To are required inclusive UTC dates in YYYY-MM-DD format.
	// From must be no earlier than three calendar months before today;
	// To must be no later than today and no earlier than From.
	From string `url:"from,omitempty"`
	To   string `url:"to,omitempty"`

	// Interval is "day" (the default) or "month".
	Interval string `url:"interval,omitempty"`

	// Percentile is an integer from 1 to 100. Zero omits the parameter,
	// allowing the server default of 95.
	Percentile int `url:"percentile,omitempty"`

	// AgentType is "self_hosted" (the default) or "hosted".
	AgentType string `url:"filter[agent_type],omitempty"`
}

// ConcurrentAgentsUsageSummary describes the requested period's daily sums of
// each scoped organization's recorded peak. These sums do not represent
// simultaneous concurrency across organizations.
type ConcurrentAgentsUsageSummary struct {
	MaxConcurrentAgents int64 `json:"max_concurrent_agents"`

	// PercentileConcurrentAgents is the discrete percentile across all
	// zero-filled daily sums, even when the report's interval is monthly.
	// It is not a percentile of instantaneous readings or monthly bars.
	PercentileConcurrentAgents int64 `json:"percentile_concurrent_agents"`
}

// ConcurrentAgentsUsageItem is a daily sum of scoped organizations' peaks, or
// the maximum of those daily sums in a calendar month.
type ConcurrentAgentsUsageItem struct {
	// Date is YYYY-MM-DD. Monthly buckets use the first of the calendar month,
	// even for partial-month requests; counts honor the exact requested range.
	Date                string `json:"date"`
	MaxConcurrentAgents int64  `json:"max_concurrent_agents"`
}

// ConcurrentAgentsUsage reports recorded concurrent-agent usage using the
// Usage page's calculations. Late events can revise values on cache refresh.
type ConcurrentAgentsUsage struct {
	Description     string               `json:"description"`
	ReportingPeriod UsageReportingPeriod `json:"reporting_period"`

	// GeneratedAt is the calculation time, retained in cached reports,
	// not an ingestion watermark.
	GeneratedAt *Timestamp `json:"generated_at"`

	// Scope is "organization" when there are no billing children, otherwise
	// "organization_and_billing_children".
	Scope string `json:"scope"`

	Interval   string                       `json:"interval"`
	Unit       string                       `json:"unit"`
	Definition string                       `json:"definition"`
	AgentType  string                       `json:"agent_type"`
	Percentile int                          `json:"percentile"`
	Summary    ConcurrentAgentsUsageSummary `json:"summary"`

	// Items is a bounded, zero-filled series without pagination. Request
	// Interval "day" to obtain the daily inputs used for the percentile.
	Items []ConcurrentAgentsUsageItem `json:"items"`
}

// GetConcurrentAgentsUsage returns recorded concurrent-agent usage for an
// organization (slug or UUID) and its billing children. Requests require
// read_organizations and access to the requested organization. Invalid options
// (400), a disabled usage-breakdowns feature (404), and unavailable usage data
// (503) remain API errors rather than empty reports.
func (os *OrganizationsService) GetConcurrentAgentsUsage(ctx context.Context, org string, opt *ConcurrentAgentsUsageOptions) (ConcurrentAgentsUsage, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/usage/concurrent-agents", org)
	u, err := addOptions(u, opt)
	if err != nil {
		return ConcurrentAgentsUsage{}, nil, err
	}

	req, err := os.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return ConcurrentAgentsUsage{}, nil, err
	}

	var usage ConcurrentAgentsUsage
	resp, err := os.client.Do(req, &usage)
	if err != nil {
		return ConcurrentAgentsUsage{}, resp, err
	}

	return usage, resp, nil
}
