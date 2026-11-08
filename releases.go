package buildkite

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// ReleasesService reads infrastructure registrations, observations and the
// release artifact catalog. These artifacts are distinct from job artifacts.
type ReleasesService struct {
	client *Client
}

// ReleaseInstallationRef selects by name or ID, never both. UUID-shaped names
// remain names; IDs must be selected explicitly.
type ReleaseInstallationRef struct {
	Name string `json:"name,omitempty"`
	ID   string `json:"id,omitempty"`
}

// ReleaseTargetRef selects a target by ID or by installation and target name.
type ReleaseTargetRef struct {
	ID           string                 `json:"id,omitempty"`
	Installation ReleaseInstallationRef `json:"installation,omitempty,omitzero"`
	Name         string                 `json:"name,omitempty"`
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

type ReleasePageLinks struct {
	Self     string `json:"self"`
	First    string `json:"first,omitempty"`
	Previous string `json:"prev,omitempty"`
	Next     string `json:"next,omitempty"`
}

type ReleaseVersionOrdering struct {
	Basis     string `json:"basis"`
	Direction string `json:"direction"`
}

type ReleaseInstallation struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	URL        string     `json:"url"`
	IDURL      string     `json:"id_url"`
	CreatedAt  *Timestamp `json:"created_at"`
	ArchivedAt *Timestamp `json:"archived_at"`
}

type ReleaseTarget struct {
	ID         string            `json:"id"`
	Name       string            `json:"name"`
	Labels     map[string]string `json:"labels"`
	URL        string            `json:"url"`
	IDURL      string            `json:"id_url"`
	CreatedAt  *Timestamp        `json:"created_at"`
	ArchivedAt *Timestamp        `json:"archived_at"`
}

type ReleaseComponent struct {
	ID         string                  `json:"id"`
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
	Name          string `json:"name"`
	VersionScheme string `json:"version_scheme"`
}

type ReleaseArtifact struct {
	ID            string            `json:"id"`
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
	ID      string  `json:"id"`
	Version string  `json:"version"`
	Digest  *string `json:"digest"`
	Yanked  bool    `json:"yanked"`
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

func (r ReleaseInstallationRef) path() (string, error) {
	if (r.ID == "") == (r.Name == "") {
		return "", fmt.Errorf("installation requires exactly one of name or ID")
	}
	if r.ID != "" {
		return "installations/by-id/" + releaseSegment(r.ID), nil
	}
	return "installations/" + releaseSegment(r.Name), nil
}

func (r ReleaseTargetRef) path() (string, error) {
	if r.ID != "" {
		if r.Name != "" || r.Installation != (ReleaseInstallationRef{}) {
			return "", fmt.Errorf("target ID cannot be combined with installation or name")
		}
		return "targets/by-id/" + releaseSegment(r.ID), nil
	}
	if r.Name == "" {
		return "", fmt.Errorf("target requires an ID or installation and name")
	}
	parent, err := r.Installation.path()
	if err != nil {
		return "", err
	}
	return parent + "/targets/" + releaseSegment(r.Name), nil
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

// releasesRead shares request construction and response/error handling across
// the read-only endpoints. Selector errors are checked before any HTTP request.
func releasesRead[T any](ctx context.Context, s *ReleasesService, org, path string, pathErr error, opt any) (T, *Response, error) {
	var zero T
	if pathErr != nil {
		return zero, nil, pathErr
	}
	u := "v2/organizations/" + releaseSegment(org) + "/" + path
	var err error
	if opt != nil {
		u, err = addOptions(u, opt)
		if err != nil {
			return zero, nil, err
		}
	}
	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
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

// ListInstallations returns installations, including archived registrations.
func (s *ReleasesService) ListInstallations(ctx context.Context, org string, opt *ReleasesListOptions) (ReleasesPage[ReleaseInstallation], *Response, error) {
	return releasesRead[ReleasesPage[ReleaseInstallation]](ctx, s, org, "installations", nil, opt)
}

func (s *ReleasesService) GetInstallation(ctx context.Context, org string, ref ReleaseInstallationRef) (ReleaseInstallation, *Response, error) {
	path, err := ref.path()
	return releasesRead[ReleaseInstallation](ctx, s, org, path, err, nil)
}

func (s *ReleasesService) ListTargets(ctx context.Context, org string, ref ReleaseInstallationRef, opt *ReleasesListOptions) (ReleasesPage[ReleaseTarget], *Response, error) {
	path, err := ref.path()
	return releasesRead[ReleasesPage[ReleaseTarget]](ctx, s, org, path+"/targets", err, opt)
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
