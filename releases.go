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
	Name     string              `json:"name"`
	Artifact *ReleaseArtifactRef `json:"artifact,omitempty"`
}

type ReleaseArtifactCreate struct {
	Name          string            `json:"name"`
	VersionScheme string            `json:"version_scheme,omitempty"`
	Pipeline      *ReleaseReference `json:"pipeline,omitempty"`
}

type ReleaseVersionPublish struct {
	Version string `json:"version"`
	Digest  string `json:"digest"`
	// An unset SourceBuild uses the agent OIDC token's build, if any.
	// Null[ReleaseReference]() explicitly records no source build.
	SourceBuild  Optional[ReleaseReference] `json:"source_build,omitzero"`
	SourceURL    Optional[string]           `json:"source_url,omitzero"`
	PullRequests []ReleasePullRequest       `json:"pull_requests,omitempty"`
}

type ReleaseVersionYank struct {
	Reason string `json:"reason,omitempty"`
}

// ReleaseObservationCreate reports a fact, not desired state. EventID must be a
// caller-supplied UUIDv7, reused with the same payload if the outcome is unknown.
// Type is TargetAppeared, ComponentAppeared, ComponentDisappeared, or
// ComponentVersionObserved. Only the latter accepts Version and Digest;
// Component is required for all but TargetAppeared.
type ReleaseObservationCreate struct {
	EventID    string     `json:"event_id"`
	Type       string     `json:"type"`
	Component  string     `json:"component,omitempty"`
	Version    string     `json:"version,omitempty"`
	Digest     string     `json:"digest,omitempty"`
	OccurredAt *Timestamp `json:"occurred_at,omitempty"`
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

// ReleaseComponentRef selects by ID or by target and component name.
type ReleaseComponentRef struct {
	ID     string           `json:"id,omitempty"`
	Target ReleaseTargetRef `json:"target,omitempty,omitzero"`
	Name   string           `json:"name,omitempty"`
}

// ReleaseArtifactRef selects an artifact by name or ID, never both.
type ReleaseArtifactRef struct {
	Name string `json:"name,omitempty"`
	ID   string `json:"id,omitempty"`
}

// ReleaseVersionRef selects by ID or by artifact and exact version label.
type ReleaseVersionRef struct {
	ID       string             `json:"id,omitempty"`
	Artifact ReleaseArtifactRef `json:"artifact,omitempty,omitzero"`
	Version  string             `json:"version,omitempty"`
}

// ReleasesListOptions specifies cursor pagination, not page-number pagination.
type ReleasesListOptions struct {
	After   string `url:"after,omitempty"`
	Before  string `url:"before,omitempty"`
	PerPage int    `url:"per_page,omitempty"`
}

type ReleasesObservationListOptions struct {
	ReleasesListOptions
	Component string `url:"component,omitempty"`
	Type      string `url:"type,omitempty"`
}

type ReleasesVersionListOptions struct {
	ReleasesListOptions
	Order     string `url:"order,omitempty"`
	Direction string `url:"direction,omitempty"`
	NewerThan string `url:"newer_than,omitempty"`
	OlderThan string `url:"older_than,omitempty"`
	// A non-nil empty suffix selects unsuffixed versions; nil selects all.
	Suffix *string `url:"suffix,omitempty"`
	Yanked *bool   `url:"yanked,omitempty"`
}

// ReleasesComparisonOptions selects left relative to right. Supply exactly one
// label or ID for each side. The API validates comparison selectors and schemes.
type ReleasesComparisonOptions struct {
	Left    string `url:"left,omitempty"`
	LeftID  string `url:"left_id,omitempty"`
	Right   string `url:"right,omitempty"`
	RightID string `url:"right_id,omitempty"`
}

type ReleasesComparisonListOptions struct {
	ReleasesListOptions
	ReleasesComparisonOptions
}

type ReleasesPage[T any] struct {
	Items    []T                     `json:"items"`
	Links    ReleasePageLinks        `json:"links"`
	Ordering *ReleaseVersionOrdering `json:"ordering,omitempty"`
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

type ReleaseVersionOrdering struct {
	Basis     string `json:"basis"`
	Direction string `json:"direction"`
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
	ID         string                  `json:"id"`
	GraphQLID  string                  `json:"graphql_id,omitempty"`
	Name       string                  `json:"name"`
	TargetID   string                  `json:"target_id"`
	Artifact   *ReleaseArtifactSummary `json:"artifact"`
	URL        string                  `json:"url"`
	IDURL      string                  `json:"id_url"`
	CreatedAt  *Timestamp              `json:"created_at"`
	ArchivedAt *Timestamp              `json:"archived_at"`
}

type ReleaseArtifactSummary struct {
	ID            string `json:"id"`
	GraphQLID     string `json:"graphql_id,omitempty"`
	Name          string `json:"name"`
	VersionScheme string `json:"version_scheme"`
}

type ReleaseArtifact struct {
	ID            string            `json:"id"`
	GraphQLID     string            `json:"graphql_id,omitempty"`
	Name          string            `json:"name"`
	VersionScheme string            `json:"version_scheme"`
	Pipeline      *ReleaseReference `json:"pipeline,omitempty"`
	URL           string            `json:"url"`
	IDURL         string            `json:"id_url"`
	CreatedAt     *Timestamp        `json:"created_at"`
	ArchivedAt    *Timestamp        `json:"archived_at"`
}

// ReleaseReference is an organization-scoped build or pipeline reference.
type ReleaseReference struct {
	ID string `json:"id"`
}

type ReleaseArtifactVersion struct {
	ID           string               `json:"id"`
	GraphQLID    string               `json:"graphql_id,omitempty"`
	Version      string               `json:"version"`
	Digest       *string              `json:"digest"`
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

type ReleaseComponentStatus struct {
	TargetID          string                   `json:"target_id"`
	Component         string                   `json:"component"`
	ComponentID       string                   `json:"component_id"`
	Present           *bool                    `json:"present"`
	ActiveVersion     *string                  `json:"active_version"`
	ActiveDigest      *string                  `json:"active_digest"`
	LastActiveVersion *string                  `json:"last_active_version"`
	FirstSeenAt       *Timestamp               `json:"first_seen_at"`
	LastSeenAt        *Timestamp               `json:"last_seen_at"`
	Sequence          *int64                   `json:"sequence"`
	CatalogResolution ReleaseCatalogResolution `json:"catalog_resolution"`
}

type ReleaseCatalogResolution struct {
	Status            string  `json:"status"`
	ArtifactVersionID *string `json:"artifact_version_id"`
	DigestCheck       string  `json:"digest_check"`
	Reason            *string `json:"reason"`
}

type ReleaseVersionSummary struct {
	ID        string  `json:"id"`
	GraphQLID string  `json:"graphql_id,omitempty"`
	Version   string  `json:"version"`
	Digest    *string `json:"digest"`
	Yanked    bool    `json:"yanked"`
}

// ReleaseVersionComparison compares catalog precedence, not deployment safety.
type ReleaseVersionComparison struct {
	Artifact    ReleaseArtifactSummary `json:"artifact"`
	Left        ReleaseVersionSummary  `json:"left"`
	Right       ReleaseVersionSummary  `json:"right"`
	Basis       *string                `json:"basis"`
	Result      string                 `json:"result"`
	Reason      string                 `json:"reason"`
	SameVersion bool                   `json:"same_version"`
	SameContent *bool                  `json:"same_content"`
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
	Component         *string                   `json:"component,omitempty"`
	Version           *string                   `json:"version,omitempty"`
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
	if r.ID != "" {
		if r.Name != "" || r.Target != (ReleaseTargetRef{}) {
			return "", fmt.Errorf("component ID cannot be combined with target or name")
		}
		return "components/by-id/" + releaseSegment(r.ID), nil
	}
	if r.Name == "" {
		return "", fmt.Errorf("component requires an ID or target and name")
	}
	parent, err := r.Target.path()
	if err != nil {
		return "", err
	}
	return parent + "/components/" + releaseSegment(r.Name), nil
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

func (r ReleaseVersionRef) path() (string, error) {
	if r.ID != "" {
		if r.Version != "" || r.Artifact != (ReleaseArtifactRef{}) {
			return "", fmt.Errorf("version ID cannot be combined with artifact or version label")
		}
		return "artifact-versions/by-id/" + releaseSegment(r.ID), nil
	}
	if r.Version == "" {
		return "", fmt.Errorf("version requires an ID or artifact and version label")
	}
	parent, err := r.Artifact.path()
	if err != nil {
		return "", err
	}
	return parent + "/versions/" + releaseSegment(r.Version), nil
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

func (s *ReleasesService) ListComponents(ctx context.Context, org string, ref ReleaseTargetRef, opt *ReleasesListOptions) (ReleasesPage[ReleaseComponent], *Response, error) {
	path, err := ref.path()
	return releasesRead[ReleasesPage[ReleaseComponent]](ctx, s, org, path+"/components", err, opt)
}

func (s *ReleasesService) GetComponent(ctx context.Context, org string, ref ReleaseComponentRef) (ReleaseComponent, *Response, error) {
	path, err := ref.path()
	return releasesRead[ReleaseComponent](ctx, s, org, path, err, nil)
}

func (s *ReleasesService) ListComponentStatuses(ctx context.Context, org string, ref ReleaseTargetRef, opt *ReleasesListOptions) (ReleasesPage[ReleaseComponentStatus], *Response, error) {
	path, err := ref.path()
	return releasesRead[ReleasesPage[ReleaseComponentStatus]](ctx, s, org, path+"/component-statuses", err, opt)
}

func (s *ReleasesService) GetComponentStatus(ctx context.Context, org string, ref ReleaseTargetRef, component string) (ReleaseComponentStatus, *Response, error) {
	path, err := ref.path()
	if err == nil && component == "" {
		err = fmt.Errorf("component status requires a component name")
	}
	return releasesRead[ReleaseComponentStatus](ctx, s, org, path+"/component-statuses/"+releaseSegment(component), err, nil)
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

func (s *ReleasesService) ListVersions(ctx context.Context, org string, ref ReleaseArtifactRef, opt *ReleasesVersionListOptions) (ReleasesPage[ReleaseArtifactVersion], *Response, error) {
	path, err := ref.path()
	return releasesRead[ReleasesPage[ReleaseArtifactVersion]](ctx, s, org, path+"/versions", err, opt)
}

func (s *ReleasesService) GetVersion(ctx context.Context, org string, ref ReleaseVersionRef) (ReleaseArtifactVersion, *Response, error) {
	path, err := ref.path()
	return releasesRead[ReleaseArtifactVersion](ctx, s, org, path, err, nil)
}

func (s *ReleasesService) CompareVersions(ctx context.Context, org string, ref ReleaseArtifactRef, opt *ReleasesComparisonOptions) (ReleaseVersionComparison, *Response, error) {
	path, err := ref.path()
	return releasesRead[ReleaseVersionComparison](ctx, s, org, path+"/version-comparison", err, opt)
}

func (s *ReleasesService) ListComparisonPullRequests(ctx context.Context, org string, ref ReleaseArtifactRef, opt *ReleasesComparisonListOptions) (ReleasesPage[ReleasePullRequest], *Response, error) {
	path, err := ref.path()
	return releasesRead[ReleasesPage[ReleasePullRequest]](ctx, s, org, path+"/version-comparison/pull-requests", err, opt)
}

func (s *ReleasesService) CreateTarget(ctx context.Context, org string, input ReleaseTargetCreate) (ReleaseTarget, *Response, error) {
	return releasesRequest[ReleaseTarget](ctx, s, http.MethodPost, org, "targets", nil, input)
}

// UpdateTarget replaces labels when supplied; it does not rename the target.
func (s *ReleasesService) UpdateTarget(ctx context.Context, org string, ref ReleaseTargetRef, input ReleaseTargetUpdate) (ReleaseTarget, *Response, error) {
	path, err := ref.path()
	return releasesRequest[ReleaseTarget](ctx, s, http.MethodPatch, org, path, err, input)
}

func (s *ReleasesService) CreateComponent(ctx context.Context, org string, target ReleaseTargetRef, input ReleaseComponentCreate) (ReleaseComponent, *Response, error) {
	path, err := target.path()
	if err == nil && input.Artifact != nil {
		_, err = input.Artifact.path()
	}
	return releasesRequest[ReleaseComponent](ctx, s, http.MethodPost, org, path+"/components", err, input)
}

func (s *ReleasesService) CreateArtifact(ctx context.Context, org string, input ReleaseArtifactCreate) (ReleaseArtifact, *Response, error) {
	return releasesRequest[ReleaseArtifact](ctx, s, http.MethodPost, org, "artifacts", nil, input)
}

// PublishVersion records an immutable version and digest. Replaying the same
// version and digest returns the original publication and provenance (HTTP 200).
// It does not upload or verify content.
func (s *ReleasesService) PublishVersion(ctx context.Context, org string, artifact ReleaseArtifactRef, input ReleaseVersionPublish) (ReleaseArtifactVersion, *Response, error) {
	path, err := artifact.path()
	return releasesRequest[ReleaseArtifactVersion](ctx, s, http.MethodPost, org, path+"/versions", err, input)
}

// YankVersion preserves history. Repeated calls retain the original yank time
// and reason; a yanked version identity cannot be republished.
func (s *ReleasesService) YankVersion(ctx context.Context, org string, ref ReleaseVersionRef, input ReleaseVersionYank) (ReleaseArtifactVersion, *Response, error) {
	path, err := ref.path()
	return releasesRequest[ReleaseArtifactVersion](ctx, s, http.MethodPost, org, path+"/yank", err, input)
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

// ArchiveTarget also archives its active components, without deleting history.
func (s *ReleasesService) ArchiveTarget(ctx context.Context, org string, ref ReleaseTargetRef) (ReleaseTarget, *Response, error) {
	path, err := ref.path()
	return releasesRequest[ReleaseTarget](ctx, s, http.MethodPost, org, path+"/archive", err, struct{}{})
}

// UnarchiveTarget restores only the target, not its components.
func (s *ReleasesService) UnarchiveTarget(ctx context.Context, org string, ref ReleaseTargetRef) (ReleaseTarget, *Response, error) {
	path, err := ref.path()
	return releasesRequest[ReleaseTarget](ctx, s, http.MethodPost, org, path+"/unarchive", err, struct{}{})
}

func (s *ReleasesService) ArchiveComponent(ctx context.Context, org string, ref ReleaseComponentRef) (ReleaseComponent, *Response, error) {
	path, err := ref.path()
	return releasesRequest[ReleaseComponent](ctx, s, http.MethodPost, org, path+"/archive", err, struct{}{})
}

// UnarchiveComponent requires an active target.
func (s *ReleasesService) UnarchiveComponent(ctx context.Context, org string, ref ReleaseComponentRef) (ReleaseComponent, *Response, error) {
	path, err := ref.path()
	return releasesRequest[ReleaseComponent](ctx, s, http.MethodPost, org, path+"/unarchive", err, struct{}{})
}

func (s *ReleasesService) ArchiveArtifact(ctx context.Context, org string, ref ReleaseArtifactRef) (ReleaseArtifact, *Response, error) {
	path, err := ref.path()
	return releasesRequest[ReleaseArtifact](ctx, s, http.MethodPost, org, path+"/archive", err, struct{}{})
}

func (s *ReleasesService) UnarchiveArtifact(ctx context.Context, org string, ref ReleaseArtifactRef) (ReleaseArtifact, *Response, error) {
	path, err := ref.path()
	return releasesRequest[ReleaseArtifact](ctx, s, http.MethodPost, org, path+"/unarchive", err, struct{}{})
}
