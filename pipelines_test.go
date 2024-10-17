package buildkite

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestPipelinesService_List(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-great-org/pipelines", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"id":"123"},{"id":"1234"}]`)
	})

	pipelines, _, err := client.Pipelines.List(context.Background(), "my-great-org", nil)
	if err != nil {
		t.Errorf("Pipelines.List returned error: %v", err)
	}

	want := []Pipeline{{ID: "123"}, {ID: "1234"}}
	if diff := cmp.Diff(pipelines, want); diff != "" {
		t.Errorf("Pipelines.List diff: (-got +want)\n%s", diff)
	}
}

func TestPipelinesService_Create(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	input := CreatePipeline{Name: "my-great-pipeline",
		Repository: "my-great-repo",
		Steps: []Step{
			{
				Type:    "script",
				Name:    "Build :package",
				Command: "script/release.sh",
				Plugins: Plugins{
					"my-org/docker#v3.3.0": {
						"image":   "node",
						"workdir": "/app",
					},
				},
			},
		},
		DefaultBranch: "main",
		Tags: []string{
			"well-tested",
			"great-config",
		},
	}

	server.HandleFunc("/v2/organizations/my-great-org/pipelines", func(w http.ResponseWriter, r *http.Request) {
		var v CreatePipeline
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("Error parsing json body: %v", err)
		}

		testMethod(t, r, "POST")

		if diff := cmp.Diff(v, input); diff != "" {
			t.Errorf("Request body diff: (-got +want)\n%s", diff)
		}

		fmt.Fprint(w, `{
						"name":"my-great-pipeline",
						"repository":"my-great-repo",
						"steps": [
							{
								"type": "script",
								"name": "Build :package:",
								"command": "script/release.sh",
								"plugins": {
									"my-org/docker#v3.3.0": {
										"image":   "node",
										"workdir": "/app"
									}
								}
							}
						],
						"default_branch":"main",
            "tags": [
              "well-tested",
              "great-config"
            ]
					}`)
	})

	pipeline, _, err := client.Pipelines.Create(context.Background(), "my-great-org", input)
	if err != nil {
		t.Errorf("Pipelines.Create returned error: %v", err)
	}

	want := Pipeline{
		Name:       "my-great-pipeline",
		Repository: "my-great-repo",
		Steps: []Step{
			{
				Type:    "script",
				Name:    "Build :package:",
				Command: "script/release.sh",
				Plugins: Plugins{
					"my-org/docker#v3.3.0": {
						"image":   "node",
						"workdir": "/app",
					},
				},
			},
		},
		DefaultBranch: "main",
		Tags: []string{
			"well-tested",
			"great-config",
		},
	}
	if diff := cmp.Diff(pipeline, want); diff != "" {
		t.Errorf("Pipelines.Create diff: (-got +want)\n%s", diff)
	}

}

func TestPipelinesService_CreateByConfiguration(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	input := CreatePipeline{Name: "my-great-pipeline",
		Repository:    "my-great-repo",
		Configuration: "steps:\n  - command: \"script/release.sh\"\n    label: \"Build :package:\"",
	}

	server.HandleFunc("/v2/organizations/my-great-org/pipelines", func(w http.ResponseWriter, r *http.Request) {
		var v CreatePipeline
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("Error parsing json body: %v", err)
		}

		testMethod(t, r, "POST")

		if diff := cmp.Diff(v, input); diff != "" {
			t.Errorf("Request body diff: (-got +want)\n%s", diff)
		}

		fmt.Fprint(w, `{
						"name":"my-great-pipeline",
						"repository":"my-great-repo",
						"configuration":"steps:\n  - command: \"script/release.sh\"\n    label: \"Build :package:\"",
						"steps": [
							{
								"type": "script",
								"name": "Build :package:",
								"command": "script/release.sh",
								"plugins": {
									"my-org/docker#v3.3.0": {
										"image":   "node",
										"workdir": "/app"
									}
								}
							}
						]
					}`)
	})

	pipeline, _, err := client.Pipelines.Create(context.Background(), "my-great-org", input)
	if err != nil {
		t.Errorf("Pipelines.Create returned error: %v", err)
	}

	want := Pipeline{
		Name:       "my-great-pipeline",
		Repository: "my-great-repo",
		Steps: []Step{
			{
				Type:    "script",
				Name:    "Build :package:",
				Command: "script/release.sh",
				Plugins: Plugins{
					"my-org/docker#v3.3.0": {
						"image":   "node",
						"workdir": "/app",
					},
				},
			},
		},
		Configuration: "steps:\n  - command: \"script/release.sh\"\n    label: \"Build :package:\"",
	}
	if diff := cmp.Diff(pipeline, want); diff != "" {
		t.Errorf("Pipelines.Create diff: (-got +want)\n%s", diff)
	}

}

