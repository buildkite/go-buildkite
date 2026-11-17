package buildkite

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
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

var (
	releaseCreatedAt          = NewTimestamp(time.Date(2026, 9, 20, 1, 2, 3, 0, time.UTC))
	releaseLaterAt            = NewTimestamp(time.Date(2026, 9, 21, 1, 2, 3, 0, time.UTC))
	releaseRecordedAt         = NewTimestamp(time.Date(2026, 9, 20, 1, 3, 4, 0, time.UTC))
	releaseDigest             = releaseTestDigest
	releaseAbsent             = false
	releaseSequence     int64 = 9007199254740993
	releaseAbsentReason       = "instance_absent"
	releaseTarget             = ReleaseTarget{
		ID: "target-id", GraphQLID: "opaque-target", Name: "east",
		Labels: map[string]string{"region": "au", "empty": "", "space": " "},
		URL:    "https://api.example/target", IDURL: "https://api.example/target-id",
		CreatedAt: releaseCreatedAt, ArchivedAt: releaseLaterAt,
	}
	releaseComponent = ReleaseComponent{
		ID: "component-id", GraphQLID: "opaque-component", Name: "api",
		Artifact: ReleaseArtifactSummary{ID: "artifact-id", GraphQLID: "opaque-artifact", Name: "service"},
		URL:      "https://api.example/component", IDURL: "https://api.example/component-id", CreatedAt: releaseCreatedAt,
	}
	releaseInstance = ReleaseComponentInstance{
		ID: "instance-id", GraphQLID: "opaque-instance", Name: "api-blue",
		TargetID: "target-id", ComponentID: "component-id",
		URL: "https://api.example/instance", IDURL: "https://api.example/instance-id", CreatedAt: releaseCreatedAt,
	}
	releaseArtifact = ReleaseArtifact{
		ID: "artifact-id", GraphQLID: "opaque-artifact", Name: "service",
		URL: "https://api.example/artifact", IDURL: "https://api.example/artifact-id", CreatedAt: releaseCreatedAt,
	}
	releaseVersion = ReleaseArtifactVersion{
		ID: "digest-id", GraphQLID: "opaque-digest", Digest: &releaseDigest,
		PullRequests: []ReleasePullRequest{
			{URL: "https://github.com/acme/api/pull/12", Provider: "github", RepositoryURL: "https://github.com/acme/api", ID: "12"},
			{URL: "https://example.com/pr/2"},
		},
		URL: "https://api.example/digest", IDURL: "https://api.example/digest-id", CreatedAt: releaseCreatedAt,
	}
	releaseStatus = ReleaseComponentInstanceStatus{
		TargetID: "target-id", Instance: "api-blue", InstanceID: "instance-id", Present: &releaseAbsent,
		LastActiveDigest: &releaseDigest, FirstSeenAt: releaseCreatedAt, LastSeenAt: releaseLaterAt, Sequence: &releaseSequence,
		CatalogResolution: ReleaseCatalogResolution{Status: "unknown", Reason: &releaseAbsentReason},
	}
	releaseObservation = ReleaseObservation{
		ID: "observation-id", GraphQLID: "opaque-observation", TargetID: "target-id", Sequence: 9007199254740993,
		EventID: releaseTestID, Type: "ComponentInstanceVersionObserved", Instance: "api-blue", Digest: &releaseDigest,
		OccurredAt: releaseCreatedAt, RecordedAt: releaseRecordedAt,
		CatalogResolution: &ReleaseCatalogResolution{Status: "resolved"},
	}
	releaseLinks = ReleasePageLinks{
		Self: "https://api.example/self", First: "https://api.example/first",
		Previous: "https://api.example/prev", Next: "https://api.example/next?after=opaque%2Bcursor",
	}
)

func assertReleaseResponse(t *testing.T, got, want any, fixture string) {
	t.Helper()
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("response diff (-want +got):\n%s", diff)
	}
	assertJSONEqual(t, string(must(json.Marshal(got))), fixture)
}

func TestReleasesService_ListTargets(t *testing.T) {
	t.Parallel()
	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)
	server.HandleFunc("/v2/organizations/acme/targets", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testFormValues(t, r, values{"after": "opaque+cursor/=", "per_page": "25"})
		_, _ = fmt.Fprint(w, releaseTestPage(releaseTargetJSON))
	})
	got, _, err := client.Releases.ListTargets(context.Background(), "acme", &ReleasesListOptions{After: "opaque+cursor/=", PerPage: 25})
	if err != nil {
		t.Fatal(err)
	}
	want := ReleasesPage[ReleaseTarget]{Items: []ReleaseTarget{releaseTarget}, Links: releaseLinks}
	assertReleaseResponse(t, got, want, releaseTestPage(releaseTargetJSON))
}

func TestReleasesService_GetTarget(t *testing.T) {
	t.Parallel()
	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)
	server.HandleFunc("/v2/organizations/acme/targets/"+releaseTestID, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		_, _ = fmt.Fprint(w, releaseTargetJSON)
	})
	got, _, err := client.Releases.GetTarget(context.Background(), "acme", ReleaseTargetRef{Name: releaseTestID})
	if err != nil {
		t.Fatal(err)
	}
	assertReleaseResponse(t, got, releaseTarget, releaseTargetJSON)
}

