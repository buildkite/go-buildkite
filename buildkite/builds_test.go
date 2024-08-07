package buildkite

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestBuildsService_Cancel(t *testing.T) {
	t.Parallel()

	mux, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	mux.HandleFunc("/v2/organizations/my-great-org/pipelines/sup-keith/builds/1/cancel", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		fmt.Fprint(w, `{
  "id": "1",
  "state": "cancelled"
}`)
	})

	build, err := client.Builds.Cancel("my-great-org", "sup-keith", "1")
	if err != nil {
		t.Errorf("Cancel returned error: %v", err)
	}

	want := &Build{ID: String("1"), State: String("cancelled")}
	if !reflect.DeepEqual(build, want) {
		t.Errorf("Cancel returned %+v, want %+v", build, want)
	}
}

func TestBuildsService_List(t *testing.T) {
	t.Parallel()

	mux, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	mux.HandleFunc("/v2/builds", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"id":"123"},{"id":"1234"}]`)
	})

	builds, _, err := client.Builds.List(nil)
	if err != nil {
		t.Errorf("Builds.List returned error: %v", err)
	}

	want := []Build{{ID: String("123")}, {ID: String("1234")}}
	if !reflect.DeepEqual(builds, want) {
		t.Errorf("Builds.List returned %+v, want %+v", builds, want)
	}
}

func TestBuildsService_Get(t *testing.T) {
	t.Parallel()

	buildNumber := "123"
	orgName := "my-great-org"
	pipelineName := "sup-keith"
	requestSlug := fmt.Sprintf("/v2/organizations/%s/pipelines/%s/builds/%s",
		orgName, pipelineName, buildNumber)
	t.Run("returns a build struct with expected id", func(t *testing.T) {
		t.Parallel()

		mux, client, teardown := newMockServerAndClient(t)
		t.Cleanup(teardown)

		mux.HandleFunc(requestSlug,
			func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "GET")
				_, _ = fmt.Fprintf(w, `{"id":"%s"}`, buildNumber)
			})

		build, _, err := client.Builds.Get(orgName, pipelineName, buildNumber, nil)
		if err != nil {
			t.Errorf("Builds.Get (expected id) returned error: %v", err)
		}

		want := &Build{ID: String(buildNumber)}
		if !reflect.DeepEqual(build, want) {
			t.Errorf("Builds.Get (expected id) returned %+v, want %+v", build, want)
		}
	})

	t.Run("returns a build struct with expected job containing a group key", func(t *testing.T) {
		t.Parallel()

		mux, client, teardown := newMockServerAndClient(t)
		t.Cleanup(teardown)

		expectedGroup := "job_group"
		mux.HandleFunc(requestSlug,
			func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "GET")
				_, _ = fmt.Fprintf(w, `{"id":"%s", "jobs": [ {"group_key": "%s" }]}`,
					buildNumber,
					expectedGroup,
				)
			})

		build, _, err := client.Builds.Get(orgName, pipelineName, buildNumber, nil)
		if err != nil {
			t.Errorf("Builds.Get (group key) returned error: %v", err)
		}

		want := &Build{ID: String(buildNumber), Jobs: []*Job{{GroupKey: &expectedGroup}}}
		if !reflect.DeepEqual(build, want) {
			t.Errorf("Builds.Get (group key) returned %+v, want %+v", build, want)
		}
	})

	t.Run("returns a build struct with expected manual job values", func(t *testing.T) {
		t.Parallel()

		mux, client, teardown := newMockServerAndClient(t)
		t.Cleanup(teardown)

		jobType := "manual"
		unblockedAt := "2023-01-01T15:00:00.00Z"
		parsedTime := must(time.Parse(BuildKiteDateFormat, unblockedAt))

		mux.HandleFunc(requestSlug,
			func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "GET")
				_, _ = fmt.Fprintf(w, `{"id":"%s", "jobs": [ {"type": "%s", "unblocked_at": "%s" }]}`,
					buildNumber,
					jobType,
					unblockedAt,
				)
			})

		build, _, err := client.Builds.Get(orgName, pipelineName, buildNumber, nil)
		if err != nil {
			t.Errorf("Builds.Get (manual job) returned error: %v", err)
		}

		want := &Build{ID: String(buildNumber), Jobs: []*Job{{Type: &jobType, UnblockedAt: NewTimestamp(parsedTime)}}}
		if !reflect.DeepEqual(build, want) {
			t.Errorf("Builds.Get (manual job) returned %+v, want %+v", build, want)
		}
	})
}

func TestBuildsService_List_by_status(t *testing.T) {
	t.Parallel()

	mux, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	mux.HandleFunc("/v2/builds", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"state[]": "running",
			"page":    "2",
		})
		fmt.Fprint(w, `[{"id":"123"},{"id":"1234"}]`)
	})

	opt := &BuildsListOptions{
		State:       []string{"running"},
		ListOptions: ListOptions{Page: 2},
	}
	builds, _, err := client.Builds.List(opt)
	if err != nil {
		t.Errorf("Builds.List returned error: %v", err)
	}

	want := []Build{{ID: String("123")}, {ID: String("1234")}}
	if !reflect.DeepEqual(builds, want) {
		t.Errorf("Builds.List returned %+v, want %+v", builds, want)
	}
}

func TestBuildsService_List_by_multiple_status(t *testing.T) {
	t.Parallel()

	mux, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	mux.HandleFunc("/v2/builds", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValuesList(t, r, valuesList{
			{"state[]", "running"},
			{"state[]", "scheduled"},
			{"page", "2"},
		})
		fmt.Fprint(w, `[{"id":"123"},{"id":"1234"}]`)
	})

	opt := &BuildsListOptions{
		State:       []string{"running", "scheduled"},
		ListOptions: ListOptions{Page: 2},
	}
	builds, _, err := client.Builds.List(opt)
	if err != nil {
		t.Errorf("Builds.List returned error: %v", err)
	}

	want := []Build{{ID: String("123")}, {ID: String("1234")}}
	if !reflect.DeepEqual(builds, want) {
		t.Errorf("Builds.List returned %+v, want %+v", builds, want)
	}
}

func TestBuildsService_List_by_created_date(t *testing.T) {
	t.Parallel()

	mux, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	ts, err := time.Parse(BuildKiteDateFormat, "2016-03-24T01:00:00Z")
	if err != nil {
		t.Fatal(err)
	}

	mux.HandleFunc("/v2/builds", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"created_from": "2016-03-24T01:00:00Z",
			"created_to":   "2016-03-24T02:00:00Z",
		})
		fmt.Fprint(w, `[{"id":"123"}]`)
	})

	opt := &BuildsListOptions{
		CreatedFrom: ts,
		CreatedTo:   ts.Add(time.Hour),
	}
	builds, _, err := client.Builds.List(opt)
	if err != nil {
		t.Errorf("Builds.List returned error: %v", err)
	}

	want := []Build{{ID: String("123")}}
	if !reflect.DeepEqual(builds, want) {
		t.Errorf("Builds.List returned %+v, want %+v", builds, want)
	}
}

func TestBuildsService_ListByOrg(t *testing.T) {
	t.Parallel()

	mux, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	mux.HandleFunc("/v2/organizations/my-great-org/builds", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"id":"123"},{"id":"1234"}]`)
	})

	builds, _, err := client.Builds.ListByOrg("my-great-org", nil)
	if err != nil {
		t.Errorf("Builds.List returned error: %v", err)
	}

	want := []Build{{ID: String("123")}, {ID: String("1234")}}
	if !reflect.DeepEqual(builds, want) {
		t.Errorf("Builds.List returned %+v, want %+v", builds, want)
	}
}