func TestPipelinesService_Get(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-great-org/pipelines/my-great-pipeline-slug", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":"123",
						"slug":"my-great-pipeline-slug",
						"timeout_in_minutes": "1",
						"agent_query_rules": [
							"queue=default",
							"llamas=true"
						]}`)
	})

	pipeline, _, err := client.Pipelines.Get(context.Background(), "my-great-org", "my-great-pipeline-slug")
	if err != nil {
		t.Errorf("Pipelines.Get returned error: %v", err)
	}

	want := Pipeline{ID: "123", Slug: "my-great-pipeline-slug"}
	if diff := cmp.Diff(pipeline, want); diff != "" {
		t.Errorf("Pipelines.Get diff: (-got +want)\n%s", diff)
	}
}

func TestPipelinesService_Delete(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-great-org/pipelines/my-great-pipeline-slug", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	_, err := client.Pipelines.Delete(context.Background(), "my-great-org", "my-great-pipeline-slug")
	if err != nil {
		t.Errorf("Pipelines.Delete returned error: %v", err)
	}
}

func TestPipelinesService_Update(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	input := CreatePipeline{Name: "my-great-pipeline",
		Repository: "my-great-repo",
		Steps: []Step{
			{
				Type:    "script",
				Name:    "Build :package",
				Command: "script/release.sh",
				Plugins: Plugins{
					"my-org/docker#v3.3.0": {
						"image":   "node",
						"workdir": "/app",
					},
				},
			},
		},
	}

	server.HandleFunc("/v2/organizations/my-great-org/pipelines", func(w http.ResponseWriter, r *http.Request) {
		var v CreatePipeline
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("Error parsing json body: %v", err)
		}

		testMethod(t, r, "POST")

		if diff := cmp.Diff(v, input); diff != "" {
			t.Errorf("Request body diff: (-got +want)\n%s", diff)
		}

		fmt.Fprint(w, `{
						"name":"my-great-pipeline",
						"repository":"my-great-repo",
						"steps": [
							{
								"type": "script",
								"name": "Build :package:",
								"command": "script/release.sh",
								"plugins": {
									"my-org/docker#v3.3.0": {
										"image":   "node",
										"workdir": "/app"
									}
								}
							}
						],
						"slug": "my-great-repo"
					}`)
	})

	_, _, err := client.Pipelines.Create(context.Background(), "my-great-org", input)
	if err != nil {
		t.Errorf("Pipelines.Create returned error: %v", err)
	}

	server.HandleFunc("/v2/organizations/my-great-org/pipelines/my-great-repo", func(w http.ResponseWriter, r *http.Request) {
		var v UpdatePipeline
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("Error parsing json body: %v", err)
		}

		if diff := cmp.Diff(v, UpdatePipeline{Name: "derp"}); diff != "" {
			t.Errorf("Request body diff: (-got +want)\n%s", diff)
		}

		testMethod(t, r, "PATCH")

		fmt.Fprint(w, `{
						"name":"derp",
						"repository":"my-great-repo",
						"steps": [
							{
								"type": "script",
								"name": "Build :package:",
								"command": "script/release.sh",
								"plugins": {
									"my-org/docker#v3.3.0": {
										"image":   "node",
										"workdir": "/app"
									}
								}
							}
						],
						"slug": "my-great-repo",
            "visibility": "public",
            "tags": ["fresh-tag"]
					}`)
	})

	got, _, err := client.Pipelines.Update(context.Background(), "my-great-org", "my-great-repo", UpdatePipeline{Name: "derp"})
	if err != nil {
		t.Errorf("Pipelines.Update returned error: %v", err)
	}

	want := Pipeline{
		Name:       "derp",
		Repository: "my-great-repo",
		Steps: []Step{
			{
				Type:    "script",
				Name:    "Build :package:",
				Command: "script/release.sh",
				Plugins: Plugins{
					"my-org/docker#v3.3.0": {
						"image":   "node",
						"workdir": "/app",
					},
				},
			},
		},
		Slug:       "my-great-repo",
		Visibility: "public",
		Tags:       []string{"fresh-tag"},
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Pipelines.Update diff: (-got +want)\n%s", diff)
	}
}

func TestPipelinesService_AddWebhook(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-great-org/pipelines/my-great-pipeline-slug/webhook", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
	})

	_, err := client.Pipelines.AddWebhook(context.Background(), "my-great-org", "my-great-pipeline-slug")
	if err != nil {
		t.Errorf("Pipelines.AddWebhook returned error: %v", err)
	}
}

func TestPipelinesService_Archive(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-great-org/pipelines/my-great-pipeline-slug/archive", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
	})

	_, err := client.Pipelines.Archive(context.Background(), "my-great-org", "my-great-pipeline-slug")
	if err != nil {
		t.Errorf("Pipelines.Archive returned error: %v", err)
	}
}

func TestPipelinesService_Unarchive(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-great-org/pipelines/my-great-pipeline-slug/unarchive", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
	})

	_, err := client.Pipelines.Unarchive(context.Background(), "my-great-org", "my-great-pipeline-slug")
	if err != nil {
		t.Errorf("Pipelines.UnArchive returned error: %v", err)
	}
}

func TestPluginsUnmarshal(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name string
		json string
	}{
		{
			name: "as an array of plugins",
			json: `[{"my-org/docker#v3.3.0": {"image":   "node", "workdir": "/app"}}]`,
		},
		{
			name: "as a map[string]Plugin",
			json: `{"my-org/docker#v3.3.0": {"image":   "node", "workdir": "/app"}}`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var plugins Plugins
			err := json.Unmarshal([]byte(tc.json), &plugins)
			if err != nil {
				t.Fatalf("unmarshal: %v", err)
			}

			want := Plugins{
				"my-org/docker#v3.3.0": {
					"image":   "node",
					"workdir": "/app",
				},
			}
			if diff := cmp.Diff(plugins, want); diff != "" {
				t.Errorf("Plugins.Unmarshal diff: (-got +want)\n%s", diff)
			}
		})
	}
}