func TestReleasesService_ListComponents(t *testing.T) {
	t.Parallel()
	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)
	server.HandleFunc("/v2/organizations/acme/components", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testFormValues(t, r, values{"after": "opaque+cursor/=", "per_page": "25"})
		_, _ = fmt.Fprint(w, releaseTestPage(releaseComponentJSON))
	})
	got, _, err := client.Releases.ListComponents(context.Background(), "acme", &ReleasesListOptions{After: "opaque+cursor/=", PerPage: 25})
	if err != nil {
		t.Fatal(err)
	}
	want := ReleasesPage[ReleaseComponent]{Items: []ReleaseComponent{releaseComponent}, Links: releaseLinks}
	assertReleaseResponse(t, got, want, releaseTestPage(releaseComponentJSON))
}

func TestReleasesService_GetComponent(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		name string
		ref  ReleaseComponentRef
		path string
	}{
		{name: "name", ref: ReleaseComponentRef{Name: "api"}, path: "components/api"},
		{name: "ID", ref: ReleaseComponentRef{ID: releaseTestID}, path: "components/by-id/" + releaseTestID},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server, client, teardown := newMockServerAndClient(t)
			t.Cleanup(teardown)
			server.HandleFunc("/v2/organizations/acme/"+tt.path, func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, http.MethodGet)
				_, _ = fmt.Fprint(w, releaseComponentJSON)
			})
			got, _, err := client.Releases.GetComponent(context.Background(), "acme", tt.ref)
			if err != nil {
				t.Fatal(err)
			}
			assertReleaseResponse(t, got, releaseComponent, releaseComponentJSON)
		})
	}
}

func TestReleasesService_ListComponentInstances(t *testing.T) {
	t.Parallel()
	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)
	server.HandleFunc("/v2/organizations/acme/targets/by-id/"+releaseTestID+"/instances", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testFormValues(t, r, values{"after": "opaque+cursor/=", "per_page": "25", "label[region]": "au +", "label[empty]": ""})
		_, _ = fmt.Fprint(w, releaseTestPage(releaseInstanceJSON))
	})
	got, _, err := client.Releases.ListComponentInstances(context.Background(), "acme", ReleaseTargetRef{ID: releaseTestID}, &ReleasesComponentInstanceListOptions{
		ReleasesListOptions: ReleasesListOptions{After: "opaque+cursor/=", PerPage: 25},
		Label:               ReleaseTargetLabelFilters{"region": "au +", "empty": ""},
	})
	if err != nil {
		t.Fatal(err)
	}
	want := ReleasesPage[ReleaseComponentInstance]{Items: []ReleaseComponentInstance{releaseInstance}, Links: releaseLinks}
	assertReleaseResponse(t, got, want, releaseTestPage(releaseInstanceJSON))
}

func TestReleasesService_ListOrganizationComponentInstances(t *testing.T) {
	t.Parallel()
	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)
	server.HandleFunc("/v2/organizations/acme/component-instances", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testFormValues(t, r, values{"after": "opaque+cursor/=", "per_page": "25", "label[region]": "au +", "label[empty]": ""})
		_, _ = fmt.Fprint(w, releaseTestPage(releaseInstanceJSON))
	})
	got, _, err := client.Releases.ListOrganizationComponentInstances(context.Background(), "acme", &ReleasesComponentInstanceListOptions{
		ReleasesListOptions: ReleasesListOptions{After: "opaque+cursor/=", PerPage: 25},
		Label:               ReleaseTargetLabelFilters{"region": "au +", "empty": ""},
	})
	if err != nil {
		t.Fatal(err)
	}
	want := ReleasesPage[ReleaseComponentInstance]{Items: []ReleaseComponentInstance{releaseInstance}, Links: releaseLinks}
	assertReleaseResponse(t, got, want, releaseTestPage(releaseInstanceJSON))
}

func TestReleasesService_GetComponentInstance(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		name string
		ref  ReleaseComponentInstanceRef
		path string
	}{
		{name: "Target ID", ref: ReleaseComponentInstanceRef{Target: ReleaseTargetRef{ID: releaseTestID}, Name: "api-blue"}, path: "targets/by-id/" + releaseTestID + "/instances/api-blue"},
		{name: "Instance ID", ref: ReleaseComponentInstanceRef{ID: releaseTestID}, path: "component-instances/by-id/" + releaseTestID},
		{name: "Target name", ref: ReleaseComponentInstanceRef{Target: ReleaseTargetRef{Name: "east"}, Name: "by-label"}, path: "targets/east/instances/by-label"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server, client, teardown := newMockServerAndClient(t)
			t.Cleanup(teardown)
			server.HandleFunc("/v2/organizations/acme/"+tt.path, func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, http.MethodGet)
				_, _ = fmt.Fprint(w, releaseInstanceJSON)
			})
			got, _, err := client.Releases.GetComponentInstance(context.Background(), "acme", tt.ref)
			if err != nil {
				t.Fatal(err)
			}
			assertReleaseResponse(t, got, releaseInstance, releaseInstanceJSON)
		})
	}
}

func TestReleasesService_ListComponentInstanceStatuses(t *testing.T) {
	t.Parallel()
	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)
	server.HandleFunc("/v2/organizations/acme/targets/by-id/"+releaseTestID+"/instance-statuses", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testFormValues(t, r, values{"after": "opaque+cursor/=", "per_page": "25"})
		_, _ = fmt.Fprint(w, releaseTestPage(releaseStatusJSON))
	})
	got, _, err := client.Releases.ListComponentInstanceStatuses(context.Background(), "acme", ReleaseTargetRef{ID: releaseTestID}, &ReleasesListOptions{After: "opaque+cursor/=", PerPage: 25})
	if err != nil {
		t.Fatal(err)
	}
	want := ReleasesPage[ReleaseComponentInstanceStatus]{Items: []ReleaseComponentInstanceStatus{releaseStatus}, Links: releaseLinks}
	assertReleaseResponse(t, got, want, releaseTestPage(releaseStatusJSON))
}

