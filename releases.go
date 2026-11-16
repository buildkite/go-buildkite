package buildkite

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// ReleasesService manages infrastructure registrations, observations and the
// release artifact catalog. These artifacts are distinct from job artifacts.
type ReleasesService struct {
	client *Client
}

type ReleaseTargetCreate struct {
	Name   string            `json:"name"`
	Labels map[string]string `json:"labels,omitempty"`
}

type ReleaseTargetUpdate struct {
	// Labels replaces the entire map. Some(map[string]string{}) clears it;
	// an unset value leaves labels unchanged.
	Labels Optional[map[string]string] `json:"labels,omitzero"`
}

type ReleaseComponentCreate struct {
	Name     string             `json:"name"`
	Artifact ReleaseArtifactRef `json:"artifact"`
}

// ReleaseComponentInstanceCreate registers a catalog Component on a Target.
// Name defaults to ComponentName when omitted.
type ReleaseComponentInstanceCreate struct {
	ComponentName string `json:"componentName"`
	Name          string `json:"name,omitempty"`
}

// ReleaseComponentInstanceRegister addresses the Target by UUID in the body.
type ReleaseComponentInstanceRegister struct {
	ReleaseComponentInstanceCreate
	TargetID string `json:"targetId"`
}

type ReleaseArtifactCreate struct {
	Name string `json:"name"`
}

type ReleaseDigestPublish struct {
	Digest string `json:"digest"`
	// An unset SourceBuild uses the agent OIDC token's build, if any.
	// Null[ReleaseReference]() explicitly records no source build.
	SourceBuild  Optional[ReleaseReference] `json:"source_build,omitzero"`
	SourceURL    Optional[string]           `json:"source_url,omitzero"`
	PullRequests []ReleasePullRequest       `json:"pull_requests,omitempty"`
}

type ReleaseDigestYank struct {
	Reason string `json:"reason,omitempty"`
}

// ReleaseObservationCreate reports a fact, not desired state. EventID must be a
// caller-supplied UUIDv7, reused with the same payload if the outcome is unknown.
// Type is ComponentInstanceAppeared, ComponentInstanceDisappeared, or
// ComponentInstanceDigestObserved. All require Instance. Only the latter accepts
// ObservedTag and requires Digest.
type ReleaseObservationCreate struct {
	EventID     string     `json:"event_id"`
	Type        string     `json:"type"`
	Instance    string     `json:"instance"`
	ObservedTag string     `json:"observed_tag,omitempty"`
	Digest      string     `json:"digest,omitempty"`
	OccurredAt  *Timestamp `json:"occurred_at,omitempty"`
}

type ReleaseObservationWriteResult struct {
	ReleaseObservation
	Created bool `json:"created"`
}

// ReleaseTargetRef selects an organization-owned target by name or ID, never
// both. UUID-shaped names remain names; IDs must be selected explicitly.
type ReleaseTargetRef struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// ReleaseComponentRef selects an organization-owned catalog Component by name or ID.
type ReleaseComponentRef struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// ReleaseComponentInstanceRef selects by ID or by Target and Instance name.
type ReleaseComponentInstanceRef struct {
	ID     string           `json:"id,omitempty"`
	Target ReleaseTargetRef `json:"target,omitempty,omitzero"`
	Name   string           `json:"name,omitempty"`
}

// ReleaseArtifactRef selects an artifact by name or ID, never both.
type ReleaseArtifactRef struct {
	Name string `json:"name,omitempty"`
	ID   string `json:"id,omitempty"`
}

// ReleaseDigestRef selects by ID or by Artifact and exact SHA-256 digest.
type ReleaseDigestRef struct {
	ID       string             `json:"id,omitempty"`
	Artifact ReleaseArtifactRef `json:"artifact,omitempty,omitzero"`
	Digest   string             `json:"digest,omitempty"`
}

// ReleasesListOptions specifies cursor pagination, not page-number pagination.
type ReleasesListOptions struct {
	After   string `url:"after,omitempty"`
	Before  string `url:"before,omitempty"`
	PerPage int    `url:"per_page,omitempty"`
}

type ReleasesObservationListOptions struct {
	ReleasesListOptions
	Instance string `url:"instance,omitempty"`
	Type     string `url:"type,omitempty"`
}

