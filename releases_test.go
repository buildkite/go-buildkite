package buildkite

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

const releaseTestID = "019a1dd8-a49b-70c5-99b8-b8f9f0ed6824"
const releaseTestDigest = "sha256:0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"

// Fixtures follow the API contract, including required nulls and omitted fields.
const releaseTargetJSON = `{"id":"target-id","graphql_id":"opaque-target","name":"east","labels":{"region":"au","empty":"","space":" "},"url":"https://api.example/target","id_url":"https://api.example/target-id","created_at":"2026-09-20T01:02:03Z","archived_at":"2026-09-21T01:02:03Z"}`
const releaseComponentJSON = `{"id":"component-id","graphql_id":"opaque-component","name":"api","artifact":{"id":"artifact-id","graphql_id":"opaque-artifact","name":"service"},"url":"https://api.example/component","id_url":"https://api.example/component-id","created_at":"2026-09-20T01:02:03Z","archived_at":null}`
const releaseInstanceJSON = `{"id":"instance-id","graphql_id":"opaque-instance","name":"api-blue","target_id":"target-id","component_id":"component-id","url":"https://api.example/instance","id_url":"https://api.example/instance-id","created_at":"2026-09-20T01:02:03Z","archived_at":null}`
const releaseArtifactJSON = `{"id":"artifact-id","graphql_id":"opaque-artifact","name":"service","url":"https://api.example/artifact","id_url":"https://api.example/artifact-id","created_at":"2026-09-20T01:02:03Z","archived_at":null}`
const releaseVersionJSON = `{"id":"digest-id","graphql_id":"opaque-digest","digest":"` + releaseTestDigest + `","source_url":null,"source_build":null,"pull_requests":[{"url":"https://github.com/acme/api/pull/12","provider":"github","repository_url":"https://github.com/acme/api","id":"12"},{"url":"https://example.com/pr/2"}],"created_by":null,"url":"https://api.example/digest","id_url":"https://api.example/digest-id","created_at":"2026-09-20T01:02:03Z","yanked":false,"yanked_at":null,"yank_reason":null}`
const releaseStatusJSON = `{"target_id":"target-id","instance":"api-blue","instance_id":"instance-id","present":false,"active_digest":null,"last_active_digest":"` + releaseTestDigest + `","first_seen_at":"2026-09-20T01:02:03Z","last_seen_at":"2026-09-21T01:02:03Z","sequence":9007199254740993,"catalog_resolution":{"status":"unknown","reason":"instance_absent"}}`
const releaseObservationJSON = `{"id":"observation-id","graphql_id":"opaque-observation","target_id":"target-id","sequence":9007199254740993,"event_id":"` + releaseTestID + `","type":"ComponentInstanceVersionObserved","instance":"api-blue","digest":"` + releaseTestDigest + `","occurred_at":"2026-09-20T01:02:03Z","recorded_at":"2026-09-20T01:03:04Z","catalog_resolution":{"status":"resolved","reason":null}}`

func releaseTestPage(item string) string {
	return `{"items":[` + item + `],"links":{"self":"https://api.example/self","first":"https://api.example/first","prev":"https://api.example/prev","next":"https://api.example/next?after=opaque%2Bcursor"}}`
}

