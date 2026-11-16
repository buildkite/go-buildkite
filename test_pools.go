package buildkite

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

// TestPoolsService reads Test Scheduler pools and aggregate metrics.
// Requests require read_test_pool, build and suite access, and Test Scheduler
// API enablement. They do not additionally require read_suites.
type TestPoolsService struct {
	client *Client
}

// TestPool is a scheduler pool in planning, populating, consuming, consumed,
// or errored state. MutedTests is a planning snapshot, not a pending-work list.
type TestPool struct {
	ID             string                  `json:"id"`
	OrganizationID string                  `json:"organization_id"`
	SuiteID        string                  `json:"suite_id"`
	PipelineID     string                  `json:"pipeline_id"`
	BuildID        string                  `json:"build_id"`
	Key            string                  `json:"key"`
	State          string                  `json:"state"`
	AttemptPolicy  TestPoolAttemptPolicies `json:"attempt_policy"`
	Lease          TestPoolLeasePolicy     `json:"lease"`
	ExpiresAt      *Timestamp              `json:"expires_at"`
	CreatedAt      *Timestamp              `json:"created_at"`
	Error          *TestPoolError          `json:"error,omitempty"`
	MutedTests     []TestPoolMutedTest     `json:"muted_tests,omitempty"`
	// Nil means sizing was not requested or is not yet available; zero means
	// planning found no runnable tests.
	Parallelism *int `json:"parallelism,omitempty"`
}

type TestPoolError struct {
	Message string `json:"message"`
}

type TestPoolMutedTest struct {
	Scope *string `json:"scope"`
	Name  *string `json:"name"`
	Path  *string `json:"path"`
}

type TestPoolAttemptPolicy struct {
	MinAttempts      int  `json:"min_attempts"`
	MaxAttempts      int  `json:"max_attempts"`
	MinPassed        *int `json:"min_passed"`
	MaxPassed        int  `json:"max_passed"`
	MinFailed        *int `json:"min_failed"`
	MaxFailed        int  `json:"max_failed"`
	ParallelAttempts int  `json:"parallel_attempts"`
	InitialAttempts  *int `json:"initial_attempts"`
}

// TestPoolAttemptPolicies preserves the API's two representations. A single
// policy is in Single; multiple policies are in Named, including "default".
type TestPoolAttemptPolicies struct {
	Single *TestPoolAttemptPolicy
	Named  map[string]TestPoolAttemptPolicy
}

func (p *TestPoolAttemptPolicies) UnmarshalJSON(data []byte) error {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}
	var policies TestPoolAttemptPolicies
	if _, single := fields["max_attempts"]; single {
		if err := json.Unmarshal(data, &policies.Single); err != nil {
			return err
		}
	} else if err := json.Unmarshal(data, &policies.Named); err != nil {
		return err
	}
	*p = policies
	return nil
}

func (p TestPoolAttemptPolicies) MarshalJSON() ([]byte, error) {
	if p.Single != nil {
		return json.Marshal(p.Single)
	}
	return json.Marshal(p.Named)
}

// TestPoolCosts contains the selected lease cost dimension: custom,
// duration_p50_ms, or duration_p90_ms. Duration costs are in milliseconds.
type TestPoolCosts map[string]float64

type TestPoolLeasePolicy struct {
	Costs       TestPoolCosts `json:"costs"`
	MaxAttempts int           `json:"max_attempts"`
}

// TestPoolMetrics is a replica-backed, point-in-time snapshot, not a success
// verdict. Neither Drained nor zero active leases alone establishes that all
// tests completed or passed. Metrics are unavailable for expired pools.
type TestPoolMetrics struct {
	ObservedAt *Timestamp              `json:"observed_at"`
	Pool       TestPoolMetricsPool     `json:"pool"`
	Entries    TestPoolMetricsEntries  `json:"entries"`
	Attempts   TestPoolMetricsAttempts `json:"attempts"`
	Leases     TestPoolMetricsLeases   `json:"leases"`
	Timing     TestPoolMetricsTiming   `json:"timing"`
}

type TestPoolMetricsPool struct {
	ID        string     `json:"id"`
	State     string     `json:"state"`
	Drained   bool       `json:"drained"`
	ExpiresAt *Timestamp `json:"expires_at"`
}

type TestPoolMetricsEntries struct {
	Total int64 `json:"total"`
}

type TestPoolMetricsAttempts struct {
	Total int64 `json:"total"`
	// States uses waiting, leased, completed, and canceled keys.
	States map[string]TestPoolMetricsAttemptState `json:"states"`
	// Results uses passed, failed, and errored keys.
	Results map[string]int64 `json:"results"`
}

type TestPoolMetricsAttemptState struct {
	Count int64         `json:"count"`
	Costs TestPoolCosts `json:"costs"`
}

