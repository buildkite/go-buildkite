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

// HostedAgentsUsageOptions specifies the reporting period and interval for
// [OrganizationsService.GetHostedAgentsUsage]. The API validates these values.
type HostedAgentsUsageOptions struct {
	// From and To are required inclusive UTC dates in YYYY-MM-DD format.
	// From must be no earlier than three calendar months before today;
	// To must be no later than today and no earlier than From.
	From string `url:"from,omitempty"`
	To   string `url:"to,omitempty"`

	// Interval is "day" (the default) or "month".
	Interval string `url:"interval,omitempty"`
}

// HostedAgentsUsageSummary contains CPU-weighted usage, combining instance
// shapes by resource before flooring minutes, like the Usage table/overview.
// Unrecognized shapes are excluded. These totals can differ from sums of
// independently rounded interval values.
type HostedAgentsUsageSummary struct {
	VCPUMinutes int64 `json:"vcpu_minutes"`

	// VCPUMinutesByResource always contains all four canonical keys, including
	// zeros: hosted_agents_linux_amd64_vcpu_minutes,
	// hosted_agents_linux_arm64_vcpu_minutes, hosted_agents_macos_m4_vcpu_minutes,
	// and hosted_agents_windows_amd64_vcpu_minutes.
	VCPUMinutesByResource map[string]int64 `json:"vcpu_minutes_by_resource"`
}

// HostedAgentsUsageItem contains CPU-weighted minutes floored per resource and
// interval, like the Usage chart. It uses the same four resource keys as the
// summary, including zeros.
type HostedAgentsUsageItem struct {
	// Date is YYYY-MM-DD. Monthly buckets use the first of the calendar month,
	// even for partial-month requests; values honor the exact requested range.
	Date                  string           `json:"date"`
	VCPUMinutesByResource map[string]int64 `json:"vcpu_minutes_by_resource"`
}

// HostedAgentsUsage reports hosted compute using the Usage page's calculations.
// Late events can revise values on cache refresh.
type HostedAgentsUsage struct {
	Description     string               `json:"description"`
	ReportingPeriod UsageReportingPeriod `json:"reporting_period"`

	// GeneratedAt is the calculation time, retained in cached reports,
	// not an ingestion watermark.
	GeneratedAt *Timestamp `json:"generated_at"`

	// Scope is "organization" when there are no billing children, otherwise
	// "organization_and_billing_children".
	Scope string `json:"scope"`

	Interval   string                   `json:"interval"`
	Unit       string                   `json:"unit"`
	Definition string                   `json:"definition"`
	Summary    HostedAgentsUsageSummary `json:"summary"`
	Items      []HostedAgentsUsageItem  `json:"items"`
}

// GetHostedAgentsUsage returns hosted compute usage for an organization (slug
// or UUID) and its billing children, without pagination. Requests require
// read_organizations and access to the requested organization. Invalid options
// (400), a disabled usage-breakdowns feature (404), and unavailable usage data
// (503) remain API errors rather than empty reports.
func (os *OrganizationsService) GetHostedAgentsUsage(ctx context.Context, org string, opt *HostedAgentsUsageOptions) (HostedAgentsUsage, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/usage/hosted-agents", org)
	u, err := addOptions(u, opt)
	if err != nil {
		return HostedAgentsUsage{}, nil, err
	}

	req, err := os.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return HostedAgentsUsage{}, nil, err
	}

	var usage HostedAgentsUsage
	resp, err := os.client.Do(req, &usage)
	if err != nil {
		return HostedAgentsUsage{}, resp, err
	}

	return usage, resp, nil
}

// TestExecutionsUsageOptions specifies the reporting period and interval for
// [OrganizationsService.GetTestExecutionsUsage]. The API validates these values.
type TestExecutionsUsageOptions struct {
	// From and To are required inclusive UTC dates in YYYY-MM-DD format.
	// From must be no earlier than three calendar months before today;
	// To must be no later than today and no earlier than From.
	From string `url:"from,omitempty"`
	To   string `url:"to,omitempty"`

	// Interval is "day" (the default) or "month".
	Interval string `url:"interval,omitempty"`
}