type ReleasesComponentInstanceListOptions struct {
	ReleasesListOptions
	// Label matches all supplied label pairs on the Instance's Target.
	Label ReleaseTargetLabelFilters `url:"label,omitempty"`
}

type ReleaseTargetLabelFilters map[string]string

func (labels ReleaseTargetLabelFilters) EncodeValues(key string, values *url.Values) error {
	for name, value := range labels {
		values.Set(key+"["+name+"]", value)
	}
	return nil
}

type ReleasesPage[T any] struct {
	Items []T              `json:"items"`
	Links ReleasePageLinks `json:"links"`
}

type ReleasePageLink string

// ToOptions extracts cursor pagination options from a Releases link.
// Retain any resource-specific filters when applying these options.
func (l ReleasePageLink) ToOptions() (*ReleasesListOptions, error) {
	u, err := url.Parse(string(l))
	if err != nil {
		return nil, fmt.Errorf("parsing link: %w", err)
	}

	opt := &ReleasesListOptions{
		After:  u.Query().Get("after"),
		Before: u.Query().Get("before"),
	}

	if perPage := u.Query().Get("per_page"); perPage != "" {
		opt.PerPage, err = strconv.Atoi(perPage)
		if err != nil {
			return nil, fmt.Errorf("parsing per_page: %w", err)
		}
	}

	return opt, nil
}

type ReleasePageLinks struct {
	Self     ReleasePageLink `json:"self"`
	First    ReleasePageLink `json:"first,omitempty"`
	Previous ReleasePageLink `json:"prev,omitempty"`
	Next     ReleasePageLink `json:"next,omitempty"`
}

type ReleaseTarget struct {
	ID         string            `json:"id"`
	GraphQLID  string            `json:"graphql_id,omitempty"`
	Name       string            `json:"name"`
	Labels     map[string]string `json:"labels"`
	URL        string            `json:"url"`
	IDURL      string            `json:"id_url"`
	CreatedAt  *Timestamp        `json:"created_at"`
	ArchivedAt *Timestamp        `json:"archived_at"`
}

type ReleaseComponent struct {
	ID         string                 `json:"id"`
	GraphQLID  string                 `json:"graphql_id,omitempty"`
	Name       string                 `json:"name"`
	Artifact   ReleaseArtifactSummary `json:"artifact"`
	URL        string                 `json:"url"`
	IDURL      string                 `json:"id_url"`
	CreatedAt  *Timestamp             `json:"created_at"`
	ArchivedAt *Timestamp             `json:"archived_at"`
}

type ReleaseComponentInstance struct {
	ID          string     `json:"id"`
	GraphQLID   string     `json:"graphql_id,omitempty"`
	Name        string     `json:"name"`
	TargetID    string     `json:"target_id"`
	ComponentID string     `json:"component_id"`
	URL         string     `json:"url"`
	IDURL       string     `json:"id_url"`
	CreatedAt   *Timestamp `json:"created_at"`
	ArchivedAt  *Timestamp `json:"archived_at"`
}

type ReleaseArtifactSummary struct {
	ID        string `json:"id"`
	GraphQLID string `json:"graphql_id,omitempty"`
	Name      string `json:"name"`
}

type ReleaseArtifact struct {
	ID         string     `json:"id"`
	GraphQLID  string     `json:"graphql_id,omitempty"`
	Name       string     `json:"name"`
	URL        string     `json:"url"`
	IDURL      string     `json:"id_url"`
	CreatedAt  *Timestamp `json:"created_at"`
	ArchivedAt *Timestamp `json:"archived_at"`
}

// ReleaseReference is an organization-scoped build reference.
type ReleaseReference struct {
	ID string `json:"id"`
}

type ReleaseArtifactDigest struct {
	ID           string               `json:"id"`
	GraphQLID    string               `json:"graphql_id,omitempty"`
	Digest       string               `json:"digest"`
	SourceURL    *string              `json:"source_url"`
	SourceBuild  *ReleaseReference    `json:"source_build"`
	PullRequests []ReleasePullRequest `json:"pull_requests"`
	CreatedBy    *ReleaseCreator      `json:"created_by"`
	URL          string               `json:"url"`
	IDURL        string               `json:"id_url"`
	CreatedAt    *Timestamp           `json:"created_at"`
	Yanked       bool                 `json:"yanked"`
	YankedAt     *Timestamp           `json:"yanked_at"`
	YankReason   *string              `json:"yank_reason"`
}

