package buildkite

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestBuildsService_Cancel(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-great-org/pipelines/sup-keith/builds/1/cancel", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		_, _ = fmt.Fprint(w, `{
  "id": "1",
  "state": "cancelled"
}`)
	})

	build, err := client.Builds.Cancel(context.Background(), "my-great-org", "sup-keith", "1")
	if err != nil {
		t.Errorf("Cancel returned error: %v", err)
	}

	want := Build{ID: "1", State: "cancelled"}
	if diff := cmp.Diff(build, want); diff != "" {
		t.Errorf("Cancel diff: (-got +want)\n%s", diff)
	}
}

func TestBuildsService_List(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/builds", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = fmt.Fprint(w, `[{"id":"123"},{"id":"1234"}]`)
	})

	builds, _, err := client.Builds.List(context.Background(), nil)
	if err != nil {
		t.Errorf("Builds.List returned error: %v", err)
	}

	want := []Build{{ID: "123"}, {ID: "1234"}}
	if diff := cmp.Diff(builds, want); diff != "" {
		t.Errorf("Builds.List diff: (-got +want)\n%s", diff)
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

		server, client, teardown := newMockServerAndClient(t)
		t.Cleanup(teardown)

		server.HandleFunc(requestSlug,
			func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "GET")
				_, _ = fmt.Fprintf(w, `{"id":"%s"}`, buildNumber)
			})

		build, _, err := client.Builds.Get(context.Background(), orgName, pipelineName, buildNumber, nil)
		if err != nil {
			t.Errorf("Builds.Get (expected id) returned error: %v", err)
		}

		want := Build{ID: buildNumber}
		if diff := cmp.Diff(build, want); diff != "" {
			t.Errorf("Builds.Get (expected id) diff: (-got +want)\n%s", diff)
		}
	})

	t.Run("returns a build struct with expected job containing a group key", func(t *testing.T) {
		t.Parallel()

		server, client, teardown := newMockServerAndClient(t)
		t.Cleanup(teardown)

		expectedGroup := "job_group"
		server.HandleFunc(requestSlug,
			func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "GET")
				_, _ = fmt.Fprintf(w, `{"id":"%s", "jobs": [ {"group_key": "%s" }]}`,
					buildNumber,
					expectedGroup,
				)
			})

		build, _, err := client.Builds.Get(context.Background(), orgName, pipelineName, buildNumber, nil)
		if err != nil {
			t.Errorf("Builds.Get (group key) returned error: %v", err)
		}

		want := Build{ID: buildNumber, Jobs: []Job{{GroupKey: expectedGroup}}}
		if diff := cmp.Diff(build, want); diff != "" {
			t.Errorf("Builds.Get (group key) diff: (-got +want)\n%s", diff)
		}
	})

	t.Run("returns a build struct with expected manual job values", func(t *testing.T) {
		t.Parallel()

		server, client, teardown := newMockServerAndClient(t)
		t.Cleanup(teardown)

		jobType := "manual"
		unblockedAt := "2023-01-01T15:00:00.00Z"
		parsedTime := must(time.Parse(BuildKiteDateFormat, unblockedAt))

		server.HandleFunc(requestSlug,
			func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "GET")
				_, _ = fmt.Fprintf(w, `{"id":"%s", "jobs": [ {"type": "%s", "unblocked_at": "%s" }]}`,
					buildNumber,
					jobType,
					unblockedAt,
				)
			})

		build, _, err := client.Builds.Get(context.Background(), orgName, pipelineName, buildNumber, nil)
		if err != nil {
			t.Errorf("Builds.Get (manual job) returned error: %v", err)
		}

		want := Build{ID: buildNumber, Jobs: []Job{{Type: jobType, UnblockedAt: NewTimestamp(parsedTime)}}}
		if diff := cmp.Diff(build, want); diff != "" {
			t.Errorf("Builds.Get (manual job) diff: (-got +want)\n%s", diff)
		}
	})
}

func TestBuildsService_List_by_status(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/builds", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"state[]": "running",
			"page":    "2",
		})
		_, _ = fmt.Fprint(w, `[{"id":"123"},{"id":"1234"}]`)
	})

	opt := &BuildsListOptions{
		State:       []string{"running"},
		ListOptions: ListOptions{Page: 2},
	}
	builds, _, err := client.Builds.List(context.Background(), opt)
	if err != nil {
		t.Errorf("Builds.List returned error: %v", err)
	}

	want := []Build{{ID: "123"}, {ID: "1234"}}
	if diff := cmp.Diff(builds, want); diff != "" {
		t.Errorf("Builds.List diff: (-got +want)\n%s", diff)
	}
}

func TestBuildsService_List_by_multiple_status(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/builds", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValuesList(t, r, valuesList{
			{"state[]", "running"},
			{"state[]", "scheduled"},
			{"page", "2"},
		})
		_, _ = fmt.Fprint(w, `[{"id":"123"},{"id":"1234"}]`)
	})

	opt := &BuildsListOptions{
		State:       []string{"running", "scheduled"},
		ListOptions: ListOptions{Page: 2},
	}
	builds, _, err := client.Builds.List(context.Background(), opt)
	if err != nil {
		t.Errorf("Builds.List returned error: %v", err)
	}

	want := []Build{{ID: "123"}, {ID: "1234"}}
	if diff := cmp.Diff(builds, want); diff != "" {
		t.Errorf("Builds.List diff: (-got +want)\n%s", diff)
	}
}

func TestBuildsService_List_by_created_date(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	ts, err := time.Parse(BuildKiteDateFormat, "2016-03-24T01:00:00Z")
	if err != nil {
		t.Fatal(err)
	}

	server.HandleFunc("/v2/builds", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"created_from": "2016-03-24T01:00:00Z",
			"created_to":   "2016-03-24T02:00:00Z",
		})
		_, _ = fmt.Fprint(w, `[{"id":"123"}]`)
	})

	opt := &BuildsListOptions{
		CreatedFrom: ts,
		CreatedTo:   ts.Add(time.Hour),
	}
	builds, _, err := client.Builds.List(context.Background(), opt)
	if err != nil {
		t.Errorf("Builds.List returned error: %v", err)
	}

	want := []Build{{ID: "123"}}
	if diff := cmp.Diff(builds, want); diff != "" {
		t.Errorf("Builds.List diff: (-got +want)\n%s", diff)
	}
}

func TestBuildsService_ListByOrg(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-great-org/builds", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = fmt.Fprint(w, `[{"id":"123"},{"id":"1234"}]`)
	})

	builds, _, err := client.Builds.ListByOrg(context.Background(), "my-great-org", nil)
	if err != nil {
		t.Errorf("Builds.List returned error: %v", err)
	}

	want := []Build{{ID: "123"}, {ID: "1234"}}
	if diff := cmp.Diff(builds, want); diff != "" {
		t.Errorf("Builds.List diff: (-got +want)\n%s", diff)
	}
}

func TestBuildsService_ListByOrg_branch_commit(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-great-org/builds", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"branch[]": "my-great-branch",
			"commit":   "my-commit-sha1",
		})
		_, _ = fmt.Fprint(w, `[{"id":"123"},{"id":"1234"}]`)
	})

	opt := &BuildsListOptions{
		Branch: []string{"my-great-branch"},
		Commit: "my-commit-sha1",
	}

	builds, _, err := client.Builds.ListByOrg(context.Background(), "my-great-org", opt)
	if err != nil {
		t.Errorf("Builds.List returned error: %v", err)
	}

	want := []Build{{ID: "123"}, {ID: "1234"}}
	if diff := cmp.Diff(builds, want); diff != "" {
		t.Errorf("Builds.List diff: (-got +want)\n%s", diff)
	}
}

func TestBuildsService_List_by_multiple_branches(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/builds", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValuesList(t, r, valuesList{
			{"branch[]", "my-great-branch"},
			{"branch[]", "my-other-great-branch"},
		})
		_, _ = fmt.Fprint(w, `[{"id":"123"},{"id":"1234"}]`)
	})

	opt := &BuildsListOptions{
		Branch: []string{"my-great-branch", "my-other-great-branch"},
	}
	builds, _, err := client.Builds.List(context.Background(), opt)
	if err != nil {
		t.Errorf("Builds.List returned error: %v", err)
	}

	want := []Build{{ID: "123"}, {ID: "1234"}}
	if diff := cmp.Diff(builds, want); diff != "" {
		t.Errorf("Builds.List diff: (-got +want)\n%s", diff)
	}
}

func TestBuildsService_List_with_exclude_jobs(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/builds", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"exclude_jobs": "true",
		})
		_, _ = fmt.Fprint(w, `[{"id":"123"},{"id":"1234"}]`)
	})

	opt := &BuildsListOptions{
		ExcludeJobs: true,
	}
	builds, _, err := client.Builds.List(context.Background(), opt)
	if err != nil {
		t.Errorf("Builds.List returned error: %v", err)
	}

	want := []Build{{ID: "123"}, {ID: "1234"}}
	if diff := cmp.Diff(builds, want); diff != "" {
		t.Errorf("Builds.List diff: (-got +want)\n%s", diff)
	}
}

func TestBuildsService_List_with_exclude_pipeline(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/builds", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"exclude_pipeline": "true",
		})
		_, _ = fmt.Fprint(w, `[{"id":"123"},{"id":"1234"}]`)
	})

	opt := &BuildsListOptions{
		ExcludePipeline: true,
	}
	builds, _, err := client.Builds.List(context.Background(), opt)
	if err != nil {
		t.Errorf("Builds.List returned error: %v", err)
	}

	want := []Build{{ID: "123"}, {ID: "1234"}}
	if diff := cmp.Diff(builds, want); diff != "" {
		t.Errorf("Builds.List diff: (-got +want)\n%s", diff)
	}
}

func TestBuildsService_ListByPipeline(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-great-org/pipelines/sup-keith/builds", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = fmt.Fprint(w, `[{"id":"123"},{"id":"1234"}]`)
	})

	builds, _, err := client.Builds.ListByPipeline(context.Background(), "my-great-org", "sup-keith", nil)
	if err != nil {
		t.Errorf("Builds.List returned error: %v", err)
	}

	want := []Build{{ID: "123"}, {ID: "1234"}}
	if diff := cmp.Diff(builds, want); diff != "" {
		t.Errorf("Builds.List diff: (-got +want)\n%s", diff)
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

func TestBuildService_Create(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-great-org/pipelines/sup-keith/builds", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		_, _ = fmt.Fprint(w, `{"id":"123"}`)
	})

	cb := CreateBuild{
		Commit:  "HEAD",
		Branch:  "main",
		Message: "Hello, world!",
	}
	build, _, err := client.Builds.Create(context.Background(), "my-great-org", "sup-keith", cb)
	if err != nil {
		t.Fatalf("Builds.Create returned error: %v", err)
	}

	want := Build{ID: "123"}
	if diff := cmp.Diff(build, want); diff != "" {
		t.Fatalf("Builds.Create diff: (-got +want)\n%s", diff)
	}
}

func TestJSONUnMarshallBuild(t *testing.T) {
	t.Parallel()

	// a test table of build JSON strings and their expected build struct
	tests := []struct {
		name   string
		json   string
		expect Build
	}{
		{
			name: "basic build with author info",
			json: `{
				"id": "123",
				"state": "running",
				"blocked": false,
				"message": "Hello, world!",
				"commit": "HEAD",
				"branch": "main",
				"source": "ui",
				"author": {
					"username": "foojim",
					"name": "Uhh, Jim",
					"email": "foojim@example.com"
				}
		}`,
			expect: Build{
				ID:      "123",
				State:   "running",
				Blocked: false,
				Message: "Hello, world!",
				Commit:  "HEAD",
				Branch:  "main",
				Source:  "ui",
				Author: Author{
					Email:    "foojim@example.com",
					Name:     "Uhh, Jim",
					Username: "foojim",
				},
			},
		},
		{
			name: "basic build with author email string",
			json: `{
				"id": "123",
				"state": "running",
				"blocked": false,
				"message": "Hello, world!",
				"commit": "HEAD",
				"branch": "main",
				"source": "ui",
				"author": "foojim@example.com"
		}`,
			expect: Build{
				ID:      "123",
				State:   "running",
				Blocked: false,
				Message: "Hello, world!",
				Commit:  "HEAD",
				Branch:  "main",
				Source:  "ui",
				Author: Author{
					Email: "foojim@example.com",
				},
			},
		},
		{
			name: "build with test engine run ids",
			json: `{
				"id": "456",
				"state": "passed",
				"blocked": false,
				"message": "Test run",
				"commit": "abc123",
				"branch": "main",
				"source": "api",
				"test_engine": {
					"runs": [
						{
							"id": "2fd5c15a-a3a6-45ab-b1b2-6ea59616c1cf",
							"suite": {
								"id": "b1c349b6-c7d6-4633-8941-7e615459c03d",
								"slug": "banana"
							}
						}
					]
				}
		}`,
			expect: Build{
				ID:      "456",
				State:   "passed",
				Blocked: false,
				Message: "Test run",
				Commit:  "abc123",
				Branch:  "main",
				Source:  "api",
				TestEngine: &TestEngineProperty{
					Runs: []TestEngineRun{
						{
							ID: "2fd5c15a-a3a6-45ab-b1b2-6ea59616c1cf",
							Suite: TestEngineSuite{
								ID:   "b1c349b6-c7d6-4633-8941-7e615459c03d",
								Slug: "banana",
							},
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			var b Build
			if err := json.Unmarshal([]byte(test.json), &b); err != nil {
				t.Fatalf("could not unmarshal: %v", err)
			}

			if diff := cmp.Diff(b, test.expect); diff != "" {
				t.Fatalf("Build diff: (-got +want)\n%s", diff)
			}
		})
	}
}
