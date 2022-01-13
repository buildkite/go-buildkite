package buildkite

// JobScheduled is triggered when a command step job has been scheduled
// to run on an agent
//
// Buildkite API docs: https://buildkite.com/docs/apis/webhooks/job-events
type JobScheduledEvent struct {
	Event    *string   `json:"event"`
	Build    *Build    `json:"build"`
	Job      *Job      `json:"job"`
	Pipeline *Pipeline `json:"pipeline"`
	Sender   *User     `json:"sender"`
}

// PingEvent is triggered when a webhook notification setting is changed
//
// Buildkite API docs: https://buildkite.com/docs/apis/webhooks/ping-events
type PingEvent struct {
	Event        *string       `json:"event"`
	Organization *Organization `json:"organization"`
	Sender       *User         `json:"sender"`
}
