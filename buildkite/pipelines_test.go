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

	want := []Pipeline{{ID: String("123")}, {ID: String("1234")}}
	if diff := cmp.Diff(pipelines, want); diff != "" {
		t.Errorf("Pipelines.List diff: (-got +want)\n%s", diff)
	}
}

func TestPipelinesService_Create(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	input := &CreatePipeline{Name: *String("my-great-pipeline"),
		Repository: *String("my-great-repo"),
		Steps: []Step{
			{
				Type:    String("script"),
				Name:    String("Build :package"),
				Command: String("script/release.sh"),
				Plugins: Plugins{
					"my-org/docker#v3.3.0": {
						"image":   "node",
						"workdir": "/app",
					},
				},
			},
		},
		DefaultBranch: *String("main"),
		Tags: []string{
			"well-tested",
			"great-config",
		},
	}

	server.HandleFunc("/v2/organizations/my-great-org/pipelines", func(w http.ResponseWriter, r *http.Request) {
		v := new(CreatePipeline)
		json.NewDecoder(r.Body).Decode(&v)

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

	want := &Pipeline{Name: String("my-great-pipeline"),
		Repository: String("my-great-repo"),
		Steps: []*Step{
			{
				Type:    String("script"),
				Name:    String("Build :package:"),
				Command: String("script/release.sh"),
				Plugins: Plugins{
					"my-org/docker#v3.3.0": {
						"image":   "node",
						"workdir": "/app",
					},
				},
			},
		},
		DefaultBranch: String("main"),
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

	input := &CreatePipeline{Name: *String("my-great-pipeline"),
		Repository:    *String("my-great-repo"),
		Configuration: *String("steps:\n  - command: \"script/release.sh\"\n    label: \"Build :package:\""),
	}

	server.HandleFunc("/v2/organizations/my-great-org/pipelines", func(w http.ResponseWriter, r *http.Request) {
		v := new(CreatePipeline)
		json.NewDecoder(r.Body).Decode(&v)

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

	want := &Pipeline{Name: String("my-great-pipeline"),
		Repository: String("my-great-repo"),
		Steps: []*Step{
			{
				Type:    String("script"),
				Name:    String("Build :package:"),
				Command: String("script/release.sh"),
				Plugins: Plugins{
					"my-org/docker#v3.3.0": {
						"image":   "node",
						"workdir": "/app",
					},
				},
			},
		},
		Configuration: *String("steps:\n  - command: \"script/release.sh\"\n    label: \"Build :package:\""),
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

	want := &Pipeline{ID: String("123"), Slug: String("my-great-pipeline-slug")}
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

	input := &CreatePipeline{Name: *String("my-great-pipeline"),
		Repository: *String("my-great-repo"),
		Steps: []Step{
			{
				Type:    String("script"),
				Name:    String("Build :package"),
				Command: String("script/release.sh"),
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
		v := new(CreatePipeline)
		json.NewDecoder(r.Body).Decode(&v)

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

	pipeline, _, err := client.Pipelines.Create(context.Background(), "my-great-org", input)
	if err != nil {
		t.Errorf("Pipelines.Create returned error: %v", err)
	}

	pipeline.Name = String("derp")

	server.HandleFunc("/v2/organizations/my-great-org/pipelines/my-great-repo", func(w http.ResponseWriter, r *http.Request) {
		v := new(CreatePipeline)
		json.NewDecoder(r.Body).Decode(&v)

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
            "tags": [
              "fresh-tag"
            ]
					}`)
	})

	_, err = client.Pipelines.Update(context.Background(), "my-great-org", pipeline)
	if err != nil {
		t.Errorf("Pipelines.Update returned error: %v", err)
	}

	want := &Pipeline{Name: String("derp"),
		Repository: String("my-great-repo"),
		Steps: []*Step{
			{
				Type:    String("script"),
				Name:    String("Build :package:"),
				Command: String("script/release.sh"),
				Plugins: Plugins{
					"my-org/docker#v3.3.0": {
						"image":   "node",
						"workdir": "/app",
					},
				},
			},
		},
		Slug:       String("my-great-repo"),
		Visibility: String("public"),
		Tags:       []string{"fresh-tag"},
	}

	if diff := cmp.Diff(pipeline, want); diff != "" {
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
