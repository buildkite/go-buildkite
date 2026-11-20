package buildkite

import (
	"context"
	"encoding/json"
	"fmt"
)

// JobDiagnosticsOptions specifies optional evidence to read for a job.
type JobDiagnosticsOptions struct {
	// Include is a comma-separated list of "dependencies", "events" and
	// "insights". Omit it to leave these optional reads off.
	Include string `url:"include,omitempty"`
}

// JobDiagnostics contains recorded job facts and bounded diagnostic evidence.
// Consider each section's limitations when interpreting its values. Nullable
// fields and collections distinguish unavailable evidence from zero or empty.
// Unknown fields are ignored, except within dynamic event details and payloads.
type JobDiagnostics struct {
	Resource     JobDiagnosticsResource     `json:"resource"`
	ObservedAt   *Timestamp                 `json:"observed_at"`
	Details      JobDiagnosticsDetails      `json:"details"`
	Events       JobDiagnosticsEvents       `json:"events"`
	Dependencies JobDiagnosticsDependencies `json:"dependencies"`
	Concurrency  struct {
		Items       []JobDiagnosticsConcurrencyItem `json:"items"`
		Limitations map[string]bool                 `json:"limitations"`
	} `json:"concurrency"`
	Retries  JobDiagnosticsRetries `json:"retries"`
	Capacity struct {
		Requirements *struct {
			AgentQueryRules []string `json:"agent_query_rules"`
		} `json:"requirements"`
		Evidence    json.RawMessage      `json:"evidence"`
		Limitations map[string]bool      `json:"limitations"`
		Links       []JobDiagnosticsLink `json:"links"`
	} `json:"capacity"`
	Errors  JobDiagnosticsUnsupported `json:"errors"`
	History struct {
		RecentRuns JobDiagnosticsUnsupported `json:"recent_runs"`
	} `json:"history"`
	Insights []JobDiagnosticsQueueInsights `json:"insights"`
	Context  JobDiagnosticsContext         `json:"context"`
}

type JobDiagnosticsReference struct {
	Type string  `json:"type"`
	UUID *string `json:"uuid"`
}

type JobDiagnosticsResource struct {
	JobDiagnosticsReference
	URL string `json:"url"`
}

type JobDiagnosticsStepReference struct {
	JobDiagnosticsReference
	Key *string `json:"key"`
}

type JobDiagnosticsLink struct {
	Rel string `json:"rel"`
	URL string `json:"url"`
}

// JobDiagnosticsJobDetails contains facts shared by the target and related jobs.
type JobDiagnosticsJobDetails struct {
	Name         *string                 `json:"name"`
	Type         string                  `json:"type"`
	Lifecycle    JobDiagnosticsLifecycle `json:"lifecycle"`
	Priority     *int                    `json:"priority"`
	ExitStatus   *int                    `json:"exit_status"`
	SoftFailed   bool                    `json:"soft_failed"`
	BrokenReason *string                 `json:"broken_reason"`
	Signal       *string                 `json:"signal"`
	SignalReason *string                 `json:"signal_reason"`
	RetriesCount *int                    `json:"retries_count"`
	Links        []JobDiagnosticsLink    `json:"links"`
}

type JobDiagnosticsDetails struct {
	JobDiagnosticsJobDetails
	Timings struct {
		RunnableToStartSeconds *float64 `json:"runnable_to_start_seconds"`
		ExecutionSeconds       *float64 `json:"execution_seconds"`
	} `json:"timings"`
	Hosted        *bool `json:"hosted"`
	InstanceShape *struct {
		Name         string  `json:"name"`
		MachineType  string  `json:"machine_type"`
		Architecture string  `json:"architecture"`
		CPU          int     `json:"cpu"`
		Memory       float64 `json:"memory"`
	} `json:"instance_shape"`
	Agent   *JobDiagnosticsReference `json:"agent"`
	Step    *JobDiagnosticsReference `json:"step"`
	Queue   *JobDiagnosticsReference `json:"queue"`
	Routing struct {
		Type string `json:"type"`
	} `json:"routing"`
	RetrySource *struct {
		Job       JobDiagnosticsReference `json:"job"`
		RetryType *string                 `json:"retry_type"`
	} `json:"retry_source"`
	RetriedInJob *JobDiagnosticsReference `json:"retried_in_job"`
	Unblockable  *bool                    `json:"unblockable"`
	UnblockedAt  *Timestamp               `json:"unblocked_at"`
}

