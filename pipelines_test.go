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
		_, _ = fmt.Fprint(w, `[{"id":"123"},{"id":"1234"}]`)
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

func TestPipelinesService_List_by_name(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-great-org/pipelines", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"name": "my-pipeline",
		})
		_, _ = fmt.Fprint(w, `[{"id":"123","name":"my-pipeline"}]`)
	})

	opt := &PipelineListOptions{
		Name: "my-pipeline",
	}
	pipelines, _, err := client.Pipelines.List(context.Background(), "my-great-org", opt)
	if err != nil {
		t.Errorf("Pipelines.List returned error: %v", err)
	}

	want := []Pipeline{{ID: "123", Name: "my-pipeline"}}
	if diff := cmp.Diff(pipelines, want); diff != "" {
		t.Errorf("Pipelines.List diff: (-got +want)\n%s", diff)
	}
}

func TestPipelinesService_List_by_repository(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-great-org/pipelines", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"repository": "https://github.com/buildkite/test-repo",
		})
		_, _ = fmt.Fprint(w, `[{"id":"456","repository":"https://github.com/buildkite/test-repo"}]`)
	})

	opt := &PipelineListOptions{
		Repository: "https://github.com/buildkite/test-repo",
	}
	pipelines, _, err := client.Pipelines.List(context.Background(), "my-great-org", opt)
	if err != nil {
		t.Errorf("Pipelines.List returned error: %v", err)
	}

	want := []Pipeline{{ID: "456", Repository: "https://github.com/buildkite/test-repo"}}
	if diff := cmp.Diff(pipelines, want); diff != "" {
		t.Errorf("Pipelines.List diff: (-got +want)\n%s", diff)
	}
}

func TestPipelinesService_List_by_name_and_repository(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-great-org/pipelines", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"name":       "my-pipeline",
			"repository": "https://github.com/buildkite/test-repo",
		})
		_, _ = fmt.Fprint(w, `[{"id":"789","name":"my-pipeline","repository":"https://github.com/buildkite/test-repo"}]`)
	})

	opt := &PipelineListOptions{
		Name:       "my-pipeline",
		Repository: "https://github.com/buildkite/test-repo",
	}
	pipelines, _, err := client.Pipelines.List(context.Background(), "my-great-org", opt)
	if err != nil {
		t.Errorf("Pipelines.List returned error: %v", err)
	}

	want := []Pipeline{{ID: "789", Name: "my-pipeline", Repository: "https://github.com/buildkite/test-repo"}}
	if diff := cmp.Diff(pipelines, want); diff != "" {
		t.Errorf("Pipelines.List diff: (-got +want)\n%s", diff)
	}
}