func TestBuildsService_ListByOrg_branch_commit(t *testing.T) {
	t.Parallel()

	mux, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	mux.HandleFunc("/v2/organizations/my-great-org/builds", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"branch[]": "my-great-branch",
			"commit":   "my-commit-sha1",
		})
		fmt.Fprint(w, `[{"id":"123"},{"id":"1234"}]`)
	})

	opt := &BuildsListOptions{
		Branch: []string{"my-great-branch"},
		Commit: "my-commit-sha1",
	}

	builds, _, err := client.Builds.ListByOrg("my-great-org", opt)
	if err != nil {
		t.Errorf("Builds.List returned error: %v", err)
	}

	want := []Build{{ID: String("123")}, {ID: String("1234")}}
	if !reflect.DeepEqual(builds, want) {
		t.Errorf("Builds.List returned %+v, want %+v", builds, want)
	}
}

func TestBuildsService_List_by_multiple_branches(t *testing.T) {
	t.Parallel()

	mux, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	mux.HandleFunc("/v2/builds", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValuesList(t, r, valuesList{
			{"branch[]", "my-great-branch"},
			{"branch[]", "my-other-great-branch"},
		})
		fmt.Fprint(w, `[{"id":"123"},{"id":"1234"}]`)
	})

	opt := &BuildsListOptions{
		Branch: []string{"my-great-branch", "my-other-great-branch"},
	}
	builds, _, err := client.Builds.List(opt)
	if err != nil {
		t.Errorf("Builds.List returned error: %v", err)
	}

	want := []Build{{ID: String("123")}, {ID: String("1234")}}
	if !reflect.DeepEqual(builds, want) {
		t.Errorf("Builds.List returned %+v, want %+v", builds, want)
	}
}

func TestBuildsService_ListByPipeline(t *testing.T) {
	t.Parallel()

	mux, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	mux.HandleFunc("/v2/organizations/my-great-org/pipelines/sup-keith/builds", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"id":"123"},{"id":"1234"}]`)
	})

	builds, _, err := client.Builds.ListByPipeline("my-great-org", "sup-keith", nil)
	if err != nil {
		t.Errorf("Builds.List returned error: %v", err)
	}

	want := []Build{{ID: String("123")}, {ID: String("1234")}}
	if !reflect.DeepEqual(builds, want) {
		t.Errorf("Builds.List returned %+v, want %+v", builds, want)
	}
}

func TestBuildsUnmarshalWebhook(t *testing.T) {
	// payload taken from buildkite services console
	sampleData := `{
  "event": "build.scheduled",
  "build": {
    "id": "foo",
    "url": "https://api.buildkite.com/v2/organizations/org/pipelines/greenpipe/builds/1",
    "web_url": "https://buildkite.com/org/greenpipe/builds/1",
    "number": 1,
    "state": "scheduled",
    "blocked": false,
    "message": "doot",
    "commit": "HEAD",
    "branch": "master",
    "tag": null,
    "source": "ui",
    "author": {
		"username": "foojim",
		"name": "Uhh, Jim",
		"email": "slam@space.jam"
	  },
    "creator": {
      "id": "foo",
      "name": "Uhh, Jim",
      "email": "slam@space.jam",
      "created_at": "2018-03-22 23:13:16 UTC"
    },
    "created_at": "2018-03-25 03:58:14 UTC",
    "scheduled_at": "2018-03-25 03:58:14 UTC"
  }
}`

	type webhookPayload struct {
		Event string
		Build Build
	}

	var payload webhookPayload

	if err := json.Unmarshal([]byte(sampleData), &payload); err != nil {
		t.Fatalf("could not unmarshal: %v", err)
	}
}