func TestReleasesService_HTTPContract(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	target := ReleaseTargetRef{ID: releaseTestID}
	component := ReleaseComponentRef{Name: "api"}
	instance := ReleaseComponentInstanceRef{Target: target, Name: "api-blue"}
	artifact := ReleaseArtifactRef{Name: "service"}
	pagination := ReleasesListOptions{After: "opaque+cursor/=", PerPage: 25}
	pageQuery := "after=opaque%2Bcursor%2F%3D&per_page=25"
	targetPath := "targets/by-id/" + releaseTestID
	instanceOptions := &ReleasesComponentInstanceListOptions{ReleasesListOptions: pagination, Label: ReleaseTargetLabelFilters{"region": "au +", "empty": ""}}
	tests := []struct {
		name, method, path, query, request, response string
		status                                       int
		call                                         func(*ReleasesService) (any, *Response, error)
	}{
		{"ListTargets", "GET", "targets", pageQuery, "", releaseTestPage(releaseTargetJSON), 200, func(s *ReleasesService) (any, *Response, error) { return s.ListTargets(ctx, "acme", &pagination) }},
		{"GetTarget", "GET", "targets/" + releaseTestID, "", "", releaseTargetJSON, 200, func(s *ReleasesService) (any, *Response, error) {
			return s.GetTarget(ctx, "acme", ReleaseTargetRef{Name: releaseTestID})
		}},
		{"ListComponents", "GET", "components", pageQuery, "", releaseTestPage(releaseComponentJSON), 200, func(s *ReleasesService) (any, *Response, error) { return s.ListComponents(ctx, "acme", &pagination) }},
		{"GetComponent", "GET", "components/api", "", "", releaseComponentJSON, 200, func(s *ReleasesService) (any, *Response, error) { return s.GetComponent(ctx, "acme", component) }},
		{"GetComponentByID", "GET", "components/by-id/" + releaseTestID, "", "", releaseComponentJSON, 200, func(s *ReleasesService) (any, *Response, error) {
			return s.GetComponent(ctx, "acme", ReleaseComponentRef{ID: releaseTestID})
		}},
		{"ListComponentInstances", "GET", targetPath + "/instances", pageQuery + "&label[region]=au+%2B&label[empty]=", "", releaseTestPage(releaseInstanceJSON), 200, func(s *ReleasesService) (any, *Response, error) {
			return s.ListComponentInstances(ctx, "acme", target, instanceOptions)
		}},
		{"ListOrganizationComponentInstances", "GET", "component-instances", pageQuery + "&label[region]=au+%2B&label[empty]=", "", releaseTestPage(releaseInstanceJSON), 200, func(s *ReleasesService) (any, *Response, error) {
			return s.ListOrganizationComponentInstances(ctx, "acme", instanceOptions)
		}},
		{"GetComponentInstance", "GET", targetPath + "/instances/api-blue", "", "", releaseInstanceJSON, 200, func(s *ReleasesService) (any, *Response, error) { return s.GetComponentInstance(ctx, "acme", instance) }},
		{"GetComponentInstanceByID", "GET", "component-instances/by-id/" + releaseTestID, "", "", releaseInstanceJSON, 200, func(s *ReleasesService) (any, *Response, error) {
			return s.GetComponentInstance(ctx, "acme", ReleaseComponentInstanceRef{ID: releaseTestID})
		}},
		{"GetComponentInstanceByTargetName", "GET", "targets/east/instances/by-label", "", "", releaseInstanceJSON, 200, func(s *ReleasesService) (any, *Response, error) {
			return s.GetComponentInstance(ctx, "acme", ReleaseComponentInstanceRef{Target: ReleaseTargetRef{Name: "east"}, Name: "by-label"})
		}},
		{"ListStatuses", "GET", targetPath + "/instance-statuses", pageQuery, "", releaseTestPage(releaseStatusJSON), 200, func(s *ReleasesService) (any, *Response, error) {
			return s.ListComponentInstanceStatuses(ctx, "acme", target, &pagination)
		}},
		{"GetStatus", "GET", targetPath + "/instance-statuses/api-blue", "", "", releaseStatusJSON, 200, func(s *ReleasesService) (any, *Response, error) {
			return s.GetComponentInstanceStatus(ctx, "acme", target, "api-blue")
		}},
		{"ListObservations", "GET", targetPath + "/observations", "before=prior&instance=api-blue&type=ComponentInstanceVersionObserved", "", releaseTestPage(releaseObservationJSON), 200, func(s *ReleasesService) (any, *Response, error) {
			return s.ListObservations(ctx, "acme", target, &ReleasesObservationListOptions{ReleasesListOptions: ReleasesListOptions{Before: "prior"}, Instance: "api-blue", Type: "ComponentInstanceVersionObserved"})
		}},
		{"ListArtifacts", "GET", "artifacts", pageQuery, "", releaseTestPage(releaseArtifactJSON), 200, func(s *ReleasesService) (any, *Response, error) { return s.ListArtifacts(ctx, "acme", &pagination) }},
		{"GetArtifact", "GET", "artifacts/service", "", "", releaseArtifactJSON, 200, func(s *ReleasesService) (any, *Response, error) { return s.GetArtifact(ctx, "acme", artifact) }},
		{"ListVersions", "GET", "artifacts/by-id/" + releaseTestID + "/versions", pageQuery, "", releaseTestPage(releaseVersionJSON), 200, func(s *ReleasesService) (any, *Response, error) {
			return s.ListVersions(ctx, "acme", ReleaseArtifactRef{ID: releaseTestID}, &pagination)
		}},
		{"GetVersionByID", "GET", "artifact-versions/by-id/" + releaseTestID, "", "", releaseVersionJSON, 200, func(s *ReleasesService) (any, *Response, error) {
			return s.GetVersion(ctx, "acme", ReleaseVersionRef{ID: releaseTestID})
		}},
		{"CreateTarget", "POST", "targets", "", `{"name":"east","labels":{"region":"au"}}`, releaseTargetJSON, 201, func(s *ReleasesService) (any, *Response, error) {
			return s.CreateTarget(ctx, "acme", ReleaseTargetCreate{Name: "east", Labels: map[string]string{"region": "au"}})
		}},
		{"UpdateTarget", "PATCH", targetPath, "", `{"labels":{}}`, releaseTargetJSON, 200, func(s *ReleasesService) (any, *Response, error) {
			return s.UpdateTarget(ctx, "acme", target, ReleaseTargetUpdate{Labels: Some(map[string]string{})})
		}},
		{"CreateComponent", "POST", "components", "", `{"name":"api","artifact":{"id":"artifact-id"}}`, releaseComponentJSON, 201, func(s *ReleasesService) (any, *Response, error) {
			return s.CreateComponent(ctx, "acme", ReleaseComponentCreate{Name: "api", Artifact: ReleaseArtifactRef{ID: "artifact-id"}})
		}},
		{"CreateComponentInstance", "POST", targetPath + "/instances", "", `{"component_name":"api","name":"api-blue"}`, releaseInstanceJSON, 200, func(s *ReleasesService) (any, *Response, error) {
			return s.CreateComponentInstance(ctx, "acme", target, ReleaseComponentInstanceCreate{ComponentName: "api", Name: "api-blue"})
		}},
		{"RegisterComponentInstance", "POST", "component-instances", "", `{"component_name":"api","target_id":"` + releaseTestID + `"}`, releaseInstanceJSON, 200, func(s *ReleasesService) (any, *Response, error) {
			return s.RegisterComponentInstance(ctx, "acme", ReleaseComponentInstanceRegister{ReleaseComponentInstanceCreate: ReleaseComponentInstanceCreate{ComponentName: "api"}, TargetID: releaseTestID})
		}},
		{"CreateArtifact", "POST", "artifacts", "", `{"name":"service"}`, releaseArtifactJSON, 201, func(s *ReleasesService) (any, *Response, error) {
			return s.CreateArtifact(ctx, "acme", ReleaseArtifactCreate{Name: "service"})
		}},
		{"PublishVersion", "POST", "artifacts/service/versions", "", `{"digest":"` + releaseTestDigest + `","source_build":{"id":"build-id"},"pull_requests":[{"url":"https://example.com/pr/2"}]}`, releaseVersionJSON, 201, func(s *ReleasesService) (any, *Response, error) {
			return s.PublishVersion(ctx, "acme", artifact, ReleaseVersionPublish{Digest: releaseTestDigest, SourceBuild: Some(ReleaseReference{ID: "build-id"}), PullRequests: []ReleasePullRequest{{URL: "https://example.com/pr/2"}}})
		}},
		{"PublishByArtifactID", "POST", "artifacts/by-id/" + releaseTestID + "/versions", "", `{"digest":"` + releaseTestDigest + `"}`, releaseVersionJSON, 201, func(s *ReleasesService) (any, *Response, error) {
			return s.PublishVersion(ctx, "acme", ReleaseArtifactRef{ID: releaseTestID}, ReleaseVersionPublish{Digest: releaseTestDigest})
		}},
		{"YankVersion", "POST", "artifact-versions/by-id/" + releaseTestID + "/yank", "", `{"reason":"withdrawn"}`, releaseVersionJSON, 200, func(s *ReleasesService) (any, *Response, error) {
			return s.YankVersion(ctx, "acme", ReleaseVersionRef{ID: releaseTestID}, ReleaseVersionYank{Reason: "withdrawn"})
		}},
		{"ArchiveTarget", "POST", targetPath + "/archive", "", `{}`, releaseTargetJSON, 200, func(s *ReleasesService) (any, *Response, error) { return s.ArchiveTarget(ctx, "acme", target) }},
		{"UnarchiveTarget", "POST", targetPath + "/unarchive", "", `{}`, releaseTargetJSON, 200, func(s *ReleasesService) (any, *Response, error) { return s.UnarchiveTarget(ctx, "acme", target) }},
		{"ArchiveComponent", "POST", "components/api/archive", "", `{}`, releaseComponentJSON, 200, func(s *ReleasesService) (any, *Response, error) { return s.ArchiveComponent(ctx, "acme", component) }},
		{"UnarchiveComponent", "POST", "components/api/unarchive", "", `{}`, releaseComponentJSON, 200, func(s *ReleasesService) (any, *Response, error) { return s.UnarchiveComponent(ctx, "acme", component) }},
		{"ArchiveInstance", "POST", "component-instances/by-id/" + releaseTestID + "/archive", "", `{}`, releaseInstanceJSON, 200, func(s *ReleasesService) (any, *Response, error) {
			return s.ArchiveComponentInstance(ctx, "acme", ReleaseComponentInstanceRef{ID: releaseTestID})
		}},
		{"UnarchiveInstance", "POST", targetPath + "/instances/api-blue/unarchive", "", `{}`, releaseInstanceJSON, 200, func(s *ReleasesService) (any, *Response, error) {
			return s.UnarchiveComponentInstance(ctx, "acme", instance)
		}},
		{"ArchiveArtifact", "POST", "artifacts/service/archive", "", `{}`, releaseArtifactJSON, 200, func(s *ReleasesService) (any, *Response, error) { return s.ArchiveArtifact(ctx, "acme", artifact) }},
		{"UnarchiveArtifact", "POST", "artifacts/by-id/" + releaseTestID + "/unarchive", "", `{}`, releaseArtifactJSON, 200, func(s *ReleasesService) (any, *Response, error) {
			return s.UnarchiveArtifact(ctx, "acme", ReleaseArtifactRef{ID: releaseTestID})
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, tt.method)
				if r.URL.EscapedPath() != "/prefix/v2/organizations/acme/"+tt.path {
					t.Errorf("wire path = %q, want %q", r.URL.EscapedPath(), tt.path)
				}
				if r.URL.Query().Encode() != must(url.ParseQuery(tt.query)).Encode() {
					t.Errorf("query = %q, want %q", r.URL.RawQuery, tt.query)
				}
				if r.Header.Get("Authorization") != "Bearer release-test-token" {
					t.Error("missing authentication")
				}
				if tt.request != "" {
					if r.Header.Get("Content-Type") != "application/json" {
						t.Error("missing JSON content type")
					}
					var body json.RawMessage
					if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
						t.Error(err)
					}
					assertJSONEqual(t, string(body), tt.request)
				}
				w.WriteHeader(tt.status)
				_, _ = fmt.Fprint(w, tt.response)
			}))
			t.Cleanup(server.Close)
			client := must(NewClient(WithBaseURL(server.URL + "/prefix/")))
			client.authHeader = "Bearer release-test-token"
			got, resp, err := tt.call(client.Releases)
			if err != nil || resp == nil || resp.StatusCode != tt.status {
				t.Fatalf("response=%#v, error=%v", resp, err)
			}
			assertJSONEqual(t, string(must(json.Marshal(got))), tt.response)
			if status, ok := got.(ReleaseComponentInstanceStatus); ok && *status.Sequence != 9007199254740993 {
				t.Error("sequence lost precision")
			}
		})
	}
}