// TestExecutionsUsageSummary contains the Usage page's recorded execution
// counts, including merged partial/correction rows, not managed-test usage.
type TestExecutionsUsageSummary struct {
	TestExecutions int64 `json:"test_executions"`
	ActiveSuites   int64 `json:"active_suites"`

	// AverageTestExecutions is the total divided by the number of requested
	// day/month buckets, floored like the Usage panels.
	AverageTestExecutions int64 `json:"average_test_executions"`
}

// TestExecutionsUsageItem contains recorded executions in an interval and the
// cumulative executions within the requested period through that interval.
type TestExecutionsUsageItem struct {
	// Date is YYYY-MM-DD. Monthly buckets use the first of the calendar month,
	// even for partial-month requests; counts honor the exact requested range.
	Date                     string `json:"date"`
	TestExecutions           int64  `json:"test_executions"`
	CumulativeTestExecutions int64  `json:"cumulative_test_executions"`
}

// TestExecutionsUsage reports recorded test executions using the Usage page's
// calculations. Late events can revise values on cache refresh.
type TestExecutionsUsage struct {
	Description     string               `json:"description"`
	ReportingPeriod UsageReportingPeriod `json:"reporting_period"`

	// GeneratedAt is the calculation time, retained in cached reports,
	// not an ingestion watermark.
	GeneratedAt *Timestamp `json:"generated_at"`

	// Scope is "organization" when there are no billing children, otherwise
	// "organization_and_billing_children".
	Scope string `json:"scope"`

	Interval   string                     `json:"interval"`
	Unit       string                     `json:"unit"`
	Definition string                     `json:"definition"`
	Summary    TestExecutionsUsageSummary `json:"summary"`

	// Items is not paginated. Missing intervals have zero executions and
	// carry forward the cumulative count.
	Items []TestExecutionsUsageItem `json:"items"`
}

// GetTestExecutionsUsage returns recorded test execution usage for an
// organization (slug or UUID) and its billing children. Requests require
// read_organizations and access to the requested organization. Invalid options
// (400), a disabled usage-breakdowns feature (404), and unavailable usage data
// (503) remain API errors rather than empty reports.
func (os *OrganizationsService) GetTestExecutionsUsage(ctx context.Context, org string, opt *TestExecutionsUsageOptions) (TestExecutionsUsage, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/usage/test-executions", org)
	u, err := addOptions(u, opt)
	if err != nil {
		return TestExecutionsUsage{}, nil, err
	}

	req, err := os.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return TestExecutionsUsage{}, nil, err
	}

	var usage TestExecutionsUsage
	resp, err := os.client.Do(req, &usage)
	if err != nil {
		return TestExecutionsUsage{}, resp, err
	}

	return usage, resp, nil
}

// UsageReportLinks contains cursor pagination URLs. Next is nil on the final
// page and is serialized as null, not omitted.
type UsageReportLinks struct {
	Self string  `json:"self"`
	Next *string `json:"next"`
}

// ConcurrentAgentsUsageByOrganizationOptions specifies the report parameters
// and forward-only organization cursor. The API validates all values.
type ConcurrentAgentsUsageByOrganizationOptions struct {
	ConcurrentAgentsUsageOptions

	// PerPage is from 1 to 100; zero omits it for the server default of 30.
	PerPage int `url:"per_page,omitempty"`

	// After is an opaque cursor. No snapshot is guaranteed between pages.
	After string `url:"after,omitempty"`
}

// ConcurrentAgentsUsageOrganizationItem contains an organization's maximum
// and discrete percentile across daily peaks in the clipped interval. On a
// daily report the percentile equals the maximum; monthly percentiles still
// use daily peaks, not instantaneous readings.
type ConcurrentAgentsUsageOrganizationItem struct {
	Date                       string `json:"date"`
	MaxConcurrentAgents        int64  `json:"max_concurrent_agents"`
	PercentileConcurrentAgents int64  `json:"percentile_concurrent_agents"`
}

// ConcurrentAgentsUsageOrganization contains an organization's zero-filled
// series and its summary across the entire requested period's daily peaks.
// Percentiles are not additive across organizations or intervals.
type ConcurrentAgentsUsageOrganization struct {
	OrganizationUUID string `json:"organization_uuid"`
	OrganizationSlug string `json:"organization_slug"`
	OrganizationName string `json:"organization_name"`

	Summary ConcurrentAgentsUsageSummary            `json:"summary"`
	Items   []ConcurrentAgentsUsageOrganizationItem `json:"items"`
}