type ReleaseCreator struct {
	ID        string     `json:"id"`
	GraphQLID string     `json:"graphql_id"`
	Name      *string    `json:"name"`
	Email     *string    `json:"email"`
	AvatarURL *string    `json:"avatar_url"`
	CreatedAt *Timestamp `json:"created_at"`
}

// ReleasePullRequest is publisher-supplied, not a verified change inventory.
type ReleasePullRequest struct {
	URL           string `json:"url"`
	Provider      string `json:"provider,omitempty"`
	RepositoryURL string `json:"repository_url,omitempty"`
	ID            string `json:"id,omitempty"`
}

type ReleaseComponentInstanceStatus struct {
	TargetID          string                   `json:"target_id"`
	Instance          string                   `json:"instance"`
	InstanceID        string                   `json:"instance_id"`
	Present           *bool                    `json:"present"`
	ActiveDigest      *string                  `json:"active_digest"`
	LastActiveDigest  *string                  `json:"last_active_digest"`
	FirstSeenAt       *Timestamp               `json:"first_seen_at"`
	LastSeenAt        *Timestamp               `json:"last_seen_at"`
	Sequence          *int64                   `json:"sequence"`
	CatalogResolution ReleaseCatalogResolution `json:"catalog_resolution"`
}

type ReleaseCatalogResolution struct {
	Status           string  `json:"status"`
	ArtifactDigestID *string `json:"artifact_digest_id"`
	Reason           *string `json:"reason"`
}

// ReleaseObservation is a stored observation, without the write-only created
// flag. Optional payload fields are omitted, unlike nullable status fields.
type ReleaseObservation struct {
	ID                string                    `json:"id"`
	GraphQLID         string                    `json:"graphql_id,omitempty"`
	TargetID          string                    `json:"target_id"`
	Sequence          int64                     `json:"sequence"`
	EventID           string                    `json:"event_id"`
	Type              string                    `json:"type"`
	Instance          string                    `json:"instance"`
	ObservedTag       *string                   `json:"observed_tag,omitempty"`
	Digest            *string                   `json:"digest,omitempty"`
	OccurredAt        *Timestamp                `json:"occurred_at"`
	RecordedAt        *Timestamp                `json:"recorded_at"`
	CatalogResolution *ReleaseCatalogResolution `json:"catalog_resolution,omitempty"`
}

func releaseSegment(value string) string {
	return strings.ReplaceAll(url.PathEscape(value), "+", "%2B")
}

func (r ReleaseTargetRef) path() (string, error) {
	if (r.ID == "") == (r.Name == "") {
		return "", fmt.Errorf("target requires exactly one of name or ID")
	}
	if r.ID != "" {
		return "targets/by-id/" + releaseSegment(r.ID), nil
	}
	return "targets/" + releaseSegment(r.Name), nil
}

func (r ReleaseComponentRef) path() (string, error) {
	if (r.ID == "") == (r.Name == "") {
		return "", fmt.Errorf("component requires exactly one of name or ID")
	}
	if r.ID != "" {
		return "components/by-id/" + releaseSegment(r.ID), nil
	}
	return "components/" + releaseSegment(r.Name), nil
}

func (r ReleaseComponentInstanceRef) path() (string, error) {
	if r.ID != "" {
		if r.Name != "" || r.Target != (ReleaseTargetRef{}) {
			return "", fmt.Errorf("instance ID cannot be combined with target or name")
		}
		return "component_instance/by-id/" + releaseSegment(r.ID), nil
	}
	if r.Name == "" {
		return "", fmt.Errorf("instance requires an ID or target and name")
	}
	parent, err := r.Target.path()
	if err != nil {
		return "", err
	}
	return parent + "/instances/" + releaseSegment(r.Name), nil
}

func (r ReleaseArtifactRef) path() (string, error) {
	if (r.ID == "") == (r.Name == "") {
		return "", fmt.Errorf("artifact requires exactly one of name or ID")
	}
	if r.ID != "" {
		return "artifacts/by-id/" + releaseSegment(r.ID), nil
	}
	return "artifacts/" + releaseSegment(r.Name), nil
}

