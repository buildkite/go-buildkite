package buildkite

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestAgentsService_List(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-great-org/agents", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = fmt.Fprint(w, `[{"id":"123"},{"id":"1234"}]`)
	})

	agents, _, err := client.Agents.List(context.Background(), "my-great-org", nil)
	if err != nil {
		t.Errorf("Agents.List returned error: %v", err)
	}

	want := []Agent{{ID: "123"}, {ID: "1234"}}
	if diff := cmp.Diff(agents, want); diff != "" {
		t.Errorf("Agents.List: diff: (-got +want)\n%s", diff)
	}
}

func TestAgentsService_Get(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-great-org/agents/123", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = fmt.Fprint(w, `{"id":"123"}`)
	})

	agent, _, err := client.Agents.Get(context.Background(), "my-great-org", "123")
	if err != nil {
		t.Errorf("Agents.Get returned error: %v", err)
	}

	want := Agent{ID: "123"}
	if diff := cmp.Diff(agent, want); diff != "" {
		t.Errorf("Agents.List(): diff: (-got +want)\n%s", diff)
	}
}

func TestAgentsService_Create(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	input := Agent{Name: "new_agent_bob"}

	server.HandleFunc("/v2/organizations/my-great-org/agents", func(w http.ResponseWriter, r *http.Request) {
		var v Agent
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("Error parsing json body: %v", err)
		}

		testMethod(t, r, "POST")

		if diff := cmp.Diff(v, input); diff != "" {
			t.Errorf("Request body diff: (-got +want)\n%s", diff)
		}

		_, _ = fmt.Fprint(w, `{"id":"123"}`)
	})

	agent, _, err := client.Agents.Create(context.Background(), "my-great-org", input)
	if err != nil {
		t.Errorf("Agents.Create returned error: %v", err)
	}

	want := Agent{ID: "123"}
	if diff := cmp.Diff(agent, want); diff != "" {
		t.Errorf("Agents.Create() diff: (-got +want)\n%s", diff)
	}
}

func TestAgentsService_Delete(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-great-org/agents/123", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	_, err := client.Agents.Delete(context.Background(), "my-great-org", "123")
	if err != nil {
		t.Errorf("Agents.Delete returned error: %v", err)
	}
}

func TestAgentsService_Stop(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-great-org/agents/123/stop", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("could not read request body: %v", err)
		}
		s := strings.TrimSpace(string(body))
		expected := `{"force":true}`
		if s != expected {
			t.Errorf("Request body: %s, want %s", s, expected)
		}
	})

	_, err := client.Agents.Stop(context.Background(), "my-great-org", "123", true)
	if err != nil {
		t.Errorf("Agents.Stop returned error: %v", err)
	}
}

func TestAgentsService_Pause(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-great-org/agents/123/pause", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("could not read request body: %v", err)
		}
		expectedBody := "{}"
		actualBody := strings.TrimSpace(string(body))
		if actualBody != expectedBody {
			t.Errorf("Expected request body %q, got %q", expectedBody, actualBody)
		}
	})

	_, err := client.Agents.Pause(context.Background(), "my-great-org", "123")
	if err != nil {
		t.Errorf("Agents.Pause returned error: %v", err)
	}
}

func TestAgentsService_PauseWithOptions(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-great-org/agents/123/pause", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("could not read request body: %v", err)
		}

		var requestBody AgentPauseOptions
		err = json.Unmarshal(body, &requestBody)
		if err != nil {
			t.Errorf("could not unmarshal request body: %v", err)
		}

		if requestBody.Note != "Maintenance scheduled" {
			t.Errorf("Expected note %q, got %q", "Maintenance scheduled", requestBody.Note)
		}

		if requestBody.TimeoutInMinutes != 60 {
			t.Errorf("Expected timeout %d, got %d", 60, requestBody.TimeoutInMinutes)
		}
	})

	opts := &AgentPauseOptions{
		Note:             "Maintenance scheduled",
		TimeoutInMinutes: 60,
	}

	_, err := client.Agents.PauseWithOptions(context.Background(), "my-great-org", "123", opts)
	if err != nil {
		t.Errorf("Agents.PauseWithOptions returned error: %v", err)
	}
}

func TestAgentsService_Resume(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-great-org/agents/123/resume", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("could not read request body: %v", err)
		}
		if strings.TrimSpace(string(body)) != "{}" {
			t.Errorf("Expected request body %q, got %q", "{}", strings.TrimSpace(string(body)))
		}
	})

	_, err := client.Agents.Resume(context.Background(), "my-great-org", "123")
	if err != nil {
		t.Errorf("Agents.Resume returned error: %v", err)
	}
}