type JobDiagnosticsLifecycle struct {
	State       string     `json:"state"`
	CreatedAt   *Timestamp `json:"created_at"`
	ScheduledAt *Timestamp `json:"scheduled_at"`
	RunnableAt  *Timestamp `json:"runnable_at"`
	StartedAt   *Timestamp `json:"started_at"`
	EndedAt     *Timestamp `json:"ended_at"`
}

type JobDiagnosticsEvents struct {
	Window    ClusterQueueInsightsWindow `json:"window"`
	Selection struct {
		Order    string `json:"order"`
		MaxItems int    `json:"max_items"`
	} `json:"selection"`
	Items       []JobDiagnosticsEvent `json:"items"`
	Limitations map[string]bool       `json:"limitations"`
}

// JobDiagnosticsEvent preserves type-specific and future event detail fields,
// including nested values and integers too large for float64.
type JobDiagnosticsEvent struct {
	UUID       string                     `json:"uuid"`
	Type       string                     `json:"type"`
	OccurredAt *Timestamp                 `json:"occurred_at"`
	Actor      JobDiagnosticsReference    `json:"actor"`
	Details    map[string]json.RawMessage `json:"details"`
}

type JobDiagnosticsDependencies struct {
	Selection struct {
		Order    string `json:"order"`
		MaxSteps int    `json:"max_steps"`
		MaxEdges int    `json:"max_edges"`
	} `json:"selection"`
	Items           []JobDiagnosticsDependency    `json:"items"`
	UnexpandedSteps []JobDiagnosticsStepReference `json:"unexpanded_steps,omitempty"`
	Limitations     map[string]bool               `json:"limitations"`
	Links           []JobDiagnosticsLink          `json:"links"`
}

type JobDiagnosticsDependency struct {
	RequiredBy    JobDiagnosticsStepReference  `json:"required_by"`
	Step          JobDiagnosticsStepReference  `json:"step"`
	Source        string                       `json:"source"`
	AllowFailure  bool                         `json:"allow_failure"`
	InheritedFrom *JobDiagnosticsStepReference `json:"inherited_from"`
}

type JobDiagnosticsConcurrencyItem struct {
	Requirements JobDiagnosticsConcurrencyRequirements `json:"requirements"`
	Evidence     *JobDiagnosticsConcurrencyEvidence    `json:"evidence"`
	Limitations  map[string]bool                       `json:"limitations"`
}

// JobDiagnosticsConcurrencyRequirements describes a concurrency_group or
// platform constraint. Method is present only for concurrency groups.
type JobDiagnosticsConcurrencyRequirements struct {
	Type     string                            `json:"type"`
	Resource JobDiagnosticsConcurrencyResource `json:"resource"`
	Limit    *int                              `json:"limit"`
	Method   *string                           `json:"method"`
}

func (r JobDiagnosticsConcurrencyRequirements) MarshalJSON() ([]byte, error) {
	if r.Type == "platform" {
		return json.Marshal(struct {
			Type     string                  `json:"type"`
			Resource JobDiagnosticsReference `json:"resource"`
			Limit    *int                    `json:"limit"`
		}{r.Type, r.Resource.JobDiagnosticsReference, r.Limit})
	}
	type requirements JobDiagnosticsConcurrencyRequirements
	return json.Marshal(requirements(r))
}

// JobDiagnosticsConcurrencyResource includes nullable group identity fields.
// Platform requirements use an organization reference instead.
type JobDiagnosticsConcurrencyResource struct {
	JobDiagnosticsReference
	Key   *string `json:"key"`
	Scope *string `json:"scope"`
}

// JobDiagnosticsConcurrencyEvidence contains either group selection and jobs
// ahead, or a nullable platform_limited_at timestamp for platform constraints.
type JobDiagnosticsConcurrencyEvidence struct {
	ObservedAt        *Timestamp                          `json:"observed_at"`
	Selection         *JobDiagnosticsConcurrencySelection `json:"selection"`
	JobsAhead         []JobDiagnosticsReference           `json:"jobs_ahead"`
	PlatformLimitedAt *Timestamp                          `json:"platform_limited_at"`
}