func (r ReleaseDigestRef) path() (string, error) {
	if r.ID != "" {
		if r.Digest != "" || r.Artifact != (ReleaseArtifactRef{}) {
			return "", fmt.Errorf("digest ID cannot be combined with artifact or digest")
		}
		return "artifact-digests/by-id/" + releaseSegment(r.ID), nil
	}
	if r.Digest == "" {
		return "", fmt.Errorf("digest requires an ID or artifact and digest")
	}
	parent, err := r.Artifact.path()
	if err != nil {
		return "", err
	}
	return parent + "/digests/" + releaseSegment(r.Digest), nil
}

func releasesRead[T any](ctx context.Context, s *ReleasesService, org, path string, pathErr error, opt any) (T, *Response, error) {
	if pathErr == nil && opt != nil {
		path, pathErr = addOptions(path, opt)
	}
	return releasesRequest[T](ctx, s, http.MethodGet, org, path, pathErr, nil)
}

// releasesRequest shares request construction and response/error handling.
// Selector errors are checked before any HTTP request.
func releasesRequest[T any](ctx context.Context, s *ReleasesService, method, org, path string, pathErr error, body any) (T, *Response, error) {
	var zero T
	if pathErr != nil {
		return zero, nil, pathErr
	}
	u := "v2/organizations/" + releaseSegment(org) + "/" + path
	req, err := s.client.NewRequest(ctx, method, u, body)
	if err != nil {
		return zero, nil, err
	}
	// NewRequest's legacy relative URL resolver discards RawPath. Restore
	// escaped segments here without changing other services' URL behavior.
	escapedPath, _, _ := strings.Cut(u, "?")
	req.URL.RawPath = strings.TrimSuffix(s.client.BaseURL.EscapedPath(), "/") + "/" + escapedPath
	var result T
	resp, err := s.client.Do(req, &result)
	if err != nil {
		return zero, resp, err
	}
	return result, resp, nil
}

// ListTargets returns organization-owned targets, including archived registrations.
func (s *ReleasesService) ListTargets(ctx context.Context, org string, opt *ReleasesListOptions) (ReleasesPage[ReleaseTarget], *Response, error) {
	return releasesRead[ReleasesPage[ReleaseTarget]](ctx, s, org, "targets", nil, opt)
}

func (s *ReleasesService) GetTarget(ctx context.Context, org string, ref ReleaseTargetRef) (ReleaseTarget, *Response, error) {
	path, err := ref.path()
	return releasesRead[ReleaseTarget](ctx, s, org, path, err, nil)
}

func (s *ReleasesService) ListComponents(ctx context.Context, org string, opt *ReleasesListOptions) (ReleasesPage[ReleaseComponent], *Response, error) {
	return releasesRead[ReleasesPage[ReleaseComponent]](ctx, s, org, "components", nil, opt)
}

func (s *ReleasesService) ListComponentInstances(ctx context.Context, org string, target ReleaseTargetRef, opt *ReleasesComponentInstanceListOptions) (ReleasesPage[ReleaseComponentInstance], *Response, error) {
	path, err := target.path()
	return releasesRead[ReleasesPage[ReleaseComponentInstance]](ctx, s, org, path+"/instances", err, opt)
}

func (s *ReleasesService) ListOrganizationComponentInstances(ctx context.Context, org string, opt *ReleasesComponentInstanceListOptions) (ReleasesPage[ReleaseComponentInstance], *Response, error) {
	return releasesRead[ReleasesPage[ReleaseComponentInstance]](ctx, s, org, "component_instance", nil, opt)
}

func (s *ReleasesService) GetComponentInstance(ctx context.Context, org string, ref ReleaseComponentInstanceRef) (ReleaseComponentInstance, *Response, error) {
	path, err := ref.path()
	return releasesRead[ReleaseComponentInstance](ctx, s, org, path, err, nil)
}

func (s *ReleasesService) GetComponent(ctx context.Context, org string, ref ReleaseComponentRef) (ReleaseComponent, *Response, error) {
	path, err := ref.path()
	return releasesRead[ReleaseComponent](ctx, s, org, path, err, nil)
}

func (s *ReleasesService) ListComponentInstanceStatuses(ctx context.Context, org string, ref ReleaseTargetRef, opt *ReleasesListOptions) (ReleasesPage[ReleaseComponentInstanceStatus], *Response, error) {
	path, err := ref.path()
	return releasesRead[ReleasesPage[ReleaseComponentInstanceStatus]](ctx, s, org, path+"/instance-statuses", err, opt)
}

