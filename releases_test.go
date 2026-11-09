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
	"sync/atomic"
	"testing"
	"time"
)

const releaseTestID = "019a1dd8-a49b-70c5-99b8-b8f9f0ed6824"

// Fixtures follow the Rails presenters, including required null fields and
// optional fields that must remain absent when re-encoded for an MCP consumer.
const releaseInstallationJSON = `{"id":"installation-id","name":"prod","url":"https://api.example/installations/prod","id_url":"https://api.example/installations/by-id/installation-id","created_at":"2026-09-20T01:02:03Z","archived_at":null}`
const releaseTargetJSON = `{"id":"target-id","name":"east","labels":{"region":"us-east-1","empty":"","space":" "},"url":"https://api.example/target","id_url":"https://api.example/target-id","created_at":"2026-09-20T01:02:03Z","archived_at":"2026-09-21T01:02:03Z"}`
const releaseComponentJSON = `{"id":"component-id","name":"api","target_id":"target-id","artifact":{"id":"artifact-id","name":"api","version_scheme":"semver"},"url":"https://api.example/component","id_url":"https://api.example/component-id","created_at":"2026-09-20T01:02:03Z","archived_at":null}`
const releaseArtifactJSON = `{"id":"artifact-id","name":"api","version_scheme":"build_number","pipeline":{"id":"pipeline-id"},"url":"https://api.example/artifact","id_url":"https://api.example/artifact-id","created_at":"2026-09-20T01:02:03Z","archived_at":null}`
const releaseVersionJSON = `{"id":"version-id","version":"1.2.3+build.4","digest":null,"source_url":null,"source_build":null,"pull_requests":[{"url":"https://github.com/acme/api/pull/12","provider":"github","repository_url":"https://github.com/acme/api","id":"12"},{"url":"https://example.com/pr/2"}],"created_by":null,"url":"https://api.example/version","id_url":"https://api.example/version-id","created_at":"2026-09-20T01:02:03Z","yanked":false,"yanked_at":null,"yank_reason":null}`
const releaseStatusJSON = `{"target_id":"target-id","component":"api","component_id":"component-id","present":false,"active_version":null,"active_digest":null,"last_active_version":"1.2.3","first_seen_at":"2026-09-20T01:02:03Z","last_seen_at":"2026-09-21T01:02:03Z","sequence":0,"catalog_resolution":{"status":"unknown","artifact_version_id":null,"digest_check":"not_checked","reason":"component_absent"}}`
const releaseComparisonJSON = `{"artifact":{"id":"artifact-id","name":"api","version_scheme":"semver"},"left":{"id":"left-id","version":"2.0.0","digest":null,"yanked":false},"right":{"id":"right-id","version":"1.0.0","digest":"sha256:abc","yanked":true},"basis":"semver","result":"newer","reason":"higher_semver_precedence","same_version":false,"same_content":null}`