type JobDiagnosticsConcurrencySelection struct {
	Population string   `json:"population"`
	Order      []string `json:"order"`
}

func (e JobDiagnosticsConcurrencyEvidence) MarshalJSON() ([]byte, error) {
	if e.Selection != nil {
		return json.Marshal(struct {
			ObservedAt *Timestamp                          `json:"observed_at"`
			Selection  *JobDiagnosticsConcurrencySelection `json:"selection"`
			JobsAhead  []JobDiagnosticsReference           `json:"jobs_ahead"`
		}{e.ObservedAt, e.Selection, e.JobsAhead})
	}
	return json.Marshal(struct {
		ObservedAt        *Timestamp `json:"observed_at"`
		PlatformLimitedAt *Timestamp `json:"platform_limited_at"`
	}{e.ObservedAt, e.PlatformLimitedAt})
}

type JobDiagnosticsRetries struct {
	Requirements *struct {
		Automatic []JobDiagnosticsRetryRule `json:"automatic"`
		Manual    struct {
			Allowed        bool    `json:"allowed"`
			PermitOnPassed bool    `json:"permit_on_passed"`
			Reason         *string `json:"reason"`
		} `json:"manual"`
	} `json:"requirements"`
	Evidence *struct {
		ObservedAt                    *Timestamp `json:"observed_at"`
		AutomaticRetriesBeforeAttempt []struct {
			RuleIndex int `json:"rule_index"`
			Count     int `json:"count"`
		} `json:"automatic_retries_before_attempt"`
	} `json:"evidence"`
	Limitations map[string]bool      `json:"limitations"`
	Links       []JobDiagnosticsLink `json:"links"`
}

// JobDiagnosticsRetryRule retains matcher unions as JSON: exit_status is a
// number, "*", or number array; signal and signal_reason are strings or string
// arrays. Matchers may be null when no value was recorded.
type JobDiagnosticsRetryRule struct {
	ExitStatus   json.RawMessage `json:"exit_status"`
	Signal       json.RawMessage `json:"signal"`
	SignalReason json.RawMessage `json:"signal_reason"`
	Limit        *int            `json:"limit"`
}

// JobDiagnosticsUnsupported reserves payloads whose schemas are not defined yet.
type JobDiagnosticsUnsupported struct {
	Items       json.RawMessage `json:"items"`
	Limitations map[string]bool `json:"limitations"`
}

// JobDiagnosticsQueueInsights reuses queue metrics with a nullable current
// snapshot when insights were not requested. Limitations is absent when read.
type JobDiagnosticsQueueInsights struct {
	Resource JobDiagnosticsReference `json:"resource"`
	Current  *struct {
		Agents ClusterQueueInsightsSnapshot `json:"agents"`
		Jobs   ClusterQueueInsightsSnapshot `json:"jobs"`
	} `json:"current"`
	Recent      *ClusterQueueInsightsRecent `json:"recent"`
	Limitations map[string]bool             `json:"limitations,omitempty"`
	Links       []JobDiagnosticsLink        `json:"links"`
}

type JobDiagnosticsContext struct {
	Scope struct {
		Organization struct {
			Resource struct {
				JobDiagnosticsResource
				Slug string `json:"slug"`
			} `json:"resource"`
			ObservedAt *Timestamp `json:"observed_at"`
			Details    struct {
				DispatchPaused bool `json:"dispatch_paused"`
			} `json:"details"`
		} `json:"organization"`
		Pipeline struct {
			Resource struct {
				JobDiagnosticsResource
				Slug string `json:"slug"`
			} `json:"resource"`
		} `json:"pipeline"`
		Build   JobDiagnosticsBuildRecord `json:"build"`
		Step    *JobDiagnosticsStepRecord `json:"step"`
		Cluster *struct {
			Resource JobDiagnosticsResource `json:"resource"`
		} `json:"cluster"`
		Queue *JobDiagnosticsQueueRecord `json:"queue"`
	} `json:"scope"`
	Steps       []JobDiagnosticsStepRecord  `json:"steps"`
	Agents      []JobDiagnosticsAgentRecord `json:"agents"`
	Queues      []JobDiagnosticsQueueRecord `json:"queues"`
	Jobs        []JobDiagnosticsJobRecord   `json:"jobs"`
	Builds      []JobDiagnosticsBuildRecord `json:"builds"`
	Limitations map[string]bool             `json:"limitations"`
}

