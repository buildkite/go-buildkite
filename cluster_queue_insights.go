package buildkite

import (
	"context"
	"fmt"
)

// ClusterQueueInsightsOptions specifies optional parameters for
// [ClusterQueuesService.Insights]. Invalid values and combinations are
// validated by the API.
type ClusterQueueInsightsOptions struct {
	// EndAtLatestPublishedMinute aligns history comparisons to the latest common
	// publication cutoff. True requires "history_comparison" and a "1h" time
	// period (the default). Nil or false retains the requested windows.
	EndAtLatestPublishedMinute *bool `url:"end_at_latest_published_minute,omitempty"`

	// Include is a comma-separated list of "history", "history_comparison"
	// and "hourly_changes". Each can be requested independently.
	// "history_comparison" returns both history periods and their changes.
	// "hourly_changes" returns hourly averages for now, yesterday and last week.
	Include string `url:"include,omitempty"`

	// TimePeriod is "1h" (the default), "24h" or "7d". It requires "history"
	// or "history_comparison" and does not affect hourly changes.
	TimePeriod string `url:"time_period,omitempty"`
}

// ClusterQueueInsights contains current counts, recent events and optional
// historical comparisons. Nil metric values mean unavailable, not zero.
// Limitations describe missing or incomplete evidence and must be considered
// when interpreting values.
type ClusterQueueInsights struct {
	Resource struct {
		Type string `json:"type"`
		UUID string `json:"uuid"`
	} `json:"resource"`
	Current struct {
		Agents ClusterQueueInsightsSnapshot `json:"agents"`
		Jobs   ClusterQueueInsightsSnapshot `json:"jobs"`
	} `json:"current"`
	Recent        *ClusterQueueInsightsRecent        `json:"recent"`
	History       *ClusterQueueInsightsHistory       `json:"history,omitempty"`
	HourlyChanges *ClusterQueueInsightsHourlyChanges `json:"hourly_changes,omitempty"`
}

// ClusterQueueInsightsSnapshot contains current agent or job counts.
// Agent values are connected, busy, idle, paused and draining. Job values are
// counts by state, with waiting reported as blocked_on_dependency.
type ClusterQueueInsightsSnapshot struct {
	ObservedAt            *Timestamp        `json:"observed_at"`
	SourceHighWatermarkAt *Timestamp        `json:"source_high_watermark_at"`
	Values                map[string]*int64 `json:"values"`
	// Computed is present for jobs only, with awaiting_execution and in_progress.
	Computed    map[string]*int64 `json:"computed,omitempty"`
	Limitations map[string]bool   `json:"limitations"`
}

// ClusterQueueInsightsWindow identifies a UTC time window. IntervalSeconds is
// present for history periods only; buckets start at From and exclude To.
type ClusterQueueInsightsWindow struct {
	From            *Timestamp `json:"from"`
	To              *Timestamp `json:"to"`
	IntervalSeconds int        `json:"interval_seconds,omitempty"`
}

// ClusterQueueInsightsRecent contains the most recent published event metrics.
type ClusterQueueInsightsRecent struct {
	JobCompletions ClusterQueueInsightsCompletions `json:"job_completions"`
	WaitTimes      ClusterQueueInsightsWaitTimes   `json:"wait_times"`
}

// ClusterQueueInsightsCompletions contains total, passed and failed job counts.
type ClusterQueueInsightsCompletions struct {
	Window      ClusterQueueInsightsWindow `json:"window"`
	ObservedAt  *Timestamp                 `json:"observed_at"`
	Values      map[string]*int64          `json:"values"`
	Limitations map[string]bool            `json:"limitations"`
}

// ClusterQueueInsightsWaitTimes contains min, p50, p95 and max wait times in
// seconds. Samples contains valid_count and unusable_duration_count.
type ClusterQueueInsightsWaitTimes struct {
	Window      ClusterQueueInsightsWindow `json:"window"`
	ObservedAt  *Timestamp                 `json:"observed_at"`
	Unit        string                     `json:"unit"`
	Values      map[string]*float64        `json:"values"`
	Samples     map[string]*int64          `json:"samples"`
	Limitations map[string]bool            `json:"limitations"`
}

// ClusterQueueInsightsHistory contains the selected period. PreviousPeriod and
// Changes are present only when history_comparison is requested.
type ClusterQueueInsightsHistory struct {
	SelectedPeriod ClusterQueueInsightsPeriod          `json:"selected_period"`
	PreviousPeriod *ClusterQueueInsightsPeriod         `json:"previous_period,omitempty"`
	Changes        *ClusterQueueInsightsHistoryChanges `json:"changes,omitempty"`
}

// ClusterQueueInsightsPeriod contains aligned buckets in chronological order.
// Limitations in series map codes to zero-based bucket indices.
type ClusterQueueInsightsPeriod struct {
	Window ClusterQueueInsightsWindow `json:"window"`
	Counts struct {
		Aggregation string                          `json:"aggregation"`
		Agents      ClusterQueueInsightsCountSeries `json:"agents"`
		Jobs        ClusterQueueInsightsCountSeries `json:"jobs"`
	} `json:"counts"`
	JobCompletions ClusterQueueInsightsCountSeries    `json:"job_completions"`
	WaitTimes      ClusterQueueInsightsWaitTimeSeries `json:"wait_times"`
}

// ClusterQueueInsightsCountSeries contains connected agents, awaiting_execution
// and in_progress jobs, or total, passed and failed completions.
// Counts use the last value in each bucket; completions count events per bucket.
type ClusterQueueInsightsCountSeries struct {
	// PublicationCoveredSeconds is present for job completions only. It records
	// publication coverage per bucket, not whether metric values are available.
	PublicationCoveredSeconds []float64           `json:"publication_covered_seconds,omitempty"`
	Values                    map[string][]*int64 `json:"values"`
	Limitations               map[string][]int    `json:"limitations"`
}