func (s *ReleasesService) GetComponentInstanceStatus(ctx context.Context, org string, ref ReleaseTargetRef, instance string) (ReleaseComponentInstanceStatus, *Response, error) {
	path, err := ref.path()
	if err == nil && instance == "" {
		err = fmt.Errorf("instance status requires an instance name")
	}
	return releasesRead[ReleaseComponentInstanceStatus](ctx, s, org, path+"/instance-statuses/"+releaseSegment(instance), err, nil)
}

func (s *ReleasesService) ListObservations(ctx context.Context, org string, ref ReleaseTargetRef, opt *ReleasesObservationListOptions) (ReleasesPage[ReleaseObservation], *Response, error) {
	path, err := ref.path()
	return releasesRead[ReleasesPage[ReleaseObservation]](ctx, s, org, path+"/observations", err, opt)
}

func (s *ReleasesService) ListArtifacts(ctx context.Context, org string, opt *ReleasesListOptions) (ReleasesPage[ReleaseArtifact], *Response, error) {
	return releasesRead[ReleasesPage[ReleaseArtifact]](ctx, s, org, "artifacts", nil, opt)
}

func (s *ReleasesService) GetArtifact(ctx context.Context, org string, ref ReleaseArtifactRef) (ReleaseArtifact, *Response, error) {
	path, err := ref.path()
	return releasesRead[ReleaseArtifact](ctx, s, org, path, err, nil)
}

func (s *ReleasesService) ListDigests(ctx context.Context, org string, ref ReleaseArtifactRef, opt *ReleasesListOptions) (ReleasesPage[ReleaseArtifactDigest], *Response, error) {
	path, err := ref.path()
	return releasesRead[ReleasesPage[ReleaseArtifactDigest]](ctx, s, org, path+"/digests", err, opt)
}

func (s *ReleasesService) GetDigest(ctx context.Context, org string, ref ReleaseDigestRef) (ReleaseArtifactDigest, *Response, error) {
	path, err := ref.path()
	return releasesRead[ReleaseArtifactDigest](ctx, s, org, path, err, nil)
}

func (s *ReleasesService) CreateTarget(ctx context.Context, org string, input ReleaseTargetCreate) (ReleaseTarget, *Response, error) {
	return releasesRequest[ReleaseTarget](ctx, s, http.MethodPost, org, "targets", nil, input)
}

// UpdateTarget replaces labels when supplied; it does not rename the target.
func (s *ReleasesService) UpdateTarget(ctx context.Context, org string, ref ReleaseTargetRef, input ReleaseTargetUpdate) (ReleaseTarget, *Response, error) {
	path, err := ref.path()
	return releasesRequest[ReleaseTarget](ctx, s, http.MethodPatch, org, path, err, input)
}

func (s *ReleasesService) CreateComponent(ctx context.Context, org string, input ReleaseComponentCreate) (ReleaseComponent, *Response, error) {
	_, err := input.Artifact.path()
	return releasesRequest[ReleaseComponent](ctx, s, http.MethodPost, org, "components", err, input)
}

func (s *ReleasesService) CreateComponentInstance(ctx context.Context, org string, target ReleaseTargetRef, input ReleaseComponentInstanceCreate) (ReleaseComponentInstance, *Response, error) {
	path, err := target.path()
	return releasesRequest[ReleaseComponentInstance](ctx, s, http.MethodPost, org, path+"/instances", err, input)
}

func (s *ReleasesService) RegisterComponentInstance(ctx context.Context, org string, input ReleaseComponentInstanceRegister) (ReleaseComponentInstance, *Response, error) {
	return releasesRequest[ReleaseComponentInstance](ctx, s, http.MethodPost, org, "component_instance", nil, input)
}

func (s *ReleasesService) CreateArtifact(ctx context.Context, org string, input ReleaseArtifactCreate) (ReleaseArtifact, *Response, error) {
	return releasesRequest[ReleaseArtifact](ctx, s, http.MethodPost, org, "artifacts", nil, input)
}

