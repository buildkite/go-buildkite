package buildkite

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestPipelinesService_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/organizations/my-great-org/pipelines", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"id":"123"},{"id":"1234"}]`)
	})

	pipelines, _, err := client.Pipelines.List("my-great-org", nil)
	if err != nil {
		t.Errorf("Pipelines.List returned error: %v", err)
	}

	want := []Pipeline{{ID: String("123")}, {ID: String("1234")}}
	if !reflect.DeepEqual(pipelines, want) {
		t.Errorf("Pipelines.List returned %+v, want %+v", pipelines, want)
	}
}

func TestPipelinesService_Create(t *testing.T) {
	setup()
	defer teardown()

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
	}

	mux.HandleFunc("/v2/organizations/my-great-org/pipelines", func(w http.ResponseWriter, r *http.Request) {
		v := new(CreatePipeline)
		json.NewDecoder(r.Body).Decode(&v)

		testMethod(t, r, "POST")

		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
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
						"default_branch":"main"
					}`)
	})

	pipeline, _, err := client.Pipelines.Create("my-great-org", input)
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
	}
	if !reflect.DeepEqual(pipeline, want) {
		t.Errorf("Pipelines.Create returned %+v, want %+v", pipeline, want)
	}

}

func TestPipelinesService_CreateByConfiguration(t *testing.T) {
	setup()
	defer teardown()

	input := &CreatePipeline{Name: *String("my-great-pipeline"),
		Repository:    *String("my-great-repo"),
		Configuration: *String("steps:\n  - command: \"script/release.sh\"\n    label: \"Build :package:\""),
	}

	mux.HandleFunc("/v2/organizations/my-great-org/pipelines", func(w http.ResponseWriter, r *http.Request) {
		v := new(CreatePipeline)
		json.NewDecoder(r.Body).Decode(&v)

		testMethod(t, r, "POST")

		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
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

	pipeline, _, err := client.Pipelines.Create("my-great-org", input)
	if err != nil {
		t.Errorf("Pipelines.Create returned error: %v", err)
	}

	want := &Pipeline{Name: String("my-great-pipeline"),
		Repository: String("my-great-repo"),
		Steps: []*Step{
			&Step{
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
	if !reflect.DeepEqual(pipeline, want) {
		t.Errorf("Pipelines.Create returned %+v, want %+v", pipeline.Configuration, want.Configuration)
	}

}

func TestPipelinesService_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/organizations/my-great-org/pipelines/my-great-pipeline-slug", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"id":"123",
						"slug":"my-great-pipeline-slug",
						"timeout_in_minutes": "1",
						"agent_query_rules": [
							"queue=default",
							"llamas=true"
						]}`)
	})

	pipeline, _, err := client.Pipelines.Get("my-great-org", "my-great-pipeline-slug")
	if err != nil {
		t.Errorf("Pipelines.Get returned error: %v", err)
	}

	want := &Pipeline{ID: String("123"), Slug: String("my-great-pipeline-slug")}
	if !reflect.DeepEqual(pipeline, want) {
		t.Errorf("Pipelines.Get returned %+v, want %+v", pipeline, want)
	}
}

func TestPipelinesService_Delete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/organizations/my-great-org/pipelines/my-great-pipeline-slug", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	_, err := client.Pipelines.Delete("my-great-org", "my-great-pipeline-slug")
	if err != nil {
		t.Errorf("Pipelines.Delete returned error: %v", err)
	}
}

func TestPipelinesService_Update(t *testing.T) {
	setup()
	defer teardown()

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

	mux.HandleFunc("/v2/organizations/my-great-org/pipelines", func(w http.ResponseWriter, r *http.Request) {
		v := new(CreatePipeline)
		json.NewDecoder(r.Body).Decode(&v)

		testMethod(t, r, "POST")

		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
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

	pipeline, _, err := client.Pipelines.Create("my-great-org", input)
	if err != nil {
		t.Errorf("Pipelines.Create returned error: %v", err)
	}

	pipeline.Name = String("derp")

	mux.HandleFunc("/v2/organizations/my-great-org/pipelines/my-great-repo", func(w http.ResponseWriter, r *http.Request) {
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
						"slug": "my-great-repo"
					}`)
	})

	_, err = client.Pipelines.Update("my-great-org", pipeline)
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
		Slug: String("my-great-repo"),
	}

	if !reflect.DeepEqual(pipeline, want) {
		t.Errorf("Pipelines.Update returned %+v, want %+v", pipeline, want)
	}
}

func TestPipelinesService_AddWebhook(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/organizations/my-great-org/pipelines/my-great-pipeline-slug/webhook", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
	})

	_, err := client.Pipelines.AddWebhook("my-great-org", "my-great-pipeline-slug")
	if err != nil {
		t.Errorf("Pipelines.AddWebhook returned error: %v", err)
	}
}

func TestPipelinesService_Archive(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/organizations/my-great-org/pipelines/my-great-pipeline-slug/archive", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
	})

	_, err := client.Pipelines.Archive("my-great-org", "my-great-pipeline-slug")
	if err != nil {
		t.Errorf("Pipelines.Archive returned error: %v", err)
	}
}

func TestPipelinesService_Unarchive(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v2/organizations/my-great-org/pipelines/my-great-pipeline-slug/unarchive", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
	})

	_, err := client.Pipelines.Unarchive("my-great-org", "my-great-pipeline-slug")
	if err != nil {
		t.Errorf("Pipelines.UnArchive returned error: %v", err)
	}
}