// ClusterQueueInsightsWaitTimeSeries contains min, mean, p50, p75, p90, p95 and max wait times.
// Samples contains valid_count and unusable_duration_count per bucket.
type ClusterQueueInsightsWaitTimeSeries struct {
	Unit string `json:"unit"`
	// PublicationCoveredSeconds records publication coverage per bucket,
	// independently of whether wait-time values are available.
	PublicationCoveredSeconds []float64             `json:"publication_covered_seconds"`
	Values                    map[string][]*float64 `json:"values"`
	Samples                   map[string][]*int64   `json:"samples"`
	Limitations               map[string][]int      `json:"limitations"`
}

// ClusterQueueInsightsHistoryChanges compares the selected and previous periods.
type ClusterQueueInsightsHistoryChanges struct {
	Counts         ClusterQueueInsightsCountChanges `json:"counts"`
	JobCompletions struct {
		Basis    string                     `json:"basis"`
		Total    ClusterQueueInsightsChange `json:"total"`
		PassRate struct {
			Unit                  string          `json:"unit"`
			Selected              *float64        `json:"selected"`
			Previous              *float64        `json:"previous"`
			PercentagePointChange *float64        `json:"percentage_point_change"`
			Limitations           map[string]bool `json:"limitations"`
		} `json:"pass_rate"`
	} `json:"job_completions"`
	WaitTimes struct {
		Basis string                     `json:"basis"`
		Unit  string                     `json:"unit"`
		Min   ClusterQueueInsightsChange `json:"min"`
		Mean  ClusterQueueInsightsChange `json:"mean"`
		P50   ClusterQueueInsightsChange `json:"p50"`
		P75   ClusterQueueInsightsChange `json:"p75"`
		P90   ClusterQueueInsightsChange `json:"p90"`
		P95   ClusterQueueInsightsChange `json:"p95"`
		Max   ClusterQueueInsightsChange `json:"max"`
	} `json:"wait_times"`
}

// ClusterQueueInsightsChange describes a numeric change. RelativeChangePercent
// is nil when the previous value is zero or either value is unavailable.
type ClusterQueueInsightsChange struct {
	Selected              *float64        `json:"selected"`
	Previous              *float64        `json:"previous"`
	Difference            *float64        `json:"difference"`
	RelativeChangePercent *float64        `json:"relative_change_percent"`
	Limitations           map[string]bool `json:"limitations"`
}

// ClusterQueueInsightsCountChanges compares connected agents and
// awaiting_execution and in_progress jobs. Basis is last_bucket for history;
// hourly comparisons specify completed_hour in their enclosing changes object.
type ClusterQueueInsightsCountChanges struct {
	Basis  string                                `json:"basis,omitempty"`
	Agents map[string]ClusterQueueInsightsChange `json:"agents"`
	Jobs   map[string]ClusterQueueInsightsChange `json:"jobs"`
}

// ClusterQueueInsightsHourlyChanges compares time-weighted average counts.
// Current contains the last completed hour and the partial current hour.
// Yesterday and LastWeek contain three completed hours around ReferenceAt.
type ClusterQueueInsightsHourlyChanges struct {
	Aggregation     string                           `json:"aggregation"`
	IntervalSeconds int                              `json:"interval_seconds"`
	Current         ClusterQueueInsightsHourlyPeriod `json:"current"`
	Yesterday       ClusterQueueInsightsHourlyPeriod `json:"yesterday"`
	LastWeek        ClusterQueueInsightsHourlyPeriod `json:"last_week"`
	Changes         struct {
		Basis     string                           `json:"basis"`
		Yesterday ClusterQueueInsightsCountChanges `json:"yesterday"`
		LastWeek  ClusterQueueInsightsCountChanges `json:"last_week"`
	} `json:"changes"`
}

// ClusterQueueInsightsHourlyPeriod contains hourly averages in chronological
// order. ReferenceAt is the current hour's start or its historical equivalent.
type ClusterQueueInsightsHourlyPeriod struct {
	ReferenceAt *Timestamp                       `json:"reference_at"`
	Window      ClusterQueueInsightsWindow       `json:"window"`
	Agents      ClusterQueueInsightsHourlySeries `json:"agents"`
	Jobs        ClusterQueueInsightsHourlySeries `json:"jobs"`
}

// ClusterQueueInsightsHourlySeries contains connected agents or
// awaiting_execution and in_progress jobs. CoveredSeconds records actual
// publication coverage per bucket, which can be partial or unavailable.
type ClusterQueueInsightsHourlySeries struct {
	Values         map[string][]*float64 `json:"values"`
	CoveredSeconds []*float64            `json:"covered_seconds"`
	Limitations    map[string][]int      `json:"limitations"`
}

// Insights returns a queue's current counts and recent metrics, plus requested
// history or hourly changes. It requires read_clusters scope, Cluster Insights
// entitlement and the Queue Insights REST API feature to be enabled.
func (cqs *ClusterQueuesService) Insights(ctx context.Context, org, clusterID, queueID string, opt *ClusterQueueInsightsOptions) (ClusterQueueInsights, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/clusters/%s/queues/%s/insights", org, clusterID, queueID)
	u, err := addOptions(u, opt)
	if err != nil {
		return ClusterQueueInsights{}, nil, err
	}

	req, err := cqs.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return ClusterQueueInsights{}, nil, err
	}

	var insights ClusterQueueInsights
	resp, err := cqs.client.Do(req, &insights)
	if err != nil {
		return ClusterQueueInsights{}, resp, err
	}

	return insights, resp, nil
}