// ConcurrentAgentsUsageByOrganization is a page of organizations ordered by
// UUID, each with its own recorded concurrent-agent usage. Late events can
// revise values; pages are not a consistent snapshot.
type ConcurrentAgentsUsageByOrganization struct {
	Description     string               `json:"description"`
	ReportingPeriod UsageReportingPeriod `json:"reporting_period"`

	// GeneratedAt is the calculation time, retained in cached reports,
	// not an ingestion watermark or cross-page snapshot identifier.
	GeneratedAt *Timestamp `json:"generated_at"`

	// Scope is "organization" when there are no billing children, otherwise
	// "organization_and_billing_children".
	Scope string `json:"scope"`

	Interval   string                              `json:"interval"`
	Unit       string                              `json:"unit"`
	Definition string                              `json:"definition"`
	AgentType  string                              `json:"agent_type"`
	Percentile int                                 `json:"percentile"`
	Items      []ConcurrentAgentsUsageOrganization `json:"items"`
	Links      UsageReportLinks                    `json:"links"`
}

// GetConcurrentAgentsUsageByOrganization returns one page of organizations
// (not dates), including zero-usage organizations. Org accepts a slug or UUID.
// Requests require read_organizations and access to the requested organization.
// Disabled usage breakdowns and other API failures remain errors. Use Links.Next
// for forward pagination; the standard Link header is also available on Response.
func (os *OrganizationsService) GetConcurrentAgentsUsageByOrganization(ctx context.Context, org string, opt *ConcurrentAgentsUsageByOrganizationOptions) (ConcurrentAgentsUsageByOrganization, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/usage/concurrent-agents/organizations", org)
	u, err := addOptions(u, opt)
	if err != nil {
		return ConcurrentAgentsUsageByOrganization{}, nil, err
	}

	req, err := os.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return ConcurrentAgentsUsageByOrganization{}, nil, err
	}

	var usage ConcurrentAgentsUsageByOrganization
	resp, err := os.client.Do(req, &usage)
	if err != nil {
		return ConcurrentAgentsUsageByOrganization{}, resp, err
	}

	return usage, resp, nil
}

// HostedAgentsUsageByPipelineOptions specifies the reporting period, interval,
// and forward-only grouped-row cursor. The API validates all values.
type HostedAgentsUsageByPipelineOptions struct {
	HostedAgentsUsageOptions

	// PerPage counts grouped rows, not pipelines or dates. Valid values are
	// 1 to 100; zero omits it for the server default of 30.
	PerPage int `url:"per_page,omitempty"`

	// After is an opaque cursor. No snapshot is guaranteed between pages.
	After string `url:"after,omitempty"`
}

// HostedAgentsUsagePipeline is a sparse date/organization/pipeline/shape row
// using the CSV's CPU weighting and flooring. Summing these rounded rows may
// differ from the overview. Unknown nonempty shapes retain raw duration and
// job counts with zero VCPUMinutes; self-hosted empty shapes are excluded.
type HostedAgentsUsagePipeline struct {
	Date                 string  `json:"date"`
	OrganizationUUID     string  `json:"organization_uuid"`
	PipelineID           int64   `json:"pipeline_id"`
	PipelineUUID         *string `json:"pipeline_uuid"`
	PipelineSlug         *string `json:"pipeline_slug"`
	InstanceShape        string  `json:"instance_shape"`
	DurationMicroseconds int64   `json:"duration_microseconds"`
	JobCount             int64   `json:"job_count"`
	VCPUMinutes          int64   `json:"vcpu_minutes"`
}

// HostedAgentsUsageByPipeline is one sparse page of CSV detail, ordered by date,
// legacy organization ID, pipeline ID, and instance shape. Pipeline UUID and
// slug are nil when unavailable or not visible to the caller; usage is retained.
// Late events can revise values; pages are not a consistent snapshot.
type HostedAgentsUsageByPipeline struct {
	Description     string               `json:"description"`
	ReportingPeriod UsageReportingPeriod `json:"reporting_period"`

	// GeneratedAt is the calculation time retained in cached reports,
	// not an ingestion watermark or cross-page snapshot identifier.
	GeneratedAt *Timestamp `json:"generated_at"`

	// Scope is "organization" when there are no billing children, otherwise
	// "organization_and_billing_children".
	Scope string `json:"scope"`

	Interval   string                      `json:"interval"`
	Unit       string                      `json:"unit"`
	Definition string                      `json:"definition"`
	Items      []HostedAgentsUsagePipeline `json:"items"`
	Links      UsageReportLinks            `json:"links"`
}

