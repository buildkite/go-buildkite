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
		fmt.Fprint(w, `[{"id":"123"},{"id":"1234"}]`)
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
		fmt.Fprint(w, `{"id":"123"}`)
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

		fmt.Fprint(w, `{"id":"123"}`)
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