func TestReleasesService_Writes(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	installation := ReleaseInstallationRef{ID: releaseTestID}
	target := ReleaseTargetRef{Installation: installation, Name: "east"}
	component := ReleaseComponentRef{ID: releaseTestID}
	artifact := ReleaseArtifactRef{Name: "api"}
	version := ReleaseVersionRef{Artifact: artifact, Version: "1.2.3+build.4"}
	const digest = "sha256:0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	const observationJSON = `{"id":"observation-id","target_id":"target-id","sequence":9007199254740993,"event_id":"019a1dd8-a49b-70c5-99b8-b8f9f0ed6824","type":"ComponentVersionObserved","component":"api","version":"1.2.3","occurred_at":"2026-09-20T01:02:03Z","recorded_at":"2026-09-20T01:03:04Z","created":false}`
	tests := []struct {
		name, method, path, request, response string
		status                                int
		call                                  func(*ReleasesService) (any, *Response, error)
	}{
		{
			name: "CreateInstallation", method: "POST", path: "installations", request: `{"name":"prod"}`, response: releaseInstallationJSON, status: 201,
			call: func(s *ReleasesService) (any, *Response, error) {
				return s.CreateInstallation(ctx, "acme", ReleaseInstallationCreate{Name: "prod"})
			},
		},
		{
			name: "CreateTarget", method: "POST", path: "installations/by-id/" + releaseTestID + "/targets", request: `{"name":"east","labels":{"region":"au","empty":""}}`, response: releaseTargetJSON, status: 201,
			call: func(s *ReleasesService) (any, *Response, error) {
				return s.CreateTarget(ctx, "acme", installation, ReleaseTargetCreate{Name: "east", Labels: map[string]string{"region": "au", "empty": ""}})
			},
		},
		{
			name: "UpdateTarget", method: "PATCH", path: "installations/by-id/" + releaseTestID + "/targets/east", request: `{"labels":{}}`, response: releaseTargetJSON, status: 200,
			call: func(s *ReleasesService) (any, *Response, error) {
				return s.UpdateTarget(ctx, "acme", target, ReleaseTargetUpdate{Labels: Some(map[string]string{})})
			},
		},
		{
			name: "CreateComponent", method: "POST", path: "targets/by-id/" + releaseTestID + "/components", request: `{"name":"api","artifact":{"id":"artifact-id"}}`, response: releaseComponentJSON, status: 201,
			call: func(s *ReleasesService) (any, *Response, error) {
				return s.CreateComponent(ctx, "acme", ReleaseTargetRef{ID: releaseTestID}, ReleaseComponentCreate{Name: "api", Artifact: &ReleaseArtifactRef{ID: "artifact-id"}})
			},
		},
		{
			name: "CreateArtifact", method: "POST", path: "artifacts", request: `{"name":"api","version_scheme":"build_number","pipeline":{"id":"pipeline-id"}}`, response: releaseArtifactJSON, status: 201,
			call: func(s *ReleasesService) (any, *Response, error) {
				return s.CreateArtifact(ctx, "acme", ReleaseArtifactCreate{Name: "api", VersionScheme: "build_number", Pipeline: &ReleaseReference{ID: "pipeline-id"}})
			},
		},
		{
			name: "PublishVersion", method: "POST", path: "artifacts/api/versions", request: `{"version":"1.2.3","digest":"` + digest + `","source_build":{"id":"build-id"},"source_url":"https://example.com/source","pull_requests":[{"url":"https://github.com/acme/api/pull/12","provider":"github","repository_url":"https://github.com/acme/api","id":"12"}]}`, response: releaseVersionJSON, status: 201,
			call: func(s *ReleasesService) (any, *Response, error) {
				return s.PublishVersion(ctx, "acme", artifact, ReleaseVersionPublish{Version: "1.2.3", Digest: digest, SourceBuild: Some(ReleaseReference{ID: "build-id"}), SourceURL: Some("https://example.com/source"), PullRequests: []ReleasePullRequest{{URL: "https://github.com/acme/api/pull/12", Provider: "github", RepositoryURL: "https://github.com/acme/api", ID: "12"}}})
			},
		},
		{
			name: "YankVersion", method: "POST", path: "artifacts/api/versions/1.2.3%2Bbuild.4/yank", request: `{"reason":"incorrect content"}`, response: releaseVersionJSON, status: 200,
			call: func(s *ReleasesService) (any, *Response, error) {
				return s.YankVersion(ctx, "acme", version, ReleaseVersionYank{Reason: "incorrect content"})
			},
		},
		{
			name: "AppendObservations", method: "POST", path: "installations/by-id/" + releaseTestID + "/targets/east/observations", request: `{"events":[{"event_id":"019a1dd8-a49b-70c5-99b8-b8f9f0ed6824","type":"ComponentVersionObserved","component":"api","version":"1.2.3","digest":"` + digest + `","occurred_at":"2026-09-20T01:02:03Z"},{"event_id":"019a1dd8-a49b-70c5-99b8-b8f9f0ed6825","type":"TargetAppeared"}]}`, response: `[` + observationJSON + `,{"id":"second","target_id":"target-id","sequence":9007199254740994,"event_id":"019a1dd8-a49b-70c5-99b8-b8f9f0ed6825","type":"TargetAppeared","occurred_at":"2026-09-20T01:02:03Z","recorded_at":"2026-09-20T01:03:04Z","created":true}]`, status: 200,
			call: func(s *ReleasesService) (any, *Response, error) {
				return s.AppendObservations(ctx, "acme", target, []ReleaseObservationCreate{
					{EventID: releaseTestID, Type: "ComponentVersionObserved", Component: "api", Version: "1.2.3", Digest: digest, OccurredAt: NewTimestamp(must(time.Parse(time.RFC3339, "2026-09-20T01:02:03Z")))},
					{EventID: "019a1dd8-a49b-70c5-99b8-b8f9f0ed6825", Type: "TargetAppeared"},
				})
			},
		},
		{
			name: "ArchiveInstallation", method: "POST", path: "installations/by-id/" + releaseTestID + "/archive", request: `{}`, response: releaseInstallationJSON, status: 200,
			call: func(s *ReleasesService) (any, *Response, error) {
				return s.ArchiveInstallation(ctx, "acme", installation)
			},
		},
		{
			name: "UnarchiveInstallation", method: "POST", path: "installations/prod/unarchive", request: `{}`, response: releaseInstallationJSON, status: 200,
			call: func(s *ReleasesService) (any, *Response, error) {
				return s.UnarchiveInstallation(ctx, "acme", ReleaseInstallationRef{Name: "prod"})
			},
		},
		{
			name: "ArchiveTarget", method: "POST", path: "installations/by-id/" + releaseTestID + "/targets/east/archive", request: `{}`, response: releaseTargetJSON, status: 200,
			call: func(s *ReleasesService) (any, *Response, error) { return s.ArchiveTarget(ctx, "acme", target) },
		},
		{
			name: "UnarchiveTarget", method: "POST", path: "targets/by-id/" + releaseTestID + "/unarchive", request: `{}`, response: releaseTargetJSON, status: 200,
			call: func(s *ReleasesService) (any, *Response, error) {
				return s.UnarchiveTarget(ctx, "acme", ReleaseTargetRef{ID: releaseTestID})
			},
		},
		{
			name: "ArchiveComponent", method: "POST", path: "components/by-id/" + releaseTestID + "/archive", request: `{}`, response: releaseComponentJSON, status: 200,
			call: func(s *ReleasesService) (any, *Response, error) { return s.ArchiveComponent(ctx, "acme", component) },
		},
		{
			name: "UnarchiveComponent", method: "POST", path: "installations/by-id/" + releaseTestID + "/targets/east/components/api/unarchive", request: `{}`, response: releaseComponentJSON, status: 200,
			call: func(s *ReleasesService) (any, *Response, error) {
				return s.UnarchiveComponent(ctx, "acme", ReleaseComponentRef{Target: target, Name: "api"})
			},
		},
		{
			name: "ArchiveArtifact", method: "POST", path: "artifacts/api/archive", request: `{}`, response: releaseArtifactJSON, status: 200,
			call: func(s *ReleasesService) (any, *Response, error) { return s.ArchiveArtifact(ctx, "acme", artifact) },
		},
		{
			name: "UnarchiveArtifact", method: "POST", path: "artifacts/by-id/" + releaseTestID + "/unarchive", request: `{}`, response: releaseArtifactJSON, status: 200,
			call: func(s *ReleasesService) (any, *Response, error) {
				return s.UnarchiveArtifact(ctx, "acme", ReleaseArtifactRef{ID: releaseTestID})
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, tt.method)
				if r.URL.EscapedPath() != "/prefix/v2/organizations/acme/"+tt.path {
					t.Errorf("wire path = %q, want %q", r.URL.EscapedPath(), tt.path)
				}
				if r.Header.Get("Authorization") != "Bearer release-test-token" || r.Header.Get("Content-Type") != "application/json" {
					t.Error("missing authentication or JSON content type")
				}
				var body json.RawMessage
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					t.Error(err)
				}
				assertJSONEqual(t, string(body), tt.request)
				w.WriteHeader(tt.status)
				_, _ = fmt.Fprint(w, tt.response)
			}))
			t.Cleanup(server.Close)
			client := must(NewClient(WithBaseURL(server.URL + "/prefix/")))
			client.authHeader = "Bearer release-test-token"
			got, resp, err := tt.call(client.Releases)
			if err != nil || resp == nil || resp.StatusCode != tt.status {
				t.Fatalf("response = %#v, err = %v", resp, err)
			}
			assertJSONEqual(t, string(must(json.Marshal(got))), tt.response)
			if results, ok := got.([]ReleaseObservationWriteResult); ok && (results[0].Sequence != 9007199254740993 || results[0].Created || !results[1].Created) {
				t.Errorf("observation order, precision or replay state lost: %#v", results)
			}
		})
	}
}

