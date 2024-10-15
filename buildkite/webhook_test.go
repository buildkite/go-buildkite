package buildkite

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParseWebHook(t *testing.T) {
	tests := []struct {
		payload      interface{}
		messageType  string
		errorMessage string
	}{
		{
			payload:     &AgentConnectedEvent{},
			messageType: "agent.connected",
		},
		{
			payload:     &AgentDisconnectedEvent{},
			messageType: "agent.disconnected",
		},
		{
			payload:     &AgentLostEvent{},
			messageType: "agent.lost",
		},
		{
			payload:     &AgentStoppedEvent{},
			messageType: "agent.stopped",
		},
		{
			payload:     &AgentStoppingEvent{},
			messageType: "agent.stopping",
		},
		{
			payload:     &BuildFailingEvent{},
			messageType: "build.failing",
		},
		{
			payload:     &BuildFinishedEvent{},
			messageType: "build.finished",
		},
		{
			payload:     &BuildRunningEvent{},
			messageType: "build.running",
		},
		{
			payload:     &BuildScheduledEvent{},
			messageType: "build.scheduled",
		},
		{
			payload:     &JobActivatedEvent{},
			messageType: "job.activated",
		},
		{
			payload:     &JobFinishedEvent{},
			messageType: "job.finished",
		},
		{
			payload:     &JobScheduledEvent{},
			messageType: "job.scheduled",
		},
		{
			payload:     &JobStartedEvent{},
			messageType: "job.started",
		},
		{
			payload:     &PingEvent{},
			messageType: "ping",
		},
		{
			payload:      &PingEvent{},
			messageType:  "invalid",
			errorMessage: "unknown X-Buildkite-Event in message: invalid",
		},
	}

	for _, test := range tests {
		p, err := json.Marshal(test.payload)
		if err != nil {
			t.Fatalf("Marshal(%#v): %v", test.payload, err)
		}

		got, err := ParseWebHook(test.messageType, p)
		if err != nil {
			if test.errorMessage != "" {
				if err.Error() != test.errorMessage {
					t.Errorf("ParseWebHook(%#v, %#v) expected error, got %#v", test.messageType, test.payload, err.Error())
				}
				continue
			}
			t.Fatalf("ParseWebHook: %v", err)
		}

		if diff := cmp.Diff(got, test.payload); diff != "" {
			t.Errorf("ParseWebHook(%q, []byte(%q)) returned unexpected output. diff: (-got +want)\n%s", test.messageType, string(p), diff)
		}
	}
}

func TestValidatePayload(t *testing.T) {
	const defaultBody = `{"event":"ping","service":{"id":"c9f8372d-c0cd-43dc-9274-768a875cf6ca","provider":"webhook","settings":{"url":"https://server.com/webhooks"}},"organization":{"id":"49801950-1df0-474f-bb56-ad6a930c5cb9","graphql_id":"T3JnYW5pemF0aW9uLS0tZTBmMzk3MgsTksGkxOWYtZTZjNzczZTJiYjEy","url":"https://api.buildkite.com/v2/organizations/acme-inc","web_url":"https://buildkite.com/acme-inc","name":"ACME Inc","slug":"acme-inc","agents_url":"https://api.buildkite.com/v2/organizations/acme-inc/agents","emojis_url":"https://api.buildkite.com/v2/organizations/acme-inc/emojis","created_at":"2021-02-03T20:34:10.486Z","pipelines_url":"https://api.buildkite.com/v2/organizations/acme-inc/pipelines"},"sender":{"id":"c9f8372d-c0cd-43dc-9269-bcbb7f308e3f","name":"ACME Man"}}`
	const defaultSignature = "timestamp=1642080837,signature=582d496ac2d869dd97a3101c4cda346288c49a742592daf582ec64c86449f79c"
	const errorDecodingSignature = "error decoding signature"
	const invalidSignatureHeader = "X-Buildkite-Signature format is incorrect."
	const missingSignatureHeader = "No X-Buildkite-Signature header present on request"
	const payloadSignatureError = "payload signature check failed"
	secretKey := []byte("29b1ff5779c76bd48ba6705eb99ff970")

	tests := []struct {
		signature   string
		event       string
		wantError   string
		wantEvent   string
		wantPayload string
	}{
		// The following tests generate expected errors:
		// Missing signature
		{
			signature: "",
			wantError: missingSignatureHeader,
		},
		// Invalid signature format
		{
			signature: "invalid",
			wantError: invalidSignatureHeader,
		},
		// Signature not hex string
		{
			signature: "timestamp=1642080837,signature=yo",
			wantError: errorDecodingSignature,
		},
		// Signature not valid
		{
			signature: strings.Replace(defaultSignature, "f", "a", 1),
			wantError: payloadSignatureError,
		},
		// The following tests expect a valid result
		{
			signature:   defaultSignature,
			event:       "ping",
			wantEvent:   "ping",
			wantPayload: defaultBody,
		},
	}

	for _, test := range tests {
		buf := bytes.NewBufferString(defaultBody)
		req, err := http.NewRequest("POST", "http://localhost/webhook", buf)
		if err != nil {
			t.Fatalf("NewRequest: %v", err)
		}

		if test.signature != "" {
			req.Header.Set(SignatureHeader, test.signature)
		}

		req.Header.Set("Content-Type", "application/json")

		got, err := ValidatePayload(req, secretKey)
		if err != nil {
			if test.wantPayload != "" {
				t.Errorf("ValidatePayload(%#v): err = %v, want nil", test, err)
			}

			if !strings.Contains(err.Error(), test.wantError) {
				t.Errorf("ValidatePayload(%#v): err = %s, want err = %s", test, err.Error(), test.wantError)
			}

			continue
		}

		if string(got) != test.wantPayload {
			t.Errorf("ValidatePayload = %q, want %q", got, test.wantPayload)
		}
	}
}

func TestWebHookType(t *testing.T) {
	eventType := "ping"

	req, err := http.NewRequest("POST", "http://localhost", nil)
	if err != nil {
		t.Fatalf("Error building requet: %v", err)
	}

	req.Header.Set(EventTypeHeader, eventType)

	got := WebHookType(req)
	if got != eventType {
		t.Errorf("WebHookType(%#v) = %q, want %q", req, got, eventType)
	}
}
