package buildkite

// buildEvent is a wrapper for a build event notification
//
// Buildkite API docs: https://buildkite.com/docs/apis/webhooks/build-events
type buildEvent struct {
	Event    *string   `json:"event"`
	Build    *Build    `json:"build"`
	Pipeline *Pipeline `json:"pipeline"`
	Sender   *User     `json:"sender"`
}

// BuildFinishedEvent is triggered when a build finishes
//
// Buildkite API docs: https://buildkite.com/docs/apis/webhooks/build-events
type BuildFinishedEvent struct {
	buildEvent
}

// BuildRunningEvent is triggered when a build starts running
//
// Buildkite API docs: https://buildkite.com/docs/apis/webhooks/build-events
type BuildRunningEvent struct {
	buildEvent
}

// BuildScheduledEvent is triggered when a build is scheduled
//
// Buildkite API docs: https://buildkite.com/docs/apis/webhooks/build-events
type BuildScheduledEvent struct {
	buildEvent
}

// jobEvent is a wrapper for a job event notification
//
// Buildkite API docs: https://buildkite.com/docs/apis/webhooks/job-events
type jobEvent struct {
	Event    *string   `json:"event"`
	Build    *Build    `json:"build"`
	Job      *Job      `json:"job"`
	Pipeline *Pipeline `json:"pipeline"`
	Sender   *User     `json:"sender"`
}

// JobActivatedEvent is triggered when a job is activated
//
// Buildkite API docs: https://buildkite.com/docs/apis/webhooks/job-events
type JobActivatedEvent struct {
	jobEvent
}

// JobFinishedEvent is triggered when a job is finished
//
// Buildkite API docs: https://buildkite.com/docs/apis/webhooks/job-events
type JobFinishedEvent struct {
	jobEvent
}

// JobScheduledEvent is triggered when a job is scheduled
//
// Buildkite API docs: https://buildkite.com/docs/apis/webhooks/job-events
type JobScheduledEvent struct {
	jobEvent
}

// JobStartedEvent is triggered when a job is started
//
// Buildkite API docs: https://buildkite.com/docs/apis/webhooks/job-events
type JobStartedEvent struct {
	jobEvent
}

// PingEvent is triggered when a webhook notification setting is changed
//
// Buildkite API docs: https://buildkite.com/docs/apis/webhooks/ping-events
type PingEvent struct {
	Event        *string       `json:"event"`
	Organization *Organization `json:"organization"`
	Sender       *User         `json:"sender"`
}