func TestReleasesService_WriteOptionalFields(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	target := ReleaseTargetRef{ID: releaseTestID}
	artifact := ReleaseArtifactRef{ID: releaseTestID}
	tests := []struct {
		name, request string
		call          func(*ReleasesService) (any, *Response, error)
	}{
		{"default target labels", `{"name":"east"}`, func(s *ReleasesService) (any, *Response, error) {
			return s.CreateTarget(ctx, "acme", ReleaseInstallationRef{Name: "prod"}, ReleaseTargetCreate{Name: "east"})
		}},
		{"unchanged labels", `{}`, func(s *ReleasesService) (any, *Response, error) {
			return s.UpdateTarget(ctx, "acme", target, ReleaseTargetUpdate{})
		}},
		{"replacement labels", `{"labels":{"region":"au"}}`, func(s *ReleasesService) (any, *Response, error) {
			return s.UpdateTarget(ctx, "acme", target, ReleaseTargetUpdate{Labels: Some(map[string]string{"region": "au"})})
		}},
		{"unbound component", `{"name":"api"}`, func(s *ReleasesService) (any, *Response, error) {
			return s.CreateComponent(ctx, "acme", target, ReleaseComponentCreate{Name: "api"})
		}},
		{"artifact name binding", `{"name":"api","artifact":{"name":"api"}}`, func(s *ReleasesService) (any, *Response, error) {
			return s.CreateComponent(ctx, "acme", target, ReleaseComponentCreate{Name: "api", Artifact: &ReleaseArtifactRef{Name: "api"}})
		}},
		{"default version scheme", `{"name":"api"}`, func(s *ReleasesService) (any, *Response, error) {
			return s.CreateArtifact(ctx, "acme", ReleaseArtifactCreate{Name: "api"})
		}},
		{"OIDC source build default", `{"version":"1.2.3","digest":"sha256:abc"}`, func(s *ReleasesService) (any, *Response, error) {
			return s.PublishVersion(ctx, "acme", artifact, ReleaseVersionPublish{Version: "1.2.3", Digest: "sha256:abc"})
		}},
		{"explicit null provenance", `{"version":"1.2.3","digest":"sha256:abc","source_build":null,"source_url":null}`, func(s *ReleasesService) (any, *Response, error) {
			return s.PublishVersion(ctx, "acme", artifact, ReleaseVersionPublish{Version: "1.2.3", Digest: "sha256:abc", SourceBuild: Null[ReleaseReference](), SourceURL: Null[string]()})
		}},
		{"empty source URL", `{"version":"1.2.3","digest":"sha256:abc","source_url":""}`, func(s *ReleasesService) (any, *Response, error) {
			return s.PublishVersion(ctx, "acme", artifact, ReleaseVersionPublish{Version: "1.2.3", Digest: "sha256:abc", SourceURL: Some("")})
		}},
		{"no yank reason", `{}`, func(s *ReleasesService) (any, *Response, error) {
			return s.YankVersion(ctx, "acme", ReleaseVersionRef{ID: releaseTestID}, ReleaseVersionYank{})
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var body json.RawMessage
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					t.Error(err)
				}
				assertJSONEqual(t, string(body), tt.request)
				_, _ = fmt.Fprint(w, `{}`)
			}))
			t.Cleanup(server.Close)
			client := must(NewClient(WithBaseURL(server.URL)))
			if _, _, err := tt.call(client.Releases); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestReleasesService_ObservationWriteResponses(t *testing.T) {
	t.Parallel()
	for _, status := range []int{201, 200, 409, 422, 500} {
		t.Run(fmt.Sprint(status), func(t *testing.T) {
			t.Parallel()
			var calls atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls.Add(1)
				w.WriteHeader(status)
				if status < 300 {
					_, _ = fmt.Fprintf(w, `[{"id":"observation-id","event_id":%q,"type":"TargetAppeared","target_id":"target-id","sequence":1,"occurred_at":"2026-09-20T01:02:03Z","recorded_at":"2026-09-20T01:03:04Z","created":%t}]`, releaseTestID, status == 201)
				} else {
					_, _ = fmt.Fprint(w, `{"message":"observation rejected","errors":[{"field":"/events/0/event_id","code":"conflict"}]}`)
				}
			}))
			t.Cleanup(server.Close)
			client := must(NewClient(WithBaseURL(server.URL)))
			got, resp, err := client.Releases.AppendObservations(context.Background(), "acme", ReleaseTargetRef{ID: releaseTestID}, []ReleaseObservationCreate{{EventID: releaseTestID, Type: "TargetAppeared"}})
			if resp == nil || resp.StatusCode != status || calls.Load() != 1 {
				t.Fatalf("response=%#v, calls=%d, error=%v", resp, calls.Load(), err)
			}
			if status < 300 {
				if err != nil || len(got) != 1 || got[0].EventID != releaseTestID || got[0].Created != (status == 201) {
					t.Fatalf("results=%#v, error=%v", got, err)
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

func releaseTestPage(item string) string {
	return `{"items":[` + item + `],"links":{"self":"https://api.example/self","first":"https://api.example/first","prev":"https://api.example/prev","next":"https://api.example/next?after=opaque%2Bcursor"}}`
}

func TestReleasePageLink_ToOptions(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		name    string
		link    ReleasePageLink
		want    ReleasesListOptions
		wantErr bool
	}{
		{"next", "https://api.example/installations?after=opaque%2Bcursor%2F%3D&per_page=25", ReleasesListOptions{After: "opaque+cursor/=", PerPage: 25}, false},
		{"previous", "https://api.example/installations?before=prior%2Bcursor&per_page=7", ReleasesListOptions{Before: "prior+cursor", PerPage: 7}, false},
		{"first", "https://api.example/installations?per_page=10", ReleasesListOptions{PerPage: 10}, false},
		{"invalid link", "https://api.example/%zz", ReleasesListOptions{}, true},
		{"invalid page size", "https://api.example/installations?per_page=many", ReleasesListOptions{}, true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := tt.link.ToOptions()
			if tt.wantErr {
				if err == nil {
					t.Fatal("ToOptions returned nil error")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if *got != tt.want {
				t.Errorf("options = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestReleasesService_Reads(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	installation := ReleaseInstallationRef{Name: "prod"}
	target := ReleaseTargetRef{ID: releaseTestID}
	artifact := ReleaseArtifactRef{Name: "api"}
	pagination := ReleasesListOptions{After: "opaque+cursor/=", PerPage: 25}
	pageQuery := "after=opaque%2Bcursor%2F%3D&per_page=25"
	emptySuffix, notYanked := "", false
	tests := []struct {
		name, path, query, body string
		call                    func(*ReleasesService) (any, *Response, error)
	}{
		{
			name: "ListInstallations", path: "installations", query: pageQuery, body: releaseTestPage(releaseInstallationJSON),
			call: func(s *ReleasesService) (any, *Response, error) { return s.ListInstallations(ctx, "acme", &pagination) },
		},
		{
			name: "GetInstallation", path: "installations/prod", body: releaseInstallationJSON,
			call: func(s *ReleasesService) (any, *Response, error) { return s.GetInstallation(ctx, "acme", installation) },
		},
		{
			name: "ListTargets", path: "installations/by-id/" + releaseTestID + "/targets", query: pageQuery, body: releaseTestPage(releaseTargetJSON),
			call: func(s *ReleasesService) (any, *Response, error) {
				return s.ListTargets(ctx, "acme", ReleaseInstallationRef{ID: releaseTestID}, &pagination)
			},
		},
		{
			name: "GetTarget", path: "installations/prod/targets/east", body: releaseTargetJSON,
			call: func(s *ReleasesService) (any, *Response, error) {
				return s.GetTarget(ctx, "acme", ReleaseTargetRef{Installation: installation, Name: "east"})
			},
		},
		{
			name: "ListComponents", path: "targets/by-id/" + releaseTestID + "/components", query: pageQuery, body: releaseTestPage(releaseComponentJSON),
			call: func(s *ReleasesService) (any, *Response, error) {
				return s.ListComponents(ctx, "acme", target, &pagination)
			},
		},
		{
			name: "GetComponent", path: "components/by-id/" + releaseTestID, body: releaseComponentJSON,
			call: func(s *ReleasesService) (any, *Response, error) {
				return s.GetComponent(ctx, "acme", ReleaseComponentRef{ID: releaseTestID})
			},
		},
		{
			name: "ListComponentStatuses", path: "targets/by-id/" + releaseTestID + "/component-statuses", query: pageQuery, body: releaseTestPage(releaseStatusJSON),
			call: func(s *ReleasesService) (any, *Response, error) {
				return s.ListComponentStatuses(ctx, "acme", target, &pagination)
			},
		},
		{
			name: "GetComponentStatus", path: "targets/by-id/" + releaseTestID + "/component-statuses/api", body: releaseStatusJSON,
			call: func(s *ReleasesService) (any, *Response, error) {
				return s.GetComponentStatus(ctx, "acme", target, "api")
			},
		},
		{
			name: "ListObservations", path: "targets/by-id/" + releaseTestID + "/observations", query: "before=prior&component=api&per_page=10&type=ComponentAppeared%2CComponentVersionObserved",
			body: releaseTestPage(`{"id":"observation-id","target_id":"target-id","sequence":0,"event_id":"event-id","type":"TargetAppeared","occurred_at":"2026-09-20T01:02:03Z","recorded_at":"2026-09-20T01:03:04Z"}`),
			call: func(s *ReleasesService) (any, *Response, error) {
				return s.ListObservations(ctx, "acme", target, &ReleasesObservationListOptions{ReleasesListOptions: ReleasesListOptions{Before: "prior", PerPage: 10}, Component: "api", Type: "ComponentAppeared,ComponentVersionObserved"})
			},
		},
		{
			name: "ListArtifacts", path: "artifacts", query: pageQuery, body: releaseTestPage(releaseArtifactJSON),
			call: func(s *ReleasesService) (any, *Response, error) { return s.ListArtifacts(ctx, "acme", &pagination) },
		},
		{
			name: "GetArtifact", path: "artifacts/api", body: releaseArtifactJSON,
			call: func(s *ReleasesService) (any, *Response, error) { return s.GetArtifact(ctx, "acme", artifact) },
		},
		{
			name: "ListVersions", path: "artifacts/by-id/" + releaseTestID + "/versions", query: pageQuery + "&direction=asc&newer_than=baseline&order=precedence&suffix=&yanked=false",
			body: `{"items":[` + releaseVersionJSON + `],"links":{"self":"https://api.example/versions"},"ordering":{"basis":"build_number","direction":"asc"}}`,
			call: func(s *ReleasesService) (any, *Response, error) {
				return s.ListVersions(ctx, "acme", ReleaseArtifactRef{ID: releaseTestID}, &ReleasesVersionListOptions{ReleasesListOptions: pagination, Order: "precedence", Direction: "asc", NewerThan: "baseline", Suffix: &emptySuffix, Yanked: &notYanked})
			},
		},
		{
			name: "GetVersion", path: "artifact-versions/by-id/" + releaseTestID, body: releaseVersionJSON,
			call: func(s *ReleasesService) (any, *Response, error) {
				return s.GetVersion(ctx, "acme", ReleaseVersionRef{ID: releaseTestID})
			},
		},
		{
			name: "CompareVersions", path: "artifacts/api/version-comparison", query: "left=2.0.0%2Bbuild.4&right_id=baseline", body: releaseComparisonJSON,
			call: func(s *ReleasesService) (any, *Response, error) {
				return s.CompareVersions(ctx, "acme", artifact, &ReleasesComparisonOptions{Left: "2.0.0+build.4", RightID: "baseline"})
			},
		},
		{
			name: "ListComparisonPullRequests", path: "artifacts/by-id/" + releaseTestID + "/version-comparison/pull-requests", query: pageQuery + "&left_id=candidate&right=1.0.0%2Bbase",
			body: releaseTestPage(`{"url":"https://github.com/acme/api/pull/12","provider":"github","repository_url":"https://github.com/acme/api","id":"12"}`),
			call: func(s *ReleasesService) (any, *Response, error) {
				return s.ListComparisonPullRequests(ctx, "acme", ReleaseArtifactRef{ID: releaseTestID}, &ReleasesComparisonListOptions{ReleasesListOptions: pagination, ReleasesComparisonOptions: ReleasesComparisonOptions{LeftID: "candidate", Right: "1.0.0+base"}})
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server, client, teardown := newMockServerAndClient(t)
			t.Cleanup(teardown)
			if client.Releases == nil {
				t.Fatal("Client.Releases is nil")
			}
			client.authHeader = "Bearer release-test-token"
			server.HandleFunc("/v2/organizations/acme/"+tt.path, func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, http.MethodGet)
				if r.Header.Get("Authorization") != "Bearer release-test-token" {
					t.Error("existing authentication was not retained")
				}
				wantQuery := must(url.ParseQuery(tt.query))
				if r.URL.Query().Encode() != wantQuery.Encode() {
					t.Errorf("query = %q, want %q", r.URL.RawQuery, wantQuery.Encode())
				}
				_, _ = fmt.Fprint(w, tt.body)
			})
			got, resp, err := tt.call(client.Releases)
			if err != nil || resp == nil || resp.StatusCode != http.StatusOK {
				t.Fatalf("response = %#v, err = %v", resp, err)
			}
			encoded, err := json.Marshal(got)
			if err != nil {
				t.Fatal(err)
			}
			assertJSONEqual(t, string(encoded), tt.body)
		})
	}
}

func TestReleasesService_Addressing(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	installation := ReleaseInstallationRef{ID: releaseTestID}
	target := ReleaseTargetRef{Installation: installation, Name: "east"}
	tests := []struct {
		name, path string
		call       func(*ReleasesService) (any, *Response, error)
	}{
		{"installation ID", "installations/by-id/" + releaseTestID, func(s *ReleasesService) (any, *Response, error) { return s.GetInstallation(ctx, "acme", installation) }},
		{"UUID name is not ID", "installations/" + releaseTestID, func(s *ReleasesService) (any, *Response, error) {
			return s.GetInstallation(ctx, "acme", ReleaseInstallationRef{Name: releaseTestID})
		}},
		{"target ID", "targets/by-id/" + releaseTestID, func(s *ReleasesService) (any, *Response, error) {
			return s.GetTarget(ctx, "acme", ReleaseTargetRef{ID: releaseTestID})
		}},
		{"target under installation ID", "installations/by-id/" + releaseTestID + "/targets/east", func(s *ReleasesService) (any, *Response, error) { return s.GetTarget(ctx, "acme", target) }},
		{"component under target ID", "targets/by-id/" + releaseTestID + "/components/api", func(s *ReleasesService) (any, *Response, error) {
			return s.GetComponent(ctx, "acme", ReleaseComponentRef{Target: ReleaseTargetRef{ID: releaseTestID}, Name: "api"})
		}},
		{"component under installation ID", "installations/by-id/" + releaseTestID + "/targets/east/components/api", func(s *ReleasesService) (any, *Response, error) {
			return s.GetComponent(ctx, "acme", ReleaseComponentRef{Target: target, Name: "api"})
		}},
		{"component all names", "installations/prod/targets/east/components/api", func(s *ReleasesService) (any, *Response, error) {
			return s.GetComponent(ctx, "acme", ReleaseComponentRef{Target: ReleaseTargetRef{Installation: ReleaseInstallationRef{Name: "prod"}, Name: "east"}, Name: "api"})
		}},
		{"artifact ID", "artifacts/by-id/" + releaseTestID, func(s *ReleasesService) (any, *Response, error) {
			return s.GetArtifact(ctx, "acme", ReleaseArtifactRef{ID: releaseTestID})
		}},
		{"artifact UUID name", "artifacts/" + releaseTestID, func(s *ReleasesService) (any, *Response, error) {
			return s.GetArtifact(ctx, "acme", ReleaseArtifactRef{Name: releaseTestID})
		}},
		{"version metadata", "artifacts/api/versions/1.2.3%2Bbuild.4", func(s *ReleasesService) (any, *Response, error) {
			return s.GetVersion(ctx, "acme", ReleaseVersionRef{Artifact: ReleaseArtifactRef{Name: "api"}, Version: "1.2.3+build.4"})
		}},
		{"version by artifact ID", "artifacts/by-id/" + releaseTestID + "/versions/1.2.3%2Bbuild.4", func(s *ReleasesService) (any, *Response, error) {
			return s.GetVersion(ctx, "acme", ReleaseVersionRef{Artifact: ReleaseArtifactRef{ID: releaseTestID}, Version: "1.2.3+build.4"})
		}},
		{"escaped version segment", "artifacts/api/versions/label%2Fpart%3F%23%25%20%2B", func(s *ReleasesService) (any, *Response, error) {
			return s.GetVersion(ctx, "acme", ReleaseVersionRef{Artifact: ReleaseArtifactRef{Name: "api"}, Version: "label/part?#% +"})
		}},
		{"escaped organization", "../unused", func(s *ReleasesService) (any, *Response, error) { return s.ListArtifacts(ctx, "acme/name+space", nil) }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			wantPath := "/prefix/v2/organizations/acme/" + tt.path
			if tt.name == "escaped organization" {
				wantPath = "/prefix/v2/organizations/acme%2Fname%2Bspace/artifacts"
			}
			// A raw handler observes the wire path without ServeMux redirects or
			// path cleaning masking segment-escaping regressions.
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.EscapedPath() != wantPath {
					t.Errorf("wire path = %q, want %q", r.URL.EscapedPath(), wantPath)
				}
				_, _ = fmt.Fprint(w, `{}`)
			}))
			t.Cleanup(server.Close)
			client := must(NewClient(WithBaseURL(server.URL + "/prefix/")))
			if _, _, err := tt.call(client.Releases); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestReleasesService_InvalidSelectors(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	tests := []struct {
		name string
		call func(*ReleasesService) (any, *Response, error)
	}{
		{"empty installation", func(s *ReleasesService) (any, *Response, error) {
			return s.GetInstallation(ctx, "acme", ReleaseInstallationRef{})
		}},
		{"ambiguous installation", func(s *ReleasesService) (any, *Response, error) {
			return s.ListTargets(ctx, "acme", ReleaseInstallationRef{ID: releaseTestID, Name: "prod"}, nil)
		}},
		{"empty target", func(s *ReleasesService) (any, *Response, error) { return s.GetTarget(ctx, "acme", ReleaseTargetRef{}) }},
		{"missing installation", func(s *ReleasesService) (any, *Response, error) {
			return s.ListComponents(ctx, "acme", ReleaseTargetRef{Name: "east"}, nil)
		}},
		{"ambiguous target name", func(s *ReleasesService) (any, *Response, error) {
			return s.GetTarget(ctx, "acme", ReleaseTargetRef{ID: releaseTestID, Name: "east"})
		}},
		{"ambiguous target parent", func(s *ReleasesService) (any, *Response, error) {
			return s.ListComponentStatuses(ctx, "acme", ReleaseTargetRef{ID: releaseTestID, Installation: ReleaseInstallationRef{Name: "prod"}}, nil)
		}},
		{"invalid parent selector", func(s *ReleasesService) (any, *Response, error) {
			return s.ListObservations(ctx, "acme", ReleaseTargetRef{Name: "east", Installation: ReleaseInstallationRef{ID: releaseTestID, Name: "prod"}}, nil)
		}},
		{"empty component", func(s *ReleasesService) (any, *Response, error) {
			return s.GetComponent(ctx, "acme", ReleaseComponentRef{})
		}},
		{"component missing target", func(s *ReleasesService) (any, *Response, error) {
			return s.GetComponent(ctx, "acme", ReleaseComponentRef{Name: "api"})
		}},
		{"ambiguous component name", func(s *ReleasesService) (any, *Response, error) {
			return s.GetComponent(ctx, "acme", ReleaseComponentRef{ID: releaseTestID, Name: "api"})
		}},
		{"ambiguous component target", func(s *ReleasesService) (any, *Response, error) {
			return s.GetComponent(ctx, "acme", ReleaseComponentRef{ID: releaseTestID, Target: ReleaseTargetRef{ID: releaseTestID}})
		}},
		{"empty status component", func(s *ReleasesService) (any, *Response, error) {
			return s.GetComponentStatus(ctx, "acme", ReleaseTargetRef{ID: releaseTestID}, "")
		}},
		{"empty artifact", func(s *ReleasesService) (any, *Response, error) {
			return s.GetArtifact(ctx, "acme", ReleaseArtifactRef{})
		}},
		{"ambiguous artifact", func(s *ReleasesService) (any, *Response, error) {
			return s.ListVersions(ctx, "acme", ReleaseArtifactRef{ID: releaseTestID, Name: "api"}, nil)
		}},
		{"empty version", func(s *ReleasesService) (any, *Response, error) {
			return s.GetVersion(ctx, "acme", ReleaseVersionRef{})
		}},
		{"version missing artifact", func(s *ReleasesService) (any, *Response, error) {
			return s.GetVersion(ctx, "acme", ReleaseVersionRef{Version: "1.2.3"})
		}},
		{"ambiguous version label", func(s *ReleasesService) (any, *Response, error) {
			return s.GetVersion(ctx, "acme", ReleaseVersionRef{ID: releaseTestID, Version: "1.2.3"})
		}},
		{"ambiguous version parent", func(s *ReleasesService) (any, *Response, error) {
			return s.GetVersion(ctx, "acme", ReleaseVersionRef{ID: releaseTestID, Artifact: ReleaseArtifactRef{Name: "api"}})
		}},
		{"invalid version artifact", func(s *ReleasesService) (any, *Response, error) {
			return s.GetVersion(ctx, "acme", ReleaseVersionRef{Version: "1.2.3", Artifact: ReleaseArtifactRef{ID: releaseTestID, Name: "api"}})
		}},
		{"invalid comparison artifact", func(s *ReleasesService) (any, *Response, error) {
			return s.CompareVersions(ctx, "acme", ReleaseArtifactRef{}, nil)
		}},
		{"invalid PR comparison artifact", func(s *ReleasesService) (any, *Response, error) {
			return s.ListComparisonPullRequests(ctx, "acme", ReleaseArtifactRef{}, nil)
		}},
		{"create target without installation", func(s *ReleasesService) (any, *Response, error) {
			return s.CreateTarget(ctx, "acme", ReleaseInstallationRef{}, ReleaseTargetCreate{Name: "east"})
		}},
		{"update ambiguous target", func(s *ReleasesService) (any, *Response, error) {
			return s.UpdateTarget(ctx, "acme", ReleaseTargetRef{ID: releaseTestID, Name: "east"}, ReleaseTargetUpdate{})
		}},
		{"create component without target", func(s *ReleasesService) (any, *Response, error) {
			return s.CreateComponent(ctx, "acme", ReleaseTargetRef{}, ReleaseComponentCreate{Name: "api"})
		}},
		{"create component with ambiguous artifact", func(s *ReleasesService) (any, *Response, error) {
			return s.CreateComponent(ctx, "acme", ReleaseTargetRef{ID: releaseTestID}, ReleaseComponentCreate{Name: "api", Artifact: &ReleaseArtifactRef{ID: releaseTestID, Name: "api"}})
		}},
		{"create component with empty artifact", func(s *ReleasesService) (any, *Response, error) {
			return s.CreateComponent(ctx, "acme", ReleaseTargetRef{ID: releaseTestID}, ReleaseComponentCreate{Name: "api", Artifact: &ReleaseArtifactRef{}})
		}},
		{"publish without artifact", func(s *ReleasesService) (any, *Response, error) {
			return s.PublishVersion(ctx, "acme", ReleaseArtifactRef{}, ReleaseVersionPublish{})
		}},
		{"yank without version", func(s *ReleasesService) (any, *Response, error) {
			return s.YankVersion(ctx, "acme", ReleaseVersionRef{}, ReleaseVersionYank{})
		}},
		{"append without target", func(s *ReleasesService) (any, *Response, error) {
			return s.AppendObservations(ctx, "acme", ReleaseTargetRef{}, nil)
		}},
		{"archive without installation", func(s *ReleasesService) (any, *Response, error) {
			return s.ArchiveInstallation(ctx, "acme", ReleaseInstallationRef{})
		}},
		{"unarchive without component", func(s *ReleasesService) (any, *Response, error) {
			return s.UnarchiveComponent(ctx, "acme", ReleaseComponentRef{})
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server, client, teardown := newMockServerAndClient(t)
			t.Cleanup(teardown)
			_, resp, err := tt.call(client.Releases)
			if err == nil || resp != nil || len(server.calls) != 0 {
				t.Fatalf("invalid selector: response=%#v error=%v HTTP calls=%v", resp, err, server.calls)
			}
		})
	}
}

func TestReleasesService_VersionOptionalFilters(t *testing.T) {
	t.Parallel()
	suffix, yanked := "debug", true
	tests := []struct {
		name, query string
		opt         *ReleasesVersionListOptions
	}{
		{"nil options", "", nil},
		{"omitted filters", "", &ReleasesVersionListOptions{}},
		{"explicit filters", "direction=desc&older_than=baseline&suffix=debug&yanked=true", &ReleasesVersionListOptions{Direction: "desc", OlderThan: "baseline", Suffix: &suffix, Yanked: &yanked}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server, client, teardown := newMockServerAndClient(t)
			t.Cleanup(teardown)
			server.HandleFunc("/v2/organizations/acme/artifacts/api/versions", func(w http.ResponseWriter, r *http.Request) {
				if r.URL.RawQuery != tt.query {
					t.Errorf("query=%q, want %q", r.URL.RawQuery, tt.query)
				}
				_, _ = fmt.Fprint(w, `{"items":[],"links":{"self":"https://api.example/versions"}}`)
			})
			got, _, err := client.Releases.ListVersions(context.Background(), "acme", ReleaseArtifactRef{Name: "api"}, tt.opt)
			if err != nil || got.Ordering != nil || got.Items == nil || len(got.Items) != 0 {
				t.Fatalf("page=%#v, error=%v", got, err)
			}
		})
	}
}

func TestReleasesService_APIErrorsAndContext(t *testing.T) {
	t.Parallel()
	for _, status := range []int{http.StatusUnauthorized, http.StatusForbidden, http.StatusNotFound, http.StatusUnprocessableEntity} {
		t.Run(fmt.Sprint(status), func(t *testing.T) {
			t.Parallel()
			server, client, teardown := newMockServerAndClient(t)
			t.Cleanup(teardown)
			server.HandleFunc("/v2/organizations/acme/artifacts/api/version-comparison", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(status)
				_, _ = fmt.Fprint(w, `{"message":"not available"}`)
			})
			_, resp, err := client.Releases.CompareVersions(context.Background(), "acme", ReleaseArtifactRef{Name: "api"}, nil)
			var apiErr *ErrorResponse
			if !errors.As(err, &apiErr) || apiErr.Message != "not available" || resp == nil || resp.StatusCode != status {
				t.Fatalf("response=%#v error=%v", resp, err)
			}
		})
	}
	t.Run("canceled context", func(t *testing.T) {
		t.Parallel()
		_, client, teardown := newMockServerAndClient(t)
		t.Cleanup(teardown)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		_, _, err := client.Releases.ListInstallations(ctx, "acme", nil)
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("error=%v, want context.Canceled", err)
		}
	})
}

func TestReleasesService_ResponseVariants(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	target := ReleaseTargetRef{ID: releaseTestID}
	artifact := ReleaseArtifactRef{Name: "api"}
	tests := []struct {
		name, body string
		call       func(*ReleasesService) (any, *Response, error)
	}{
		{
			name: "unbound component",
			body: `{"id":"component-id","name":"api","target_id":"target-id","artifact":null,"url":"https://api.example/component","id_url":"https://api.example/component-id","created_at":"2026-09-20T01:02:03Z","archived_at":null}`,
			call: func(s *ReleasesService) (any, *Response, error) {
				return s.GetComponent(ctx, "acme", ReleaseComponentRef{ID: releaseTestID})
			},
		},
		{
			name: "artifact without publishing pipeline",
			body: `{"id":"artifact-id","name":"api","version_scheme":"semver","url":"https://api.example/artifact","id_url":"https://api.example/artifact-id","created_at":"2026-09-20T01:02:03Z","archived_at":null}`,
			call: func(s *ReleasesService) (any, *Response, error) { return s.GetArtifact(ctx, "acme", artifact) },
		},
		{
			name: "version provenance and yank state",
			body: `{"id":"version-id","version":"1.2.3","digest":"sha256:0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef","source_url":"opaque source","source_build":{"id":"build-id"},"pull_requests":[],"created_by":{"id":"user-id","graphql_id":"graphql-id","name":null,"email":"person@example.com","avatar_url":null,"created_at":"2026-09-20T01:02:03Z"},"url":"https://api.example/version","id_url":"https://api.example/version-id","created_at":"2026-09-20T01:02:03Z","yanked":true,"yanked_at":"2026-09-21T01:02:03Z","yank_reason":"withdrawn"}`,
			call: func(s *ReleasesService) (any, *Response, error) {
				return s.GetVersion(ctx, "acme", ReleaseVersionRef{ID: releaseTestID})
			},
		},
		{
			name: "never observed status remains null",
			body: `{"target_id":"target-id","component":"api","component_id":"component-id","present":null,"active_version":null,"active_digest":null,"last_active_version":null,"first_seen_at":null,"last_seen_at":null,"sequence":null,"catalog_resolution":{"status":"unknown","artifact_version_id":null,"digest_check":"not_checked","reason":"no_active_version"}}`,
			call: func(s *ReleasesService) (any, *Response, error) {
				return s.GetComponentStatus(ctx, "acme", target, "api")
			},
		},
		{
			name: "resolved present status",
			body: `{"target_id":"target-id","component":"api","component_id":"component-id","present":true,"active_version":"1.2.3","active_digest":"sha256:0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef","last_active_version":"1.2.3","first_seen_at":"2026-09-20T01:02:03Z","last_seen_at":"2026-09-21T01:02:03Z","sequence":9007199254740993,"catalog_resolution":{"status":"resolved","artifact_version_id":"version-id","digest_check":"match","reason":null}}`,
			call: func(s *ReleasesService) (any, *Response, error) {
				return s.GetComponentStatus(ctx, "acme", target, "api")
			},
		},
		{
			name: "unknown comparison with known different content",
			body: `{"artifact":{"id":"artifact-id","name":"api","version_scheme":"sha"},"left":{"id":"left-id","version":"abcdef0","digest":"sha256:abc","yanked":false},"right":{"id":"right-id","version":"abcdef1","digest":"sha256:def","yanked":false},"basis":null,"result":"unknown","reason":"no_precedence_for_scheme","same_version":false,"same_content":false}`,
			call: func(s *ReleasesService) (any, *Response, error) {
				return s.CompareVersions(ctx, "acme", artifact, &ReleasesComparisonOptions{LeftID: "left-id", RightID: "right-id"})
			},
		},
		{
			name: "component observation payload presence",
			body: `{"items":[
				{"id":"one","target_id":"target-id","sequence":1,"event_id":"event-one","type":"ComponentAppeared","component":"api","occurred_at":"2026-09-20T01:02:03Z","recorded_at":"2026-09-20T01:03:04Z"},
				{"id":"two","target_id":"target-id","sequence":2,"event_id":"event-two","type":"ComponentDisappeared","component":"api","occurred_at":"2026-09-20T01:02:03Z","recorded_at":"2026-09-20T01:03:04Z"},
				{"id":"three","target_id":"target-id","sequence":3,"event_id":"event-three","type":"ComponentVersionObserved","component":"api","version":"1.2.3","occurred_at":"2026-09-20T01:02:03Z","recorded_at":"2026-09-20T01:03:04Z","catalog_resolution":{"status":"resolved","artifact_version_id":"version-id","digest_check":"not_reported","reason":null}},
				{"id":"four","target_id":"target-id","sequence":4,"event_id":"event-four","type":"ComponentVersionObserved","component":"api","version":"1.2.4","digest":"sha256:0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef","occurred_at":"2026-09-20T01:02:03Z","recorded_at":"2026-09-20T01:03:04Z","catalog_resolution":{"status":"unknown","artifact_version_id":null,"digest_check":"not_checked","reason":"version_not_registered"}}
			],"links":{"self":"https://api.example/observations"}}`,
			call: func(s *ReleasesService) (any, *Response, error) { return s.ListObservations(ctx, "acme", target, nil) },
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_, _ = fmt.Fprint(w, tt.body)
			}))
			t.Cleanup(server.Close)
			client := must(NewClient(WithBaseURL(server.URL)))
			got, _, err := tt.call(client.Releases)
			if err != nil {
				t.Fatal(err)
			}
			encoded := must(json.Marshal(got))
			assertJSONEqual(t, string(encoded), tt.body)
			if status, ok := got.(ReleaseComponentStatus); ok && status.Sequence != nil && *status.Sequence != 9007199254740993 {
				t.Errorf("sequence lost integer precision: %d", *status.Sequence)
			}
		})
	}
}