type JobDiagnosticsBuildRecord struct {
	Resource struct {
		JobDiagnosticsResource
		Number int `json:"number"`
	} `json:"resource"`
	ObservedAt *Timestamp `json:"observed_at"`
	Details    struct {
		State        string  `json:"state"`
		Blocked      bool    `json:"blocked"`
		CancelReason *string `json:"cancel_reason"`
		Branch       string  `json:"branch"`
		Commit       string  `json:"commit"`
	} `json:"details"`
}

type JobDiagnosticsStepRecord struct {
	Resource   JobDiagnosticsStepReference `json:"resource"`
	ObservedAt *Timestamp                  `json:"observed_at"`
	Details    struct {
		Type      string                   `json:"type"`
		Label     *string                  `json:"label"`
		Group     *JobDiagnosticsReference `json:"group"`
		Lifecycle struct {
			CreatedAt      *Timestamp `json:"created_at"`
			State          *string    `json:"state"`
			Outcome        *string    `json:"outcome"`
			StateChangedAt *Timestamp `json:"state_changed_at"`
		} `json:"lifecycle"`
	} `json:"details"`
}

type JobDiagnosticsQueueRecord struct {
	Resource struct {
		JobDiagnosticsResource
		Key string `json:"key"`
	} `json:"resource"`
	ObservedAt *Timestamp `json:"observed_at"`
	Details    struct {
		DispatchPaused     bool                     `json:"dispatch_paused"`
		DispatchPausedAt   *Timestamp               `json:"dispatch_paused_at"`
		DispatchPausedBy   *JobDiagnosticsReference `json:"dispatch_paused_by"`
		DispatchPausedNote *string                  `json:"dispatch_paused_note"`
		DispatchSuspended  bool                     `json:"dispatch_suspended"`
	} `json:"details"`
}

type JobDiagnosticsAgentRecord struct {
	Resource   JobDiagnosticsResource `json:"resource"`
	ObservedAt *Timestamp             `json:"observed_at"`
	Details    struct {
		ConnectionState string                   `json:"connection_state"`
		Metadata        []string                 `json:"metadata"`
		Queue           *JobDiagnosticsReference `json:"queue"`
		ConnectedAt     *Timestamp               `json:"connected_at"`
		DisconnectedAt  *Timestamp               `json:"disconnected_at"`
		LostAt          *Timestamp               `json:"lost_at"`
		StoppedAt       *Timestamp               `json:"stopped_at"`
		Paused          bool                     `json:"paused"`
		PausedAt        *Timestamp               `json:"paused_at"`
		PausedNote      *string                  `json:"paused_note"`
	} `json:"details"`
	Events JobDiagnosticsEvents `json:"events"`
}

type JobDiagnosticsJobRecord struct {
	Resource   JobDiagnosticsResource `json:"resource"`
	ObservedAt *Timestamp             `json:"observed_at"`
	Details    struct {
		JobDiagnosticsJobDetails
		Build JobDiagnosticsReference `json:"build"`
	} `json:"details"`
}

// GetJobDiagnosticsByOrg returns diagnostic evidence by organization and job
// UUID. Requires read_builds and the Job Diagnostics feature enabled.
func (js *JobsService) GetJobDiagnosticsByOrg(ctx context.Context, org, jobID string, opt *JobDiagnosticsOptions) (JobDiagnostics, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/jobs/%s/diagnostics", org, jobID)
	u, err := addOptions(u, opt)
	if err != nil {
		return JobDiagnostics{}, nil, err
	}
	req, err := js.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return JobDiagnostics{}, nil, err
	}

	var diagnostics JobDiagnostics
	resp, err := js.client.Do(req, &diagnostics)
	if err != nil {
		return JobDiagnostics{}, resp, err
	}
	return diagnostics, resp, nil
}