// GetHostedAgentsUsageByPipeline returns one page of grouped hosted-agent
// usage rows for an organization (slug or UUID) and its billing children.
// Requests require read_organizations and access to the requested organization.
// API failures remain errors. Use Links.Next for forward pagination; the
// standard Link header is also available on Response.
func (os *OrganizationsService) GetHostedAgentsUsageByPipeline(ctx context.Context, org string, opt *HostedAgentsUsageByPipelineOptions) (HostedAgentsUsageByPipeline, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/usage/hosted-agents/pipelines", org)
	u, err := addOptions(u, opt)
	if err != nil {
		return HostedAgentsUsageByPipeline{}, nil, err
	}

	req, err := os.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return HostedAgentsUsageByPipeline{}, nil, err
	}

	var usage HostedAgentsUsageByPipeline
	resp, err := os.client.Do(req, &usage)
	if err != nil {
		return HostedAgentsUsageByPipeline{}, resp, err
	}

	return usage, resp, nil
}

// TestExecutionsUsageBySuiteOptions specifies the reporting period, interval,
// and forward-only grouped-row cursor. The API validates all values.
type TestExecutionsUsageBySuiteOptions struct {
	TestExecutionsUsageOptions

	// PerPage counts grouped rows, not distinct suites. Valid values are
	// 1 to 100; zero omits it for the server default of 30.
	PerPage int `url:"per_page,omitempty"`

	// After is an opaque cursor. No snapshot is guaranteed between pages.
	After string `url:"after,omitempty"`
}

// TestExecutionsUsageSuite is a sparse date/organization/suite row containing
// corrected recorded executions for the clipped bucket, not managed tests.
// Counts are additive. Missing and fully corrected zero groups are omitted.
type TestExecutionsUsageSuite struct {
	Date             string  `json:"date"`
	OrganizationUUID string  `json:"organization_uuid"`
	SuiteUUID        string  `json:"suite_uuid"`
	SuiteSlug        *string `json:"suite_slug"`
	SuiteName        *string `json:"suite_name"`
	TestExecutions   int64   `json:"test_executions"`
}

// TestExecutionsUsageBySuite is one sparse page ordered by date, organization
// UUID, and suite UUID. Suite slug and name are nil when unavailable or not
// visible to the caller; UUID and usage are retained. Late events can revise
// values; pages are not a consistent snapshot.
type TestExecutionsUsageBySuite struct {
	Description     string               `json:"description"`
	ReportingPeriod UsageReportingPeriod `json:"reporting_period"`

	// GeneratedAt is the calculation time retained in cached reports,
	// not an ingestion watermark or cross-page snapshot identifier.
	GeneratedAt *Timestamp `json:"generated_at"`

	// Scope is "organization" when there are no billing children, otherwise
	// "organization_and_billing_children".
	Scope string `json:"scope"`

	Interval   string                     `json:"interval"`
	Unit       string                     `json:"unit"`
	Definition string                     `json:"definition"`
	Items      []TestExecutionsUsageSuite `json:"items"`
	Links      UsageReportLinks           `json:"links"`
}

// GetTestExecutionsUsageBySuite returns one page of grouped test-execution
// usage rows for an organization (slug or UUID) and its billing children.
// Requests require read_organizations and access to the requested organization.
// API failures remain errors. Use Links.Next for forward pagination; the
// standard Link header is also available on Response.
func (os *OrganizationsService) GetTestExecutionsUsageBySuite(ctx context.Context, org string, opt *TestExecutionsUsageBySuiteOptions) (TestExecutionsUsageBySuite, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/usage/test-executions/suites", org)
	u, err := addOptions(u, opt)
	if err != nil {
		return TestExecutionsUsageBySuite{}, nil, err
	}

	req, err := os.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return TestExecutionsUsageBySuite{}, nil, err
	}

	var usage TestExecutionsUsageBySuite
	resp, err := os.client.Do(req, &usage)
	if err != nil {
		return TestExecutionsUsageBySuite{}, resp, err
	}

	return usage, resp, nil
}