func TestReleasesService_ListObservations(t *testing.T) {
	t.Parallel()
	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)
	server.HandleFunc("/v2/organizations/acme/targets/by-id/"+releaseTestID+"/observations", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testFormValues(t, r, values{"before": "prior", "instance": "api-blue", "type": "ComponentInstanceVersionObserved"})
		_, _ = fmt.Fprint(w, releaseTestPage(releaseObservationJSON))
	})
	got, _, err := client.Releases.ListObservations(context.Background(), "acme", ReleaseTargetRef{ID: releaseTestID}, &ReleasesObservationListOptions{
		ReleasesListOptions: ReleasesListOptions{Before: "prior"}, Instance: "api-blue", Type: "ComponentInstanceVersionObserved",
	})
	if err != nil {
		t.Fatal(err)
	}
	want := ReleasesPage[ReleaseObservation]{Items: []ReleaseObservation{releaseObservation}, Links: releaseLinks}
	assertReleaseResponse(t, got, want, releaseTestPage(releaseObservationJSON))
}

func TestReleasesService_ListArtifacts(t *testing.T) {
	t.Parallel()
	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)
	server.HandleFunc("/v2/organizations/acme/artifacts", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testFormValues(t, r, values{"after": "opaque+cursor/=", "per_page": "25"})
		_, _ = fmt.Fprint(w, releaseTestPage(releaseArtifactJSON))
	})
	got, _, err := client.Releases.ListArtifacts(context.Background(), "acme", &ReleasesListOptions{After: "opaque+cursor/=", PerPage: 25})
	if err != nil {
		t.Fatal(err)
	}
	want := ReleasesPage[ReleaseArtifact]{Items: []ReleaseArtifact{releaseArtifact}, Links: releaseLinks}
	assertReleaseResponse(t, got, want, releaseTestPage(releaseArtifactJSON))
}

func TestReleasesService_GetArtifact(t *testing.T) {
	t.Parallel()
	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)
	server.HandleFunc("/v2/organizations/acme/artifacts/service", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		_, _ = fmt.Fprint(w, releaseArtifactJSON)
	})
	got, _, err := client.Releases.GetArtifact(context.Background(), "acme", ReleaseArtifactRef{Name: "service"})
	if err != nil {
		t.Fatal(err)
	}
	assertReleaseResponse(t, got, releaseArtifact, releaseArtifactJSON)
}

func TestReleasesService_ListVersions(t *testing.T) {
	t.Parallel()
	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)
	server.HandleFunc("/v2/organizations/acme/artifacts/by-id/"+releaseTestID+"/versions", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testFormValues(t, r, values{"after": "opaque+cursor/=", "per_page": "25"})
		_, _ = fmt.Fprint(w, releaseTestPage(releaseVersionJSON))
	})
	got, _, err := client.Releases.ListVersions(context.Background(), "acme", ReleaseArtifactRef{ID: releaseTestID}, &ReleasesListOptions{After: "opaque+cursor/=", PerPage: 25})
	if err != nil {
		t.Fatal(err)
	}
	want := ReleasesPage[ReleaseArtifactVersion]{Items: []ReleaseArtifactVersion{releaseVersion}, Links: releaseLinks}
	assertReleaseResponse(t, got, want, releaseTestPage(releaseVersionJSON))
}

func TestReleasesService_CreateTarget(t *testing.T) {
	t.Parallel()
	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)
	server.HandleFunc("/v2/organizations/acme/targets", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		assertRequestJSON(t, r, `{"name":"east","labels":{"region":"au"}}`)
		w.WriteHeader(http.StatusCreated)
		_, _ = fmt.Fprint(w, releaseTargetJSON)
	})
	got, resp, err := client.Releases.CreateTarget(context.Background(), "acme", ReleaseTargetCreate{Name: "east", Labels: map[string]string{"region": "au"}})
	if err != nil || resp == nil || resp.StatusCode != http.StatusCreated {
		t.Fatalf("response=%#v, error=%v", resp, err)
	}
	assertReleaseResponse(t, got, releaseTarget, releaseTargetJSON)
}

func TestReleasesService_UpdateTarget(t *testing.T) {
	t.Parallel()
	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)
	server.HandleFunc("/v2/organizations/acme/targets/by-id/"+releaseTestID, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		assertRequestJSON(t, r, `{"labels":{}}`)
		_, _ = fmt.Fprint(w, releaseTargetJSON)
	})
	got, _, err := client.Releases.UpdateTarget(context.Background(), "acme", ReleaseTargetRef{ID: releaseTestID}, ReleaseTargetUpdate{Labels: Some(map[string]string{})})
	if err != nil {
		t.Fatal(err)
	}
	assertReleaseResponse(t, got, releaseTarget, releaseTargetJSON)
}

func TestReleasesService_CreateComponent(t *testing.T) {
	t.Parallel()
	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)
	server.HandleFunc("/v2/organizations/acme/components", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		assertRequestJSON(t, r, `{"name":"api","artifact":{"id":"artifact-id"}}`)
		w.WriteHeader(http.StatusCreated)
		_, _ = fmt.Fprint(w, releaseComponentJSON)
	})
	got, resp, err := client.Releases.CreateComponent(context.Background(), "acme", ReleaseComponentCreate{Name: "api", Artifact: ReleaseArtifactRef{ID: "artifact-id"}})
	if err != nil || resp == nil || resp.StatusCode != http.StatusCreated {
		t.Fatalf("response=%#v, error=%v", resp, err)
	}
	assertReleaseResponse(t, got, releaseComponent, releaseComponentJSON)
}