func TestReleasesService_PublicationFailureIsNotRetried(t *testing.T) {
	t.Parallel()
	for _, status := range []int{500, 201} {
		t.Run(fmt.Sprint(status), func(t *testing.T) {
			calls := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				calls++
				w.WriteHeader(status)
				if status == 201 {
					_, _ = fmt.Fprint(w, `{"id":`)
				} else {
					_, _ = fmt.Fprint(w, `{"message":"outcome unknown"}`)
				}
			}))
			t.Cleanup(server.Close)
			_, _, err := must(NewClient(WithBaseURL(server.URL))).Releases.PublishVersion(context.Background(), "acme", ReleaseArtifactRef{Name: "service"}, ReleaseVersionPublish{Digest: releaseTestDigest})
			if err == nil || calls != 1 {
				t.Fatalf("error=%v, requests=%d; an ambiguous publication must not be replayed", err, calls)
			}
		})
	}
}

func TestReleasesService_OptionalFields(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	tests := []struct {
		request string
		call    func(*ReleasesService) (any, *Response, error)
	}{
		{`{"name":"east"}`, func(s *ReleasesService) (any, *Response, error) {
			return s.CreateTarget(ctx, "acme", ReleaseTargetCreate{Name: "east"})
		}},
		{`{}`, func(s *ReleasesService) (any, *Response, error) {
			return s.UpdateTarget(ctx, "acme", ReleaseTargetRef{Name: "east"}, ReleaseTargetUpdate{})
		}},
		{`{"component_name":"api"}`, func(s *ReleasesService) (any, *Response, error) {
			return s.CreateComponentInstance(ctx, "acme", ReleaseTargetRef{Name: "east"}, ReleaseComponentInstanceCreate{ComponentName: "api"})
		}},
		{`{"digest":"` + releaseTestDigest + `","source_build":null,"source_url":null}`, func(s *ReleasesService) (any, *Response, error) {
			return s.PublishVersion(ctx, "acme", ReleaseArtifactRef{Name: "service"}, ReleaseVersionPublish{Digest: releaseTestDigest, SourceBuild: Null[ReleaseReference](), SourceURL: Null[string]()})
		}},
		{`{"digest":"` + releaseTestDigest + `","source_url":""}`, func(s *ReleasesService) (any, *Response, error) {
			return s.PublishVersion(ctx, "acme", ReleaseArtifactRef{Name: "service"}, ReleaseVersionPublish{Digest: releaseTestDigest, SourceURL: Some("")})
		}},
		{`{}`, func(s *ReleasesService) (any, *Response, error) {
			return s.YankVersion(ctx, "acme", ReleaseVersionRef{ID: releaseTestID}, ReleaseVersionYank{})
		}},
	}
	for _, tt := range tests {
		t.Run(tt.request, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var body json.RawMessage
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					t.Error(err)
				}
				assertJSONEqual(t, string(body), tt.request)
				_, _ = fmt.Fprint(w, `{}`)
			}))
			t.Cleanup(server.Close)
			_, _, err := tt.call(must(NewClient(WithBaseURL(server.URL))).Releases)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestReleasesService_ObservationBatch(t *testing.T) {
	t.Parallel()
	for _, status := range []int{201, 200, 409, 422, 500} {
		t.Run(fmt.Sprint(status), func(t *testing.T) {
			calls := 0
			response := `[{"created":false,` + releaseObservationJSON[1:] + `,{"id":"second","target_id":"target-id","sequence":9007199254740994,"event_id":"019a1dd8-a49b-70c5-99b8-b8f9f0ed6825","type":"ComponentInstanceDisappeared","instance":"api-blue","occurred_at":"2026-09-20T01:02:03Z","recorded_at":"2026-09-20T01:03:04Z","created":true}]`
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls++
				if r.URL.Path != "/v2/organizations/acme/targets/east/observations" {
					t.Error(r.URL.Path)
				}
				var body json.RawMessage
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					t.Error(err)
				}
				assertJSONEqual(t, string(body), `{"events":[{"event_id":"`+releaseTestID+`","type":"ComponentInstanceVersionObserved","instance":"api-blue","digest":"`+releaseTestDigest+`"},{"event_id":"019a1dd8-a49b-70c5-99b8-b8f9f0ed6825","type":"ComponentInstanceDisappeared","instance":"api-blue"}]}`)
				w.WriteHeader(status)
				if status < 300 {
					_, _ = fmt.Fprint(w, response)
				} else {
					_, _ = fmt.Fprint(w, `{"message":"observation rejected","errors":[{"field":"/events/0/event_id","code":"conflict"}]}`)
				}
			}))
			t.Cleanup(server.Close)
			got, resp, err := must(NewClient(WithBaseURL(server.URL))).Releases.AppendObservations(context.Background(), "acme", ReleaseTargetRef{Name: "east"}, []ReleaseObservationCreate{{EventID: releaseTestID, Type: "ComponentInstanceVersionObserved", Instance: "api-blue", Digest: releaseTestDigest}, {EventID: "019a1dd8-a49b-70c5-99b8-b8f9f0ed6825", Type: "ComponentInstanceDisappeared", Instance: "api-blue"}})
			if resp == nil || resp.StatusCode != status || calls != 1 {
				t.Fatalf("response=%#v, calls=%d, error=%v", resp, calls, err)
			}
			if status < 300 {
				if err != nil {
					t.Fatal(err)
				}
				assertJSONEqual(t, string(must(json.Marshal(got))), response)
				if got[0].Sequence != 9007199254740993 || got[0].Created || !got[1].Created {
					t.Error("order, precision or replay state lost")
				}
			} else {
				var apiErr *ErrorResponse
				if !errors.As(err, &apiErr) || apiErr.Message != "observation rejected" || !strings.Contains(string(apiErr.RawBody), "/events/0/event_id") {
					t.Fatalf("structured API error lost: %v", err)
				}
			}
		})
	}
}

