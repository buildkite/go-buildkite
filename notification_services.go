package buildkite

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// NotificationServicesService handles communication with the organization
// notification service methods of the Buildkite API.
//
// Listing and getting require the read_notification_services scope. Other
// methods require write_notification_services. The token owner must be an
// organization administrator unless the organization enables Manage
// Notification Services for members.
//
// Buildkite API docs: https://buildkite.com/docs/apis/rest-api/organizations/notification-services
type NotificationServicesService struct {
	client *Client
}

// NotificationServiceSettings is the provider-specific settings JSON object of
// a notification service. Its keys depend on the provider. For example, a
// webhook accepts url, token, token_mode, version, events, tls_verify, and
// cluster_queue_ids.
type NotificationServiceSettings map[string]any

// NotificationServiceProvider identifies the type of a notification service.
type NotificationServiceProvider struct {
	// ID is the machine-readable provider, such as webhook, slack,
	// aws_event_bridge, or open_telemetry_tracing.
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// NotificationService represents an organization notification service.
type NotificationService struct {
	ID                  string                      `json:"id,omitempty"`
	GraphQLID           *string                     `json:"graphql_id,omitempty"`
	URL                 string                      `json:"url,omitempty"`
	Provider            NotificationServiceProvider `json:"provider"`
	Description         string                      `json:"description,omitempty"`
	Enabled             bool                        `json:"enabled"`
	Scope               string                      `json:"scope,omitempty"`
	ScopeUUIDs          []string                    `json:"scope_uuids,omitempty"`
	BranchConfiguration string                      `json:"branch_configuration,omitempty"`
	BuildStates         map[string]bool             `json:"build_states,omitempty"`
	Settings            NotificationServiceSettings `json:"settings,omitempty"`
	CreatedAt           *Timestamp                  `json:"created_at,omitempty"`
	CreatedBy           *User                       `json:"created_by,omitempty"`
}

// NotificationServiceCreate represents the request body for creating a
// notification service. Provider is required. Scope defaults to all when
// omitted.
type NotificationServiceCreate struct {
	Provider            string                      `json:"provider"`
	Description         string                      `json:"description,omitempty"`
	BranchConfiguration string                      `json:"branch_configuration,omitempty"`
	Scope               string                      `json:"scope,omitempty"`
	ScopeUUIDs          []string                    `json:"scope_uuids,omitempty"`
	BuildStates         map[string]bool             `json:"build_states,omitempty"`
	Settings            NotificationServiceSettings `json:"settings,omitempty"`
}

// NotificationServiceUpdate represents the request body for updating a
// notification service. Unset fields are omitted and keep their current
// values. If set, Provider must match the existing provider. Omitted settings
// keys also remain unchanged.
type NotificationServiceUpdate struct {
	Provider            Optional[string]                      `json:"provider,omitzero"`
	Description         Optional[string]                      `json:"description,omitzero"`
	BranchConfiguration Optional[string]                      `json:"branch_configuration,omitzero"`
	Scope               Optional[string]                      `json:"scope,omitzero"`
	ScopeUUIDs          Optional[[]string]                    `json:"scope_uuids,omitzero"`
	BuildStates         Optional[map[string]bool]             `json:"build_states,omitzero"`
	Settings            Optional[NotificationServiceSettings] `json:"settings,omitzero"`
}

// NotificationServicesListOptions specifies cursor pagination parameters.
type NotificationServicesListOptions struct {
	After   string `url:"after,omitempty"`
	Before  string `url:"before,omitempty"`
	PerPage int    `url:"per_page,omitempty"`
}

type NotificationServicesListLink string

// ToOptions extracts cursor pagination options from a notification services link.
func (l NotificationServicesListLink) ToOptions() (*NotificationServicesListOptions, error) {
	u, err := url.Parse(string(l))
	if err != nil {
		return nil, fmt.Errorf("parsing link: %w", err)
	}

	q := u.Query()
	opt := &NotificationServicesListOptions{
		After:  q.Get("after"),
		Before: q.Get("before"),
	}

	if perPage := q.Get("per_page"); perPage != "" {
		opt.PerPage, err = strconv.Atoi(perPage)
		if err != nil {
			return nil, fmt.Errorf("parsing per_page: %w", err)
		}
	}

	return opt, nil
}

type NotificationServicesListLinks struct {
	Self     NotificationServicesListLink `json:"self,omitempty"`
	First    NotificationServicesListLink `json:"first,omitempty"`
	Previous NotificationServicesListLink `json:"prev,omitempty"`
	Next     NotificationServicesListLink `json:"next,omitempty"`
}

type NotificationServicesList struct {
	Items []NotificationService         `json:"items"`
	Links NotificationServicesListLinks `json:"links"`
}

// List returns the notification services in an organization, ordered from
// oldest to newest.
func (s *NotificationServicesService) List(ctx context.Context, org string, opt *NotificationServicesListOptions) (NotificationServicesList, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/services", org)
	u, err := addOptions(u, opt)
	if err != nil {
		return NotificationServicesList{}, nil, err
	}

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return NotificationServicesList{}, nil, err
	}

	var services NotificationServicesList
	resp, err := s.client.Do(req, &services)
	if err != nil {
		return NotificationServicesList{}, resp, err
	}

	return services, resp, nil
}

// Get returns a notification service by UUID.
func (s *NotificationServicesService) Get(ctx context.Context, org, serviceUUID string) (NotificationService, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/services/%s", org, serviceUUID)
	return s.do(ctx, "GET", u, nil)
}

// Create creates a notification service in an organization.
func (s *NotificationServicesService) Create(ctx context.Context, org string, input NotificationServiceCreate) (NotificationService, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/services", org)
	return s.do(ctx, "POST", u, input)
}

// Update updates a notification service by UUID.
func (s *NotificationServicesService) Update(ctx context.Context, org, serviceUUID string, input NotificationServiceUpdate) (NotificationService, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/services/%s", org, serviceUUID)
	return s.do(ctx, "PATCH", u, input)
}

// Delete deletes a notification service by UUID.
func (s *NotificationServicesService) Delete(ctx context.Context, org, serviceUUID string) (*Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/services/%s", org, serviceUUID)
	req, err := s.client.NewRequest(ctx, "DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// Enable enables a notification service and clears any previous disabled or
// broken state. It returns the updated notification service.
func (s *NotificationServicesService) Enable(ctx context.Context, org, serviceUUID string) (NotificationService, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/services/%s/enable", org, serviceUUID)
	return s.do(ctx, "PUT", u, nil)
}

// Disable disables a notification service. It returns the updated
// notification service.
func (s *NotificationServicesService) Disable(ctx context.Context, org, serviceUUID string) (NotificationService, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/services/%s/disable", org, serviceUUID)
	return s.do(ctx, "PUT", u, nil)
}

func (s *NotificationServicesService) do(ctx context.Context, method, u string, body any) (NotificationService, *Response, error) {
	req, err := s.client.NewRequest(ctx, method, u, body)
	if err != nil {
		return NotificationService{}, nil, err
	}

	var service NotificationService
	resp, err := s.client.Do(req, &service)
	if err != nil {
		return NotificationService{}, resp, err
	}

	return service, resp, nil
}