func TestReleasesService_CreateComponentInstance(t *testing.T) {
	t.Parallel()
	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)
	server.HandleFunc("/v2/organizations/acme/targets/by-id/"+releaseTestID+"/instances", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		assertRequestJSON(t, r, `{"component_name":"api","name":"api-blue"}`)
		_, _ = fmt.Fprint(w, releaseInstanceJSON)
	})
	got, _, err := client.Releases.CreateComponentInstance(context.Background(), "acme", ReleaseTargetRef{ID: releaseTestID}, ReleaseComponentInstanceCreate{ComponentName: "api", Name: "api-blue"})
	if err != nil {
		t.Fatal(err)
	}
	assertReleaseResponse(t, got, releaseInstance, releaseInstanceJSON)
}

func TestReleasesService_RegisterComponentInstance(t *testing.T) {
	t.Parallel()
	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)
	server.HandleFunc("/v2/organizations/acme/component-instances", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		assertRequestJSON(t, r, `{"component_name":"api","target_id":"`+releaseTestID+`"}`)
		_, _ = fmt.Fprint(w, releaseInstanceJSON)
	})
	got, _, err := client.Releases.RegisterComponentInstance(context.Background(), "acme", ReleaseComponentInstanceRegister{
		ReleaseComponentInstanceCreate: ReleaseComponentInstanceCreate{ComponentName: "api"}, TargetID: releaseTestID,
	})
	if err != nil {
		t.Fatal(err)
	}
	assertReleaseResponse(t, got, releaseInstance, releaseInstanceJSON)
}

func TestReleasesService_CreateArtifact(t *testing.T) {
	t.Parallel()
	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)
	server.HandleFunc("/v2/organizations/acme/artifacts", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		assertRequestJSON(t, r, `{"name":"service"}`)
		w.WriteHeader(http.StatusCreated)
		_, _ = fmt.Fprint(w, releaseArtifactJSON)
	})
	got, resp, err := client.Releases.CreateArtifact(context.Background(), "acme", ReleaseArtifactCreate{Name: "service"})
	if err != nil || resp == nil || resp.StatusCode != http.StatusCreated {
		t.Fatalf("response=%#v, error=%v", resp, err)
	}
	assertReleaseResponse(t, got, releaseArtifact, releaseArtifactJSON)
}

func TestReleasesService_PublishVersion(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		name     string
		artifact ReleaseArtifactRef
		path     string
		input    ReleaseVersionPublish
		request  string
	}{
		{
			name: "Artifact name", artifact: ReleaseArtifactRef{Name: "service"}, path: "artifacts/service/versions",
			input:   ReleaseVersionPublish{Digest: releaseTestDigest, SourceBuild: Some(ReleaseReference{ID: "build-id"}), PullRequests: []ReleasePullRequest{{URL: "https://example.com/pr/2"}}},
			request: `{"digest":"` + releaseTestDigest + `","source_build":{"id":"build-id"},"pull_requests":[{"url":"https://example.com/pr/2"}]}`,
		},
		{
			name: "Artifact ID", artifact: ReleaseArtifactRef{ID: releaseTestID}, path: "artifacts/by-id/" + releaseTestID + "/versions",
			input: ReleaseVersionPublish{Digest: releaseTestDigest}, request: `{"digest":"` + releaseTestDigest + `"}`,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server, client, teardown := newMockServerAndClient(t)
			t.Cleanup(teardown)
			server.HandleFunc("/v2/organizations/acme/"+tt.path, func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, http.MethodPost)
				if r.Header.Get("Content-Type") != "application/json" {
					t.Error("missing JSON content type")
				}
				assertRequestJSON(t, r, tt.request)
				w.WriteHeader(http.StatusCreated)
				_, _ = fmt.Fprint(w, releaseVersionJSON)
			})
			got, resp, err := client.Releases.PublishVersion(context.Background(), "acme", tt.artifact, tt.input)
			if err != nil || resp == nil || resp.StatusCode != http.StatusCreated {
				t.Fatalf("response=%#v, error=%v", resp, err)
			}
			assertReleaseResponse(t, got, releaseVersion, releaseVersionJSON)
		})
	}
}

func TestReleasesService_YankVersion(t *testing.T) {
	t.Parallel()
	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)
	server.HandleFunc("/v2/organizations/acme/artifact-versions/by-id/"+releaseTestID+"/yank", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		assertRequestJSON(t, r, `{"reason":"withdrawn"}`)
		_, _ = fmt.Fprint(w, releaseVersionJSON)
	})
	got, _, err := client.Releases.YankVersion(context.Background(), "acme", ReleaseVersionRef{ID: releaseTestID}, ReleaseVersionYank{Reason: "withdrawn"})
	if err != nil {
		t.Fatal(err)
	}
	assertReleaseResponse(t, got, releaseVersion, releaseVersionJSON)
}

func TestReleasesService_ArchiveTarget(t *testing.T) {
	t.Parallel()
	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)
	server.HandleFunc("/v2/organizations/acme/targets/by-id/"+releaseTestID+"/archive", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		assertRequestJSON(t, r, `{}`)
		_, _ = fmt.Fprint(w, releaseTargetJSON)
	})
	got, _, err := client.Releases.ArchiveTarget(context.Background(), "acme", ReleaseTargetRef{ID: releaseTestID})
	if err != nil {
		t.Fatal(err)
	}
	assertReleaseResponse(t, got, releaseTarget, releaseTargetJSON)
}