func TestReleasesService_InvalidSelectors(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	calls := []func(*ReleasesService) (any, *Response, error){
		func(s *ReleasesService) (any, *Response, error) { return s.GetTarget(ctx, "acme", ReleaseTargetRef{}) },
		func(s *ReleasesService) (any, *Response, error) {
			return s.GetTarget(ctx, "acme", ReleaseTargetRef{ID: releaseTestID, Name: "east"})
		},
		func(s *ReleasesService) (any, *Response, error) {
			return s.GetComponent(ctx, "acme", ReleaseComponentRef{})
		},
		func(s *ReleasesService) (any, *Response, error) {
			return s.GetComponent(ctx, "acme", ReleaseComponentRef{ID: releaseTestID, Name: "api"})
		},
		func(s *ReleasesService) (any, *Response, error) {
			return s.GetComponentInstance(ctx, "acme", ReleaseComponentInstanceRef{Name: "api"})
		},
		func(s *ReleasesService) (any, *Response, error) {
			return s.GetComponentInstance(ctx, "acme", ReleaseComponentInstanceRef{ID: releaseTestID, Target: ReleaseTargetRef{Name: "east"}})
		},
		func(s *ReleasesService) (any, *Response, error) {
			return s.GetComponentInstance(ctx, "acme", ReleaseComponentInstanceRef{ID: releaseTestID, Name: "api"})
		},
		func(s *ReleasesService) (any, *Response, error) {
			return s.GetComponentInstanceStatus(ctx, "acme", ReleaseTargetRef{Name: "east"}, "")
		},
		func(s *ReleasesService) (any, *Response, error) {
			return s.ListVersions(ctx, "acme", ReleaseArtifactRef{ID: releaseTestID, Name: "api"}, nil)
		},
		func(s *ReleasesService) (any, *Response, error) {
			return s.GetVersion(ctx, "acme", ReleaseVersionRef{})
		},
		func(s *ReleasesService) (any, *Response, error) {
			return s.CreateComponent(ctx, "acme", ReleaseComponentCreate{Name: "api"})
		},
		func(s *ReleasesService) (any, *Response, error) {
			return s.CreateComponent(ctx, "acme", ReleaseComponentCreate{Name: "api", Artifact: ReleaseArtifactRef{ID: releaseTestID, Name: "api"}})
		},
		func(s *ReleasesService) (any, *Response, error) {
			return s.CreateComponentInstance(ctx, "acme", ReleaseTargetRef{}, ReleaseComponentInstanceCreate{})
		},
	}
	for i, call := range calls {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			server, client, teardown := newMockServerAndClient(t)
			t.Cleanup(teardown)
			_, resp, err := call(client.Releases)
			if err == nil || resp != nil || len(server.calls) != 0 {
				t.Fatalf("response=%#v, error=%v, HTTP calls=%v", resp, err, server.calls)
			}
		})
	}
}

