package buildkite

import (
	"context"
	"fmt"
	"time"
)

// RepositoryConnectionsService handles communication with the repository
// connections API.
//
// Buildkite API docs: https://buildkite.com/docs/apis/rest-api/organizations/repository-connections
type RepositoryConnectionsService struct {
	client *Client
}

// RepositoryConnection represents a source control connection available to an
// organization. ServiceAccount, Host, and RateLimit are only included in Get
// responses.
type RepositoryConnection struct {
	ID             string                              `json:"id,omitempty"`
	Type           string                              `json:"type,omitempty"`
	DisplayName    string                              `json:"display_name,omitempty"`
	URL            string                              `json:"url,omitempty"`
	ServiceAccount *RepositoryConnectionServiceAccount `json:"service_account,omitempty"`
	Host           *RepositoryConnectionHost           `json:"host,omitempty"`
	RateLimit      *RepositoryConnectionRateLimit      `json:"rate_limit,omitempty"`
}

// RepositoryConnectionServiceAccount identifies the account connected to the
// source control provider.
type RepositoryConnectionServiceAccount struct {
	Login string `json:"login,omitempty"`
}

// RepositoryConnectionHost describes the source control host. Optional fields
// vary by provider.
type RepositoryConnectionHost struct {
	Type                    string  `json:"type,omitempty"`
	URL                     string  `json:"url,omitempty"`
	WebhookAllowedAddresses *string `json:"webhook_allowed_addresses,omitempty"`
	VerifyServerCert        *bool   `json:"verify_server_cert,omitempty"`
	AuthenticationStrategy  *string `json:"authentication_strategy,omitempty"`
	ProxyAuthStrategy       *string `json:"proxy_auth_strategy,omitempty"`
}

// RepositoryConnectionRateLimit describes a GitHub installation's cached API
// quota. It is unrelated to the Buildkite API quota returned by RateLimit.Get.
type RepositoryConnectionRateLimit struct {
	Limit     int64     `json:"limit"`
	Used      int64     `json:"used"`
	Remaining int64     `json:"remaining"`
	ResetAt   time.Time `json:"reset_at"`
}

// List returns the unpaginated repository connections for an organization.
func (rcs *RepositoryConnectionsService) List(ctx context.Context, org string) ([]RepositoryConnection, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/repository_connections", org)
	req, err := rcs.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var connections []RepositoryConnection
	resp, err := rcs.client.Do(req, &connections)
	if err != nil {
		return nil, resp, err
	}

	return connections, resp, nil
}

// Get returns a repository connection by ID.
func (rcs *RepositoryConnectionsService) Get(ctx context.Context, org, id string) (RepositoryConnection, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/repository_connections/%s", org, id)
	req, err := rcs.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return RepositoryConnection{}, nil, err
	}

	var connection RepositoryConnection
	resp, err := rcs.client.Do(req, &connection)
	if err != nil {
		return RepositoryConnection{}, resp, err
	}

	return connection, resp, nil
}