func TestReleasesService_UnarchiveTarget(t *testing.T) {
	t.Parallel()
	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)
	server.HandleFunc("/v2/organizations/acme/targets/by-id/"+releaseTestID+"/unarchive", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		assertRequestJSON(t, r, `{}`)
		_, _ = fmt.Fprint(w, releaseTargetJSON)
	})
	got, _, err := client.Releases.UnarchiveTarget(context.Background(), "acme", ReleaseTargetRef{ID: releaseTestID})
	if err != nil {
		t.Fatal(err)
	}
	assertReleaseResponse(t, got, releaseTarget, releaseTargetJSON)
}

func TestReleasesService_ArchiveComponent(t *testing.T) {
	t.Parallel()
	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)
	server.HandleFunc("/v2/organizations/acme/components/api/archive", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		assertRequestJSON(t, r, `{}`)
		_, _ = fmt.Fprint(w, releaseComponentJSON)
	})
	got, _, err := client.Releases.ArchiveComponent(context.Background(), "acme", ReleaseComponentRef{Name: "api"})
	if err != nil {
		t.Fatal(err)
	}
	assertReleaseResponse(t, got, releaseComponent, releaseComponentJSON)
}

func TestReleasesService_UnarchiveComponent(t *testing.T) {
	t.Parallel()
	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)
	server.HandleFunc("/v2/organizations/acme/components/api/unarchive", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		assertRequestJSON(t, r, `{}`)
		_, _ = fmt.Fprint(w, releaseComponentJSON)
	})
	got, _, err := client.Releases.UnarchiveComponent(context.Background(), "acme", ReleaseComponentRef{Name: "api"})
	if err != nil {
		t.Fatal(err)
	}
	assertReleaseResponse(t, got, releaseComponent, releaseComponentJSON)
}

func TestReleasesService_ArchiveComponentInstance(t *testing.T) {
	t.Parallel()
	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)
	server.HandleFunc("/v2/organizations/acme/component-instances/by-id/"+releaseTestID+"/archive", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		assertRequestJSON(t, r, `{}`)
		_, _ = fmt.Fprint(w, releaseInstanceJSON)
	})
	got, _, err := client.Releases.ArchiveComponentInstance(context.Background(), "acme", ReleaseComponentInstanceRef{ID: releaseTestID})
	if err != nil {
		t.Fatal(err)
	}
	assertReleaseResponse(t, got, releaseInstance, releaseInstanceJSON)
}

func TestReleasesService_UnarchiveComponentInstance(t *testing.T) {
	t.Parallel()
	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)
	server.HandleFunc("/v2/organizations/acme/targets/by-id/"+releaseTestID+"/instances/api-blue/unarchive", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		assertRequestJSON(t, r, `{}`)
		_, _ = fmt.Fprint(w, releaseInstanceJSON)
	})
	got, _, err := client.Releases.UnarchiveComponentInstance(context.Background(), "acme", ReleaseComponentInstanceRef{Target: ReleaseTargetRef{ID: releaseTestID}, Name: "api-blue"})
	if err != nil {
		t.Fatal(err)
	}
	assertReleaseResponse(t, got, releaseInstance, releaseInstanceJSON)
}

func TestReleasesService_ArchiveArtifact(t *testing.T) {
	t.Parallel()
	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)
	server.HandleFunc("/v2/organizations/acme/artifacts/service/archive", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		assertRequestJSON(t, r, `{}`)
		_, _ = fmt.Fprint(w, releaseArtifactJSON)
	})
	got, _, err := client.Releases.ArchiveArtifact(context.Background(), "acme", ReleaseArtifactRef{Name: "service"})
	if err != nil {
		t.Fatal(err)
	}
	assertReleaseResponse(t, got, releaseArtifact, releaseArtifactJSON)
}

func TestReleasesService_UnarchiveArtifact(t *testing.T) {
	t.Parallel()
	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)
	server.HandleFunc("/v2/organizations/acme/artifacts/by-id/"+releaseTestID+"/unarchive", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		assertRequestJSON(t, r, `{}`)
		_, _ = fmt.Fprint(w, releaseArtifactJSON)
	})
	got, _, err := client.Releases.UnarchiveArtifact(context.Background(), "acme", ReleaseArtifactRef{ID: releaseTestID})
	if err != nil {
		t.Fatal(err)
	}
	assertReleaseResponse(t, got, releaseArtifact, releaseArtifactJSON)
}