func TestPipelinesService_Create(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	input := CreatePipeline{
		Name:       "my-great-pipeline",
		Repository: "my-great-repo",
		ClusterID:  "528000d8-4ee1-4479-8af1-032b143185f0",
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

		_, _ = fmt.Fprint(w, `{
						"name":"my-great-pipeline",
						"repository":"my-great-repo",
						"cluster_id":"528000d8-4ee1-4479-8af1-032b143185f0",
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
		ClusterID:  "528000d8-4ee1-4479-8af1-032b143185f0",
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

	input := CreatePipeline{
		Name:          "my-great-pipeline",
		Repository:    "my-great-repo",
		ClusterID:     "528000d8-4ee1-4479-8af1-032b143185f0",
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

		_, _ = fmt.Fprint(w, `{
						"name":"my-great-pipeline",
						"repository":"my-great-repo",
						"cluster_id":"528000d8-4ee1-4479-8af1-032b143185f0",
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
		ClusterID:  "528000d8-4ee1-4479-8af1-032b143185f0",
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
		_, _ = fmt.Fprint(w, `{"id":"123",
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

	input := CreatePipeline{
		Name:       "my-great-pipeline",
		Repository: "my-great-repo",
		ClusterID:  "528000d8-4ee1-4479-8af1-032b143185f0",
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

		_, _ = fmt.Fprint(w, `{
						"name":"my-great-pipeline",
						"repository":"my-great-repo",
						"cluster_id":"528000d8-4ee1-4479-8af1-032b143185f0",
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
		testMethod(t, r, "PATCH")
		assertRequestJSON(t, r, `{"name":"derp"}`)

		_, _ = fmt.Fprint(w, `{
						"name":"derp",
						"repository":"my-great-repo",
						"cluster_id":"528000d8-4ee1-4479-8af1-032b143185f0",
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

	got, _, err := client.Pipelines.Update(context.Background(), "my-great-org", "my-great-repo", UpdatePipeline{Name: Some("derp")})
	if err != nil {
		t.Errorf("Pipelines.Update returned error: %v", err)
	}

	want := Pipeline{
		Name:       "derp",
		Repository: "my-great-repo",
		ClusterID:  "528000d8-4ee1-4479-8af1-032b143185f0",
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

func TestPipelinesService_UpdateClearsTagsAndSendsFalse(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/my-great-org/pipelines/my-great-repo", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		assertRequestJSON(t, r, `{
			"tags": [],
			"skip_queued_branch_builds": false
		}`)

		_, _ = fmt.Fprint(w, `{
			"name": "my-great-repo",
			"repository": "git@github.com:buildkite/my-great-repo.git",
			"skip_queued_branch_builds": false,
			"cancel_running_branch_builds": true,
			"tags": []
		}`)
	})

	got, _, err := client.Pipelines.Update(context.Background(), "my-great-org", "my-great-repo", UpdatePipeline{
		Tags:                   Some([]string{}),
		SkipQueuedBranchBuilds: Some(false),
	})
	if err != nil {
		t.Errorf("Pipelines.Update returned error: %v", err)
	}

	if len(got.Tags) != 0 {
		t.Errorf("Update tags length = %d, want 0", len(got.Tags))
	}
	if got.SkipQueuedBranchBuilds {
		t.Error("Update skip_queued_branch_builds = true, want false")
	}
}

func TestPipelinesService_CreateGitHubActions(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/example/pipelines", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		assertRequestJSON(t, r, `{
			"name":"Workflows", "repository":"https://github.com/example/workflows.git",
			"cluster_id":"example-cluster", "default_branch":"main",
			"configuration":"steps:\n  - plugin: github-actions\n",
			"skip_queued_branch_builds":false, "cancel_running_branch_builds":false,
			"provider_settings":{
				"github_actions_mode":true, "github_workflow_access_tokens_enabled":true,
				"publish_commit_status":false, "build_branches":true,
				"build_pull_requests":true, "build_tags":true,
				"filter_enabled":false, "filter_condition":null,
				"pull_request_branch_filter_enabled":false,
				"pull_request_branch_filter_configuration":null,
				"ignore_default_branch_pull_requests":false,
				"skip_builds_for_closed_pull_requests":false,
				"skip_builds_for_existing_commits":false,
				"skip_pull_request_builds_for_existing_commits":false,
				"trigger_mode":"none"
			}
		}`)
		w.WriteHeader(http.StatusCreated)
		_, _ = fmt.Fprint(w, `{"slug":"workflows","branch_configuration":null,"provider":{"id":"github","settings":{
			"github_actions_mode":true,"github_workflow_access_tokens_enabled":true,
			"ignore_default_branch_pull_requests":false,"skip_builds_for_closed_pull_requests":false,
			"skip_builds_for_existing_commits":false,"build_branches":true,"build_pull_requests":true,
			"build_tags":true,"trigger_mode":"none","filter_condition":null,
			"pull_request_branch_filter_configuration":null
		}}}`)
	})

	got, _, err := client.Pipelines.Create(context.Background(), "example", CreatePipeline{
		Name:          "Workflows",
		Repository:    "https://github.com/example/workflows.git",
		ClusterID:     "example-cluster",
		DefaultBranch: "main",
		Configuration: "steps:\n  - plugin: github-actions\n",
		ProviderSettings: &GitHubSettingsInput{
			GitHubActionsMode:                       Some(true),
			GitHubWorkflowAccessTokensEnabled:       Some(true),
			PublishCommitStatus:                     Some(false),
			BuildBranches:                           Some(true),
			BuildPullRequests:                       Some(true),
			BuildTags:                               Some(true),
			FilterEnabled:                           Some(false),
			FilterCondition:                         Null[string](),
			PullRequestBranchFilterEnabled:          Some(false),
			PullRequestBranchFilterConfiguration:    Null[string](),
			IgnoreDefaultBranchPullRequests:         Some(false),
			SkipBuildsForClosedPullRequests:         Some(false),
			SkipBuildsForExistingCommits:            Some(false),
			SkipPullRequestBuildsForExistingCommits: Some(false),
			TriggerMode:                             Some("none"),
		},
	})
	if err != nil {
		t.Fatalf("Pipelines.Create: %v", err)
	}
	if diff := cmp.Diff(Pipeline{
		Slug: "workflows",
		Provider: Provider{ID: "github", Settings: &GitHubSettings{
			GitHubActionsMode:                 boolPtr(true),
			GitHubWorkflowAccessTokensEnabled: boolPtr(true),
			IgnoreDefaultBranchPullRequests:   boolPtr(false),
			SkipBuildsForClosedPullRequests:   boolPtr(false),
			SkipBuildsForExistingCommits:      boolPtr(false),
			BuildBranches:                     true,
			BuildPullRequests:                 true,
			BuildTags:                         true,
			TriggerMode:                       "none",
		}},
	}, got); diff != "" {
		t.Errorf("Create response (-want +got):\n%s", diff)
	}
}

func TestPipelinesService_UpdateGitHubSettings(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   UpdatePipeline
		want string
	}{
		{
			name: "enable mode without resetting other settings",
			in: UpdatePipeline{ProviderSettings: Some[ProviderSettings](&GitHubSettingsInput{
				GitHubActionsMode: Some(true),
			})},
			want: `{"provider_settings":{"github_actions_mode":true}}`,
		},
		{
			name: "disable mode and tokens explicitly",
			in: UpdatePipeline{ProviderSettings: Some[ProviderSettings](&GitHubSettingsInput{
				GitHubActionsMode:                 Some(false),
				GitHubWorkflowAccessTokensEnabled: Some(false),
			})},
			want: `{"provider_settings":{"github_actions_mode":false,"github_workflow_access_tokens_enabled":false}}`,
		},
		{
			name: "clear migration filters and disable intermediate builds",
			in: UpdatePipeline{
				BranchConfiguration:       Null[string](),
				SkipQueuedBranchBuilds:    Some(false),
				CancelRunningBranchBuilds: Some(false),
				ProviderSettings: Some[ProviderSettings](&GitHubSettingsInput{
					FilterEnabled:                        Some(false),
					FilterCondition:                      Null[string](),
					PullRequestBranchFilterEnabled:       Some(false),
					PullRequestBranchFilterConfiguration: Null[string](),
				}),
			},
			want: `{"branch_configuration":null,"skip_queued_branch_builds":false,"cancel_running_branch_builds":false,
				"provider_settings":{"filter_enabled":false,"filter_condition":null,
				"pull_request_branch_filter_enabled":false,"pull_request_branch_filter_configuration":null}}`,
		},
		{
			name: "empty strings remain distinct from null",
			in: UpdatePipeline{
				BranchConfiguration: Some(""),
				ProviderSettings: Some[ProviderSettings](&GitHubSettingsInput{
					FilterCondition:                      Some(""),
					PullRequestBranchFilterConfiguration: Some(""),
				}),
			},
			want: `{"branch_configuration":"","provider_settings":{"filter_condition":"","pull_request_branch_filter_configuration":""}}`,
		},
		{
			name: "ordinary pipeline filters",
			in: UpdatePipeline{
				BranchConfiguration: Some("main release/*"),
				ProviderSettings: Some[ProviderSettings](&GitHubSettingsInput{
					TriggerMode:                          Some("code"),
					PublishCommitStatus:                  Some(true),
					FilterEnabled:                        Some(true),
					FilterCondition:                      Some("build.message !~ /skip/"),
					PullRequestBranchFilterEnabled:       Some(true),
					PullRequestBranchFilterConfiguration: Some("feature/*"),
				}),
			},
			want: `{"branch_configuration":"main release/*","provider_settings":{"trigger_mode":"code",
				"publish_commit_status":true,"filter_enabled":true,"filter_condition":"build.message !~ /skip/",
				"pull_request_branch_filter_enabled":true,"pull_request_branch_filter_configuration":"feature/*"}}`,
		},
		{
			name: "ordinary update omits provider and filters",
			in:   UpdatePipeline{Description: Some("Ordinary pipeline")},
			want: `{"description":"Ordinary pipeline"}`,
		},
		{
			name: "legacy GitHub writes keep existing booleans",
			in: UpdatePipeline{ProviderSettings: Some[ProviderSettings](&GitHubSettings{
				BuildBranches:       true,
				PublishCommitStatus: true,
			})},
			want: `{"provider_settings":{
				"build_branches":true,"publish_commit_status":true,
				"build_pull_requests":false,"pull_request_branch_filter_enabled":false,
				"skip_pull_request_builds_for_existing_commits":false,"build_pull_request_forks":false,
				"prefix_pull_request_fork_branch_names":false,"build_tags":false,
				"publish_commit_status_per_step":false,"filter_enabled":false,
				"separate_pull_request_statuses":false,"publish_blocked_as_pending":false
			}}`,
		},
		{
			name: "Bitbucket writes stay unchanged",
			in: UpdatePipeline{ProviderSettings: Some[ProviderSettings](&BitbucketSettings{
				BuildTags: true,
			})},
			want: `{"provider_settings":{
				"build_pull_requests":false,"build_branches":false,"pull_request_branch_filter_enabled":false,
				"skip_pull_request_builds_for_existing_commits":false,"build_tags":true,
				"publish_commit_status":false,"publish_commit_status_per_step":false
			}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server, client, teardown := newMockServerAndClient(t)
			t.Cleanup(teardown)
			server.HandleFunc("/v2/organizations/example/pipelines/workflows", func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "PATCH")
				assertRequestJSON(t, r, tt.want)
				_, _ = fmt.Fprint(w, `{"slug":"workflows"}`)
			})
			got, _, err := client.Pipelines.Update(context.Background(), "example", "workflows", tt.in)
			if err != nil {
				t.Fatalf("Pipelines.Update: %v", err)
			}
			if got.Slug != "workflows" {
				t.Errorf("Update slug = %q, want workflows", got.Slug)
			}
		})
	}
}

func TestPipelinesService_GetGitHubActionsSettings(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name    string
		present bool
		enabled bool
	}{
		{name: "absent"},
		{name: "mode enabled", present: true, enabled: true},
		{name: "mode disabled", present: true, enabled: false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server, client, teardown := newMockServerAndClient(t)
			t.Cleanup(teardown)
			server.HandleFunc("/v2/organizations/example/pipelines/workflows", func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "GET")
				fields := ""
				if tt.present {
					fields = fmt.Sprintf(`"github_actions_mode":%t,"github_workflow_access_tokens_enabled":%t,
						"ignore_default_branch_pull_requests":%t,"skip_builds_for_closed_pull_requests":%t,
						"skip_builds_for_existing_commits":%t,`, tt.enabled, !tt.enabled, tt.enabled, !tt.enabled, tt.enabled)
				}
				_, _ = fmt.Fprintf(w, `{"branch_configuration":"release/*","provider":{"id":"github","settings":{
					%s"publish_commit_status":true,"trigger_mode":"code",
					"filter_enabled":true,"filter_condition":"build.branch == 'main'",
					"pull_request_branch_filter_enabled":true,"pull_request_branch_filter_configuration":"feature/*"
				}}}`, fields)
			})
			got, _, err := client.Pipelines.Get(context.Background(), "example", "workflows")
			if err != nil {
				t.Fatalf("Pipelines.Get: %v", err)
			}
			want := Pipeline{
				BranchConfiguration: "release/*",
				Provider: Provider{ID: "github", Settings: &GitHubSettings{
					PublishCommitStatus:                  true,
					TriggerMode:                          "code",
					FilterEnabled:                        true,
					FilterCondition:                      "build.branch == 'main'",
					PullRequestBranchFilterEnabled:       true,
					PullRequestBranchFilterConfiguration: "feature/*",
				}},
			}
			if tt.present {
				settings := want.Provider.Settings.(*GitHubSettings)
				settings.GitHubActionsMode = boolPtr(tt.enabled)
				settings.GitHubWorkflowAccessTokensEnabled = boolPtr(!tt.enabled)
				settings.IgnoreDefaultBranchPullRequests = boolPtr(tt.enabled)
				settings.SkipBuildsForClosedPullRequests = boolPtr(!tt.enabled)
				settings.SkipBuildsForExistingCommits = boolPtr(tt.enabled)
			}
			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("Get response (-want +got):\n%s", diff)
			}
			encoded, err := json.Marshal(got)
			if err != nil {
				t.Fatalf("Marshal pipeline: %v", err)
			}
			var wire struct {
				Provider struct {
					Settings map[string]json.RawMessage
				}
			}
			if err := json.Unmarshal(encoded, &wire); err != nil {
				t.Fatalf("Unmarshal pipeline JSON fields: %v", err)
			}
			for _, key := range []string{"github_actions_mode", "github_workflow_access_tokens_enabled", "ignore_default_branch_pull_requests", "skip_builds_for_closed_pull_requests", "skip_builds_for_existing_commits"} {
				if _, present := wire.Provider.Settings[key]; present != tt.present {
					t.Errorf("Pipeline JSON field %s present = %t, want %t", key, present, tt.present)
				}
			}
			var roundTrip Pipeline
			if err := json.Unmarshal(encoded, &roundTrip); err != nil {
				t.Fatalf("Unmarshal pipeline: %v", err)
			}
			if diff := cmp.Diff(want, roundTrip); diff != "" {
				t.Errorf("Pipeline JSON round trip (-want +got):\n%s", diff)
			}
		})
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