func TestReleasesService_ComponentInstanceErrors(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	for _, tt := range []struct {
		name string
		call func(*ReleasesService) error
	}{
		{name: "missing selector", call: func(s *ReleasesService) error {
			_, _, err := s.GetComponentInstance(ctx, "acme", ReleaseComponentInstanceRef{})
			return err
		}},
		{name: "conflicting selectors", call: func(s *ReleasesService) error {
			_, _, err := s.GetComponentInstance(ctx, "acme", ReleaseComponentInstanceRef{ID: releaseTestID, Name: "api"})
			return err
		}},
		{name: "missing status name", call: func(s *ReleasesService) error {
			_, _, err := s.GetComponentInstanceStatus(ctx, "acme", ReleaseTargetRef{Name: "east"}, "")
			return err
		}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			_, client, teardown := newMockServerAndClient(t)
			t.Cleanup(teardown)
			if err := tt.call(client.Releases); err == nil || !strings.Contains(err.Error(), "component instance") {
				t.Fatalf("error = %v, want component instance context", err)
			}
		})
	}
}

func TestReleasesService_ResponseVariants(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	target := ReleaseTargetRef{Name: "east"}
	tests := []struct {
		name, body string
		call       func(*ReleasesService) (any, *Response, error)
	}{
		{"never observed", `{"target_id":"target-id","instance":"api-blue","instance_id":"instance-id","present":null,"active_digest":null,"last_active_digest":null,"first_seen_at":null,"last_seen_at":null,"sequence":null,"catalog_resolution":{"status":"unknown","reason":"no_active_version"}}`, func(s *ReleasesService) (any, *Response, error) {
			return s.GetComponentInstanceStatus(ctx, "acme", target, "api-blue")
		}},
		{"resolved present", `{"target_id":"target-id","instance":"api-blue","instance_id":"instance-id","present":true,"active_digest":"` + releaseTestDigest + `","last_active_digest":"` + releaseTestDigest + `","first_seen_at":null,"last_seen_at":null,"sequence":1,"catalog_resolution":{"status":"resolved","reason":null}}`, func(s *ReleasesService) (any, *Response, error) {
			return s.GetComponentInstanceStatus(ctx, "acme", target, "api-blue")
		}},
		{"yanked provenance", `{"id":"digest-id","digest":"` + releaseTestDigest + `","source_url":"opaque source","source_build":{"id":"build-id"},"pull_requests":[],"created_by":{"id":"user-id","graphql_id":"opaque-user","name":null,"email":"person@example.com","avatar_url":null,"created_at":"2026-09-20T01:02:03Z"},"url":"https://api.example/digest","id_url":"https://api.example/digest-id","created_at":"2026-09-20T01:02:03Z","yanked":true,"yanked_at":"2026-09-21T01:02:03Z","yank_reason":"withdrawn"}`, func(s *ReleasesService) (any, *Response, error) {
			return s.GetVersion(ctx, "acme", ReleaseVersionRef{ID: releaseTestID})
		}},
		{"legacy digestless publication", `{"id":"legacy-id","digest":null,"source_url":null,"source_build":null,"pull_requests":[],"created_by":null,"url":"https://api.example/legacy","id_url":"https://api.example/legacy","created_at":"2026-09-20T01:02:03Z","yanked":false,"yanked_at":null,"yank_reason":null}`, func(s *ReleasesService) (any, *Response, error) {
			return s.GetVersion(ctx, "acme", ReleaseVersionRef{ID: releaseTestID})
		}},
		{"history payload presence", releaseTestPage(`{"id":"one","target_id":"target-id","sequence":1,"event_id":"` + releaseTestID + `","type":"ComponentInstanceAppeared","instance":"api-blue","occurred_at":"2026-09-20T01:02:03Z","recorded_at":"2026-09-20T01:03:04Z"},{"id":"two","target_id":"target-id","sequence":2,"event_id":"019a1dd8-a49b-70c5-99b8-b8f9f0ed6825","type":"ComponentInstanceVersionObserved","instance":"api-blue","occurred_at":"2026-09-20T01:02:03Z","recorded_at":"2026-09-20T01:03:04Z","catalog_resolution":{"status":"unknown","reason":"no_active_version"}}`), func(s *ReleasesService) (any, *Response, error) { return s.ListObservations(ctx, "acme", target, nil) }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { _, _ = fmt.Fprint(w, tt.body) }))
			t.Cleanup(server.Close)
			got, _, err := tt.call(must(NewClient(WithBaseURL(server.URL))).Releases)
			if err != nil {
				t.Fatal(err)
			}
			assertJSONEqual(t, string(must(json.Marshal(got))), tt.body)
		})
	}
}