func TestReleasesService_PublicationFailureIsNotRetried(t *testing.T) {
	t.Parallel()
	for _, status := range []int{500, 201} {
		t.Run(fmt.Sprint(status), func(t *testing.T) {
			calls := 0
			server, client, teardown := newMockServerAndClient(t)
			t.Cleanup(teardown)
			server.HandleFunc("/v2/organizations/acme/artifacts/service/versions", func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, http.MethodPost)
				calls++
				w.WriteHeader(status)
				if status == 201 {
					_, _ = fmt.Fprint(w, `{"id":`)
				} else {
					_, _ = fmt.Fprint(w, `{"message":"outcome unknown"}`)
				}
			})
			_, _, err := client.Releases.PublishVersion(context.Background(), "acme", ReleaseArtifactRef{Name: "service"}, ReleaseVersionPublish{Digest: releaseTestDigest})
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
		name, path, request string
		call                func(*ReleasesService) (any, *Response, error)
	}{
		{
			name: "omitted labels", path: "targets", request: `{"name":"east"}`,
			call: func(s *ReleasesService) (any, *Response, error) {
				return s.CreateTarget(ctx, "acme", ReleaseTargetCreate{Name: "east"})
			},
		},
		{
			name: "unchanged labels", path: "targets/east", request: `{}`,
			call: func(s *ReleasesService) (any, *Response, error) {
				return s.UpdateTarget(ctx, "acme", ReleaseTargetRef{Name: "east"}, ReleaseTargetUpdate{})
			},
		},
		{
			name: "default Instance name", path: "targets/east/instances", request: `{"component_name":"api"}`,
			call: func(s *ReleasesService) (any, *Response, error) {
				return s.CreateComponentInstance(ctx, "acme", ReleaseTargetRef{Name: "east"}, ReleaseComponentInstanceCreate{ComponentName: "api"})
			},
		},
		{
			name: "explicit null provenance", path: "artifacts/service/versions", request: `{"digest":"` + releaseTestDigest + `","source_build":null,"source_url":null}`,
			call: func(s *ReleasesService) (any, *Response, error) {
				return s.PublishVersion(ctx, "acme", ReleaseArtifactRef{Name: "service"}, ReleaseVersionPublish{Digest: releaseTestDigest, SourceBuild: Null[ReleaseReference](), SourceURL: Null[string]()})
			},
		},
		{
			name: "empty source", path: "artifacts/service/versions", request: `{"digest":"` + releaseTestDigest + `","source_url":""}`,
			call: func(s *ReleasesService) (any, *Response, error) {
				return s.PublishVersion(ctx, "acme", ReleaseArtifactRef{Name: "service"}, ReleaseVersionPublish{Digest: releaseTestDigest, SourceURL: Some("")})
			},
		},
		{
			name: "omitted yank reason", path: "artifact-versions/by-id/" + releaseTestID + "/yank", request: `{}`,
			call: func(s *ReleasesService) (any, *Response, error) {
				return s.YankVersion(ctx, "acme", ReleaseVersionRef{ID: releaseTestID}, ReleaseVersionYank{})
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server, client, teardown := newMockServerAndClient(t)
			t.Cleanup(teardown)
			server.HandleFunc("/v2/organizations/acme/"+tt.path, func(w http.ResponseWriter, r *http.Request) {
				assertRequestJSON(t, r, tt.request)
				_, _ = fmt.Fprint(w, `{}`)
			})
			_, _, err := tt.call(client.Releases)
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
			server, client, teardown := newMockServerAndClient(t)
			t.Cleanup(teardown)
			server.HandleFunc("/v2/organizations/acme/targets/east/observations", func(w http.ResponseWriter, r *http.Request) {
				calls++
				testMethod(t, r, http.MethodPost)
				assertRequestJSON(t, r, `{"events":[{"event_id":"`+releaseTestID+`","type":"ComponentInstanceVersionObserved","instance":"api-blue","digest":"`+releaseTestDigest+`"},{"event_id":"019a1dd8-a49b-70c5-99b8-b8f9f0ed6825","type":"ComponentInstanceDisappeared","instance":"api-blue"}]}`)
				w.WriteHeader(status)
				if status < 300 {
					_, _ = fmt.Fprint(w, response)
				} else {
					_, _ = fmt.Fprint(w, `{"message":"observation rejected","errors":[{"field":"/events/0/event_id","code":"conflict"}]}`)
				}
			})
			got, resp, err := client.Releases.AppendObservations(context.Background(), "acme", ReleaseTargetRef{Name: "east"}, []ReleaseObservationCreate{{EventID: releaseTestID, Type: "ComponentInstanceVersionObserved", Instance: "api-blue", Digest: releaseTestDigest}, {EventID: "019a1dd8-a49b-70c5-99b8-b8f9f0ed6825", Type: "ComponentInstanceDisappeared", Instance: "api-blue"}})
			if resp == nil || resp.StatusCode != status || calls != 1 {
				t.Fatalf("response=%#v, calls=%d, error=%v", resp, calls, err)
			}
			if status < 300 {
				if err != nil {
					t.Fatal(err)
				}
				want := []ReleaseObservationWriteResult{
					{ReleaseObservation: releaseObservation, Created: false},
					{
						ReleaseObservation: ReleaseObservation{
							ID: "second", TargetID: "target-id", Sequence: 9007199254740994,
							EventID: "019a1dd8-a49b-70c5-99b8-b8f9f0ed6825", Type: "ComponentInstanceDisappeared", Instance: "api-blue",
							OccurredAt: releaseCreatedAt, RecordedAt: releaseRecordedAt,
						},
						Created: true,
					},
				}
				assertReleaseResponse(t, got, want, response)
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
	for _, tt := range []struct {
		name string
		call func(*ReleasesService) (any, *Response, error)
	}{
		{name: "target missing selector", call: func(s *ReleasesService) (any, *Response, error) {
			return s.GetTarget(ctx, "acme", ReleaseTargetRef{})
		}},
		{name: "target conflicting selectors", call: func(s *ReleasesService) (any, *Response, error) {
			return s.GetTarget(ctx, "acme", ReleaseTargetRef{ID: releaseTestID, Name: "east"})
		}},
		{name: "component missing selector", call: func(s *ReleasesService) (any, *Response, error) {
			return s.GetComponent(ctx, "acme", ReleaseComponentRef{})
		}},
		{name: "component conflicting selectors", call: func(s *ReleasesService) (any, *Response, error) {
			return s.GetComponent(ctx, "acme", ReleaseComponentRef{ID: releaseTestID, Name: "api"})
		}},
		{name: "component instance missing target", call: func(s *ReleasesService) (any, *Response, error) {
			return s.GetComponentInstance(ctx, "acme", ReleaseComponentInstanceRef{Name: "api"})
		}},
		{name: "component instance ID with target", call: func(s *ReleasesService) (any, *Response, error) {
			return s.GetComponentInstance(ctx, "acme", ReleaseComponentInstanceRef{ID: releaseTestID, Target: ReleaseTargetRef{Name: "east"}})
		}},
		{name: "versions conflicting artifact selectors", call: func(s *ReleasesService) (any, *Response, error) {
			return s.ListVersions(ctx, "acme", ReleaseArtifactRef{ID: releaseTestID, Name: "api"}, nil)
		}},
		{name: "version missing ID", call: func(s *ReleasesService) (any, *Response, error) {
			return s.GetVersion(ctx, "acme", ReleaseVersionRef{})
		}},
		{name: "create component missing artifact", call: func(s *ReleasesService) (any, *Response, error) {
			return s.CreateComponent(ctx, "acme", ReleaseComponentCreate{Name: "api"})
		}},
		{name: "create component conflicting artifact selectors", call: func(s *ReleasesService) (any, *Response, error) {
			return s.CreateComponent(ctx, "acme", ReleaseComponentCreate{Name: "api", Artifact: ReleaseArtifactRef{ID: releaseTestID, Name: "api"}})
		}},
		{name: "create component instance missing target", call: func(s *ReleasesService) (any, *Response, error) {
			return s.CreateComponentInstance(ctx, "acme", ReleaseTargetRef{}, ReleaseComponentInstanceCreate{})
		}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			server, client, teardown := newMockServerAndClient(t)
			t.Cleanup(teardown)
			_, resp, err := tt.call(client.Releases)
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
		call func(*ReleasesService) (any, *Response, error)
	}{
		{name: "missing selector", call: func(s *ReleasesService) (any, *Response, error) {
			return s.GetComponentInstance(ctx, "acme", ReleaseComponentInstanceRef{})
		}},
		{name: "conflicting selectors", call: func(s *ReleasesService) (any, *Response, error) {
			return s.GetComponentInstance(ctx, "acme", ReleaseComponentInstanceRef{ID: releaseTestID, Name: "api"})
		}},
		{name: "missing status name", call: func(s *ReleasesService) (any, *Response, error) {
			return s.GetComponentInstanceStatus(ctx, "acme", ReleaseTargetRef{Name: "east"}, "")
		}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			server, client, teardown := newMockServerAndClient(t)
			t.Cleanup(teardown)
			_, resp, err := tt.call(client.Releases)
			if err == nil || !strings.Contains(err.Error(), "component instance") {
				t.Fatalf("error = %v, want component instance context", err)
			}
			if resp != nil || len(server.calls) != 0 {
				t.Fatalf("response=%#v, HTTP calls=%v; invalid selectors must not send a request", resp, server.calls)
			}
		})
	}
}

func TestReleasesService_GetComponentInstanceStatus(t *testing.T) {
	t.Parallel()
	present, sequence, reason := true, int64(1), "no_active_version"
	for _, tt := range []struct {
		name string
		body string
		want ReleaseComponentInstanceStatus
	}{
		{name: "absent", body: releaseStatusJSON, want: releaseStatus},
		{
			name: "never observed",
			body: `{"target_id":"target-id","instance":"api-blue","instance_id":"instance-id","present":null,"active_digest":null,"last_active_digest":null,"first_seen_at":null,"last_seen_at":null,"sequence":null,"catalog_resolution":{"status":"unknown","reason":"no_active_version"}}`,
			want: ReleaseComponentInstanceStatus{
				TargetID: "target-id", Instance: "api-blue", InstanceID: "instance-id",
				CatalogResolution: ReleaseCatalogResolution{Status: "unknown", Reason: &reason},
			},
		},
		{
			name: "resolved present",
			body: `{"target_id":"target-id","instance":"api-blue","instance_id":"instance-id","present":true,"active_digest":"` + releaseTestDigest + `","last_active_digest":"` + releaseTestDigest + `","first_seen_at":null,"last_seen_at":null,"sequence":1,"catalog_resolution":{"status":"resolved","reason":null}}`,
			want: ReleaseComponentInstanceStatus{
				TargetID: "target-id", Instance: "api-blue", InstanceID: "instance-id", Present: &present,
				ActiveDigest: &releaseDigest, LastActiveDigest: &releaseDigest, Sequence: &sequence,
				CatalogResolution: ReleaseCatalogResolution{Status: "resolved"},
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server, client, teardown := newMockServerAndClient(t)
			t.Cleanup(teardown)
			server.HandleFunc("/v2/organizations/acme/targets/by-id/"+releaseTestID+"/instance-statuses/api-blue", func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, http.MethodGet)
				_, _ = fmt.Fprint(w, tt.body)
			})
			got, _, err := client.Releases.GetComponentInstanceStatus(context.Background(), "acme", ReleaseTargetRef{ID: releaseTestID}, "api-blue")
			if err != nil {
				t.Fatal(err)
			}
			assertReleaseResponse(t, got, tt.want, tt.body)
		})
	}
}

func TestReleasesService_GetVersion(t *testing.T) {
	t.Parallel()
	source, reason, email := "opaque source", "withdrawn", "person@example.com"
	for _, tt := range []struct {
		name string
		body string
		want ReleaseArtifactVersion
	}{
		{name: "published", body: releaseVersionJSON, want: releaseVersion},
		{
			name: "yanked provenance",
			body: `{"id":"digest-id","digest":"` + releaseTestDigest + `","source_url":"opaque source","source_build":{"id":"build-id"},"pull_requests":[],"created_by":{"id":"user-id","graphql_id":"opaque-user","name":null,"email":"person@example.com","avatar_url":null,"created_at":"2026-09-20T01:02:03Z"},"url":"https://api.example/digest","id_url":"https://api.example/digest-id","created_at":"2026-09-20T01:02:03Z","yanked":true,"yanked_at":"2026-09-21T01:02:03Z","yank_reason":"withdrawn"}`,
			want: ReleaseArtifactVersion{
				ID: "digest-id", Digest: &releaseDigest, SourceURL: &source, SourceBuild: &ReleaseReference{ID: "build-id"},
				PullRequests: []ReleasePullRequest{},
				CreatedBy:    &ReleaseCreator{ID: "user-id", GraphQLID: "opaque-user", Email: &email, CreatedAt: releaseCreatedAt},
				URL:          "https://api.example/digest", IDURL: "https://api.example/digest-id", CreatedAt: releaseCreatedAt,
				Yanked: true, YankedAt: releaseLaterAt, YankReason: &reason,
			},
		},
		{
			name: "legacy digestless publication",
			body: `{"id":"legacy-id","digest":null,"source_url":null,"source_build":null,"pull_requests":[],"created_by":null,"url":"https://api.example/legacy","id_url":"https://api.example/legacy","created_at":"2026-09-20T01:02:03Z","yanked":false,"yanked_at":null,"yank_reason":null}`,
			want: ReleaseArtifactVersion{
				ID: "legacy-id", PullRequests: []ReleasePullRequest{},
				URL: "https://api.example/legacy", IDURL: "https://api.example/legacy", CreatedAt: releaseCreatedAt,
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server, client, teardown := newMockServerAndClient(t)
			t.Cleanup(teardown)
			server.HandleFunc("/v2/organizations/acme/artifact-versions/by-id/"+releaseTestID, func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, http.MethodGet)
				_, _ = fmt.Fprint(w, tt.body)
			})
			got, _, err := client.Releases.GetVersion(context.Background(), "acme", ReleaseVersionRef{ID: releaseTestID})
			if err != nil {
				t.Fatal(err)
			}
			assertReleaseResponse(t, got, tt.want, tt.body)
		})
	}
}