// PublishDigest records an immutable SHA-256 digest and provenance. Replaying
// the same Artifact and digest returns the original publication (HTTP 200).
// It does not upload or verify content.
func (s *ReleasesService) PublishDigest(ctx context.Context, org string, artifact ReleaseArtifactRef, input ReleaseDigestPublish) (ReleaseArtifactDigest, *Response, error) {
	path, err := artifact.path()
	return releasesRequest[ReleaseArtifactDigest](ctx, s, http.MethodPost, org, path+"/digests", err, input)
}

// YankDigest preserves history. Repeated calls retain the original yank time
// and reason; a yanked digest identity cannot be republished.
func (s *ReleasesService) YankDigest(ctx context.Context, org string, ref ReleaseDigestRef, input ReleaseDigestYank) (ReleaseArtifactDigest, *Response, error) {
	path, err := ref.path()
	return releasesRequest[ReleaseArtifactDigest](ctx, s, http.MethodPost, org, path+"/yank", err, input)
}

// AppendObservations reports 1–100 events in input order. The API returns 201
// when all are new, or 200 when any event is a replay. Callers own event IDs and
// recovery from unknown outcomes; this method does not generate IDs or retry
// transport/server errors. The client's configured rate-limit retries apply.
func (s *ReleasesService) AppendObservations(ctx context.Context, org string, target ReleaseTargetRef, events []ReleaseObservationCreate) ([]ReleaseObservationWriteResult, *Response, error) {
	path, err := target.path()
	input := struct {
		Events []ReleaseObservationCreate `json:"events"`
	}{Events: events}
	return releasesRequest[[]ReleaseObservationWriteResult](ctx, s, http.MethodPost, org, path+"/observations", err, input)
}

// ArchiveTarget also archives its active Instances, without deleting history.
func (s *ReleasesService) ArchiveTarget(ctx context.Context, org string, ref ReleaseTargetRef) (ReleaseTarget, *Response, error) {
	path, err := ref.path()
	return releasesRequest[ReleaseTarget](ctx, s, http.MethodPost, org, path+"/archive", err, struct{}{})
}

// UnarchiveTarget restores only the Target, not its Instances.
func (s *ReleasesService) UnarchiveTarget(ctx context.Context, org string, ref ReleaseTargetRef) (ReleaseTarget, *Response, error) {
	path, err := ref.path()
	return releasesRequest[ReleaseTarget](ctx, s, http.MethodPost, org, path+"/unarchive", err, struct{}{})
}

func (s *ReleasesService) ArchiveComponent(ctx context.Context, org string, ref ReleaseComponentRef) (ReleaseComponent, *Response, error) {
	path, err := ref.path()
	return releasesRequest[ReleaseComponent](ctx, s, http.MethodPost, org, path+"/archive", err, struct{}{})
}

func (s *ReleasesService) UnarchiveComponent(ctx context.Context, org string, ref ReleaseComponentRef) (ReleaseComponent, *Response, error) {
	path, err := ref.path()
	return releasesRequest[ReleaseComponent](ctx, s, http.MethodPost, org, path+"/unarchive", err, struct{}{})
}

func (s *ReleasesService) ArchiveComponentInstance(ctx context.Context, org string, ref ReleaseComponentInstanceRef) (ReleaseComponentInstance, *Response, error) {
	path, err := ref.path()
	return releasesRequest[ReleaseComponentInstance](ctx, s, http.MethodPost, org, path+"/archive", err, struct{}{})
}

// UnarchiveComponentInstance requires an active Target.
func (s *ReleasesService) UnarchiveComponentInstance(ctx context.Context, org string, ref ReleaseComponentInstanceRef) (ReleaseComponentInstance, *Response, error) {
	path, err := ref.path()
	return releasesRequest[ReleaseComponentInstance](ctx, s, http.MethodPost, org, path+"/unarchive", err, struct{}{})
}

func (s *ReleasesService) ArchiveArtifact(ctx context.Context, org string, ref ReleaseArtifactRef) (ReleaseArtifact, *Response, error) {
	path, err := ref.path()
	return releasesRequest[ReleaseArtifact](ctx, s, http.MethodPost, org, path+"/archive", err, struct{}{})
}

func (s *ReleasesService) UnarchiveArtifact(ctx context.Context, org string, ref ReleaseArtifactRef) (ReleaseArtifact, *Response, error) {
	path, err := ref.path()
	return releasesRequest[ReleaseArtifact](ctx, s, http.MethodPost, org, path+"/unarchive", err, struct{}{})
}
