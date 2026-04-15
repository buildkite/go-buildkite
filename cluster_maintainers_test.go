package buildkite

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestClusterMaintainersService_List(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57/maintainers", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = fmt.Fprint(w, `[
			{"id": "aaa", "actor": {"id": "u1", "name": "Alice", "type": "user"}},
			{"id": "bbb", "actor": {"id": "t1", "name": "Ops Team", "type": "team"}}
		]`)
	})

	got, _, err := client.ClusterMaintainers.List(context.Background(), "my-org", "b7c9bc4f-526f-4c18-a3be-dc854ab75d57", nil)
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}

	want := []ClusterMaintainerEntry{
		{ID: "aaa", Actor: ClusterMaintainerActor{ID: "u1", Name: "Alice", Type: "user"}},
		{ID: "bbb", Actor: ClusterMaintainerActor{ID: "t1", Name: "Ops Team", Type: "team"}},
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("List mismatch (-want +got):\n%s", diff)
	}
}

func TestClusterMaintainersService_Get(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57/maintainers/aaa", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = fmt.Fprint(w, `{"id": "aaa", "actor": {"id": "u1", "name": "Alice", "type": "user"}}`)
	})

	got, _, err := client.ClusterMaintainers.Get(context.Background(), "my-org", "b7c9bc4f-526f-4c18-a3be-dc854ab75d57", "aaa")
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}

	want := ClusterMaintainerEntry{ID: "aaa", Actor: ClusterMaintainerActor{ID: "u1", Name: "Alice", Type: "user"}}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Get mismatch (-want +got):\n%s", diff)
	}
}

func TestClusterMaintainersService_Create(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57/maintainers", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		w.WriteHeader(http.StatusCreated)
		_, _ = fmt.Fprint(w, `{"id": "aaa", "actor": {"id": "u1", "name": "Alice", "type": "user"}}`)
	})

	got, _, err := client.ClusterMaintainers.Create(context.Background(), "my-org", "b7c9bc4f-526f-4c18-a3be-dc854ab75d57", ClusterMaintainer{UserID: "u1"})
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	want := ClusterMaintainerEntry{ID: "aaa", Actor: ClusterMaintainerActor{ID: "u1", Name: "Alice", Type: "user"}}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Create mismatch (-want +got):\n%s", diff)
	}
}

func TestClusterMaintainersService_Delete(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-org/clusters/b7c9bc4f-526f-4c18-a3be-dc854ab75d57/maintainers/aaa", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	_, err := client.ClusterMaintainers.Delete(context.Background(), "my-org", "b7c9bc4f-526f-4c18-a3be-dc854ab75d57", "aaa")
	if err != nil {
		t.Fatalf("Delete returned error: %v", err)
	}
}