func TestReleasesService_ListObservations_PayloadPresence(t *testing.T) {
	t.Parallel()
	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)
	body := releaseTestPage(`{"id":"one","target_id":"target-id","sequence":1,"event_id":"` + releaseTestID + `","type":"ComponentInstanceAppeared","instance":"api-blue","occurred_at":"2026-09-20T01:02:03Z","recorded_at":"2026-09-20T01:03:04Z"},{"id":"two","target_id":"target-id","sequence":2,"event_id":"019a1dd8-a49b-70c5-99b8-b8f9f0ed6825","type":"ComponentInstanceVersionObserved","instance":"api-blue","occurred_at":"2026-09-20T01:02:03Z","recorded_at":"2026-09-20T01:03:04Z","catalog_resolution":{"status":"unknown","reason":"no_active_version"}}`)
	server.HandleFunc("/v2/organizations/acme/targets/east/observations", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		_, _ = fmt.Fprint(w, body)
	})
	got, _, err := client.Releases.ListObservations(context.Background(), "acme", ReleaseTargetRef{Name: "east"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	reason := "no_active_version"
	want := ReleasesPage[ReleaseObservation]{
		Items: []ReleaseObservation{
			{ID: "one", TargetID: "target-id", Sequence: 1, EventID: releaseTestID, Type: "ComponentInstanceAppeared", Instance: "api-blue", OccurredAt: releaseCreatedAt, RecordedAt: releaseRecordedAt},
			{ID: "two", TargetID: "target-id", Sequence: 2, EventID: "019a1dd8-a49b-70c5-99b8-b8f9f0ed6825", Type: "ComponentInstanceVersionObserved", Instance: "api-blue", OccurredAt: releaseCreatedAt, RecordedAt: releaseRecordedAt, CatalogResolution: &ReleaseCatalogResolution{Status: "unknown", Reason: &reason}},
		},
		Links: releaseLinks,
	}
	assertReleaseResponse(t, got, want, body)
}

func TestReleasesService_APIErrorsAndContext(t *testing.T) {
	t.Parallel()
	for _, status := range []int{401, 403, 404, 422} {
		t.Run(fmt.Sprint(status), func(t *testing.T) {
			server, client, teardown := newMockServerAndClient(t)
			t.Cleanup(teardown)
			server.HandleFunc("/v2/organizations/acme/components", func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(status)
				_, _ = fmt.Fprint(w, `{"message":"not available"}`)
			})
			_, resp, err := client.Releases.ListComponents(context.Background(), "acme", nil)
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
		if r.Header.Get("Authorization") != "Bearer release-test-token" {
			t.Error("missing authentication")
		}
		_, _ = fmt.Fprint(w, `{}`)
	}))
	t.Cleanup(server.Close)
	_, _, err := must(NewClient(WithBaseURL(server.URL+"/prefix/"), WithTokenAuth("release-test-token"))).Releases.GetTarget(context.Background(), "acme/space+name", ReleaseTargetRef{Name: "prod/blue?#% +"})
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
