package buildkite

import (
	"context"
	"fmt"
	"time"
)

// SpanAggregatesService handles communication with the span aggregate methods
// of the Buildkite Test Engine API.
type SpanAggregatesService struct {
	client *Client
}

// SpanAggregateScope contains totals for an OpenTelemetry instrumentation scope.
type SpanAggregateScope struct {
	Name      string `json:"name"`
	SpanCount int    `json:"span_count"`

	// Duration is the sum of span durations, not wall-clock time.
	Duration time.Duration `json:"duration_ns"`
}

// SpanAggregateName contains duration statistics for a span name within a scope.
type SpanAggregateName struct {
	Name      string        `json:"name"`
	SpanCount int           `json:"span_count"`
	Duration  time.Duration `json:"duration_ns"`

	// AverageDuration, P50Duration, and P95Duration are in nanoseconds and may
	// have fractional values. The percentiles are approximate.
	AverageDuration float64       `json:"avg_duration_ns"`
	P50Duration     float64       `json:"p50_duration_ns"`
	P95Duration     float64       `json:"p95_duration_ns"`
	MaximumDuration time.Duration `json:"max_duration_ns"`
}

// SpanAggregateNames contains the top span names for the requested scope.
type SpanAggregateNames struct {
	ScopeName          string              `json:"scope_name"`
	SpanNames          []SpanAggregateName `json:"span_names"`
	SpanNamesTruncated bool                `json:"span_names_truncated"`
}

// SpanAggregateScopesListOptions specifies optional parameters for
// [SpanAggregatesService.ListScopesBySuite].
type SpanAggregateScopesListOptions struct {
	// Period is one of "1hour", "4hours", "1day", "7days", "14days", or
	// "28days". Values exceeding the organization's API history allowance
	// return an API error. When omitted, the API selects the default period.
	Period string `url:"period,omitempty"`
}

// SpanAggregateNamesListOptions specifies optional parameters for
// [SpanAggregatesService.ListNamesByBuild] and is embedded in
// [SpanAggregateSuiteNamesListOptions]. Invalid values are validated by the API.
type SpanAggregateNamesListOptions struct {
	// Sort is "count", "total", "avg", "p50", "p95", or "max", optionally
	// prefixed with "-" for descending order. It defaults to "-total".
	// Ties are sorted by name ascending.
	Sort string `url:"sort,omitempty"`

	// Limit is the maximum number of span names to return, from 1 to 100.
	// It defaults to 20. SpanNamesTruncated indicates whether more names exist.
	Limit int `url:"limit,omitempty"`
}

// SpanAggregateSuiteNamesListOptions specifies optional parameters for
// [SpanAggregatesService.ListNamesBySuite]. Invalid values are validated by the API.
type SpanAggregateSuiteNamesListOptions struct {
	SpanAggregateNamesListOptions

	// Period is one of "1hour", "4hours", "1day", "7days", "14days", or
	// "28days". Values exceeding the organization's API history allowance
	// return an API error. When omitted, the API selects the default period.
	Period string `url:"period,omitempty"`
}

// ListScopesBySuite returns all instrumentation scopes for a suite in the
// requested period, ordered by total duration descending, then name ascending.
// The suite may be identified by slug or UUID. Results are not paginated.
// Aggregates have hourly precision and are retained for one month.
// OTel trace ingestion must be enabled for the organization.
func (sas *SpanAggregatesService) ListScopesBySuite(ctx context.Context, org, slug string, opt *SpanAggregateScopesListOptions) ([]SpanAggregateScope, *Response, error) {
	u := fmt.Sprintf("v2/analytics/organizations/%s/suites/%s/span-aggregates/scopes", org, slug)
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := sas.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var scopes []SpanAggregateScope
	resp, err := sas.client.Do(req, &scopes)
	if err != nil {
		return nil, resp, err
	}

	return scopes, resp, err
}

// ListScopesByBuild returns all instrumentation scopes across the build's
// suites that the caller can view, ordered by total duration descending, then
// name ascending. The caller must have access to the build's pipeline.
// buildUUID is the build's UUID, not its pipeline build number. Results are
// not paginated. OTel trace ingestion must be enabled for the organization.
func (sas *SpanAggregatesService) ListScopesByBuild(ctx context.Context, org, buildUUID string) ([]SpanAggregateScope, *Response, error) {
	u := fmt.Sprintf("v2/analytics/organizations/%s/builds/%s/span-aggregates/scopes", org, buildUUID)
	req, err := sas.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var scopes []SpanAggregateScope
	resp, err := sas.client.Do(req, &scopes)
	if err != nil {
		return nil, resp, err
	}

	return scopes, resp, err
}

// ListNamesBySuite returns the top span-name statistics for an instrumentation
// scope in a suite during the requested period. The suite may be identified by
// slug or UUID. An empty scopeName selects the unnamed scope; an unknown scope
// returns an empty span-name list. Results are not paginated. Aggregates have
// hourly precision and are retained for one month. OTel trace ingestion must
// be enabled for the organization.
func (sas *SpanAggregatesService) ListNamesBySuite(ctx context.Context, org, slug, scopeName string, opt *SpanAggregateSuiteNamesListOptions) (SpanAggregateNames, *Response, error) {
	u := fmt.Sprintf("v2/analytics/organizations/%s/suites/%s/span-aggregates/span-names", org, slug)
	params := struct {
		ScopeName string `url:"scope_name"`
		SpanAggregateSuiteNamesListOptions
	}{ScopeName: scopeName}
	if opt != nil {
		params.SpanAggregateSuiteNamesListOptions = *opt
	}
	u, err := addOptions(u, params)
	if err != nil {
		return SpanAggregateNames{}, nil, err
	}

	req, err := sas.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return SpanAggregateNames{}, nil, err
	}

	var names SpanAggregateNames
	resp, err := sas.client.Do(req, &names)
	if err != nil {
		return SpanAggregateNames{}, resp, err
	}

	return names, resp, err
}

// ListNamesByBuild returns the top span-name statistics for an instrumentation
// scope across the build's suites that the caller can view. The caller must
// have access to the build's pipeline. buildUUID is the build's UUID, not its
// pipeline build number. An empty scopeName selects the unnamed scope; an
// unknown scope returns an empty span-name list. Results are not paginated.
func (sas *SpanAggregatesService) ListNamesByBuild(ctx context.Context, org, buildUUID, scopeName string, opt *SpanAggregateNamesListOptions) (SpanAggregateNames, *Response, error) {
	u := fmt.Sprintf("v2/analytics/organizations/%s/builds/%s/span-aggregates/span-names", org, buildUUID)
	params := struct {
		ScopeName string `url:"scope_name"`
		SpanAggregateNamesListOptions
	}{ScopeName: scopeName}
	if opt != nil {
		params.SpanAggregateNamesListOptions = *opt
	}
	u, err := addOptions(u, params)
	if err != nil {
		return SpanAggregateNames{}, nil, err
	}

	req, err := sas.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return SpanAggregateNames{}, nil, err
	}

	var names SpanAggregateNames
	resp, err := sas.client.Do(req, &names)
	if err != nil {
		return SpanAggregateNames{}, resp, err
	}

	return names, resp, err
}