type TestPoolMetricsLeases struct {
	Active  TestPoolMetricsLease `json:"active"`
	Expired TestPoolMetricsLease `json:"expired"`
}

type TestPoolMetricsLease struct {
	Count        int64 `json:"count"`
	AttemptCount int64 `json:"attempt_count"`
}

type TestPoolMetricsTiming struct {
	OldestWaitingAttemptCreatedAt *Timestamp `json:"oldest_waiting_attempt_created_at"`
	OldestActiveLeaseStartedAt    *Timestamp `json:"oldest_active_lease_started_at"`
	OldestExpiredLeaseStartedAt   *Timestamp `json:"oldest_expired_lease_started_at"`
	NextActiveLeaseExpiresAt      *Timestamp `json:"next_active_lease_expires_at"`
}

// TestPoolsListOptions selects one page of pools. PipelineID and BuildID must
// be UUIDs. Key is an optional case-insensitive exact match. Supply After OR
// Before, never both. PerPage defaults to 30; the API accepts 1–100.
type TestPoolsListOptions struct {
	PipelineID string `url:"pipeline_id"`
	BuildID    string `url:"build_id"`
	Key        string `url:"key,omitempty"`
	After      string `url:"after,omitempty"`
	Before     string `url:"before,omitempty"`
	PerPage    int    `url:"per_page,omitempty"`
}

type TestPoolsListLink string

// ToOptions preserves the build, pipeline and key filters as well as the
// opaque cursor when following a continuation link.
func (l TestPoolsListLink) ToOptions() (*TestPoolsListOptions, error) {
	u, err := url.Parse(string(l))
	if err != nil {
		return nil, fmt.Errorf("parsing link: %w", err)
	}
	q := u.Query()
	opt := &TestPoolsListOptions{
		PipelineID: q.Get("pipeline_id"),
		BuildID:    q.Get("build_id"),
		Key:        q.Get("key"),
		After:      q.Get("after"),
		Before:     q.Get("before"),
	}
	if perPage := q.Get("per_page"); perPage != "" {
		opt.PerPage, err = strconv.Atoi(perPage)
		if err != nil {
			return nil, fmt.Errorf("parsing per_page: %w", err)
		}
	}
	return opt, nil
}

type TestPoolsListLinks struct {
	Self     TestPoolsListLink `json:"self"`
	First    TestPoolsListLink `json:"first,omitempty"`
	Previous TestPoolsListLink `json:"prev,omitempty"`
	Next     TestPoolsListLink `json:"next,omitempty"`
}

type TestPoolsList struct {
	Items []TestPool         `json:"items"`
	Links TestPoolsListLinks `json:"links"`
}

// List returns one cursor-paginated page, newest UUID first. It does not fetch
// subsequent pages. A 404 can mean unavailable, inaccessible, expired or missing
// resources; errors are not converted to empty lists.
func (s *TestPoolsService) List(ctx context.Context, org string, opt *TestPoolsListOptions) (TestPoolsList, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/test-scheduler/pools", org)
	u, err := addOptions(u, opt)
	if err != nil {
		return TestPoolsList{}, nil, err
	}
	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return TestPoolsList{}, nil, err
	}
	var pools TestPoolsList
	resp, err := s.client.Do(req, &pools)
	if err != nil {
		return TestPoolsList{}, resp, err
	}
	return pools, resp, nil
}

// Get returns a pool by UUID. BuildID is the build UUID, not its build number;
// it is required for user API tokens (OIDC tokens can supply an empty string).
func (s *TestPoolsService) Get(ctx context.Context, org, poolID, buildID string) (TestPool, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/test-scheduler/pools/%s", org, poolID)
	u, err := addOptions(u, struct {
		BuildID string `url:"build_id,omitempty"`
	}{buildID})
	if err != nil {
		return TestPool{}, nil, err
	}
	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return TestPool{}, nil, err
	}
	var pool TestPool
	resp, err := s.client.Do(req, &pool)
	if err != nil {
		return TestPool{}, resp, err
	}
	return pool, resp, nil
}

// GetMetrics returns aggregate metrics for a pool UUID. BuildID is required
// for user API tokens; OIDC tokens can supply an empty string.
func (s *TestPoolsService) GetMetrics(ctx context.Context, org, poolID, buildID string) (TestPoolMetrics, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/test-scheduler/pools/%s/metrics", org, poolID)
	u, err := addOptions(u, struct {
		BuildID string `url:"build_id,omitempty"`
	}{buildID})
	if err != nil {
		return TestPoolMetrics{}, nil, err
	}
	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return TestPoolMetrics{}, nil, err
	}
	var metrics TestPoolMetrics
	resp, err := s.client.Do(req, &metrics)
	if err != nil {
		return TestPoolMetrics{}, resp, err
	}
	return metrics, resp, nil
}