func TestReleasesService_APIErrorsAndContext(t *testing.T) {
	t.Parallel()
	for _, status := range []int{401, 403, 404, 422} {
		t.Run(fmt.Sprint(status), func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(status)
				_, _ = fmt.Fprint(w, `{"message":"not available"}`)
			}))
			t.Cleanup(server.Close)
			_, resp, err := must(NewClient(WithBaseURL(server.URL))).Releases.ListComponents(context.Background(), "acme", nil)
			var apiErr *ErrorResponse
			if !errors.As(err, &apiErr) || apiErr.Message != "not available" || resp == nil || resp.StatusCode != status {
				t.Fatalf("response=%#v, error=%v", resp, err)
			}
		})
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)
	_, _, err := client.Releases.ListTargets(ctx, "acme", nil)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("error=%v, want context.Canceled", err)
	}
}

func TestReleasesService_Escaping(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.EscapedPath() != "/prefix/v2/organizations/acme%2Fspace%2Bname/targets/prod%2Fblue%3F%23%25%20%2B" {
			t.Error(r.URL.EscapedPath())
		}
		_, _ = fmt.Fprint(w, `{}`)
	}))
	t.Cleanup(server.Close)
	_, _, err := must(NewClient(WithBaseURL(server.URL+"/prefix/"))).Releases.GetTarget(context.Background(), "acme/space+name", ReleaseTargetRef{Name: "prod/blue?#% +"})
	if err != nil {
		t.Fatal(err)
	}
}

func TestReleasePageLink_ToOptions(t *testing.T) {
	got, err := ReleasePageLink("https://api.example/component-instances?before=prior%2Bcursor&per_page=7&label[region]=au").ToOptions()
	if err != nil || *got != (ReleasesListOptions{Before: "prior+cursor", PerPage: 7}) {
		t.Fatalf("options=%#v, error=%v", got, err)
	}
	for _, link := range []ReleasePageLink{"https://api.example/%zz", "https://api.example/targets?per_page=many"} {
		if _, err := link.ToOptions(); err == nil {
			t.Fatalf("expected error for %q", link)
		}
	}
}