func TestPipelinesService_CreateWithTeams(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	input := CreatePipeline{
		Name:          "my-great-pipeline",
		Repository:    "my-great-repo",
		ClusterID:     "528000d8-4ee1-4479-8af1-032b143185f0",
		Configuration: "steps:\n  - command: \"script/release.sh\"",
		Teams: map[string]string{
			"b3a1c8f4-1f2e-4f0e-9d8c-0a1b2c3d4e5f": "build_and_read",
		},
	}

	server.HandleFunc("/v2/organizations/my-great-org/pipelines", func(w http.ResponseWriter, r *http.Request) {
		var v map[string]any
		if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
			t.Fatalf("Error parsing json body: %v", err)
		}

		testMethod(t, r, "POST")

		want := map[string]any{"b3a1c8f4-1f2e-4f0e-9d8c-0a1b2c3d4e5f": "build_and_read"}
		if diff := cmp.Diff(v["teams"], want); diff != "" {
			t.Errorf("teams diff: (-got +want)\n%s", diff)
		}
		if _, ok := v["team_uuids"]; ok {
			t.Error("team_uuids should be omitted when unset")
		}

		_, _ = fmt.Fprint(w, `{"name":"my-great-pipeline","repository":"my-great-repo"}`)
	})

	if _, _, err := client.Pipelines.Create(context.Background(), "my-great-org", input); err != nil {
		t.Errorf("Pipelines.Create returned error: %v", err)
	}
}
