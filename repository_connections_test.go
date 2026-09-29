package buildkite

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestRepositoryConnectionsService_List(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	server.HandleFunc("/v2/organizations/acme-inc/repository_connections", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		if r.URL.RawQuery != "" {
			t.Errorf("query = %q, want empty", r.URL.RawQuery)
		}
		_, _ = fmt.Fprint(w, `[
			{
				"id": "01234567-89ab-cdef-0123-456789abcdef",
				"type": "github_code_access_app",
				"display_name": "GitHub (acme-inc)",
				"url": "https://api.buildkite.com/v2/organizations/acme-inc/repository_connections/01234567-89ab-cdef-0123-456789abcdef"
			},
			{
				"id": "12345678-9abc-def0-1234-56789abcdef0",
				"type": "gitlab_self_managed",
				"display_name": "GitLab Self-Managed",
				"url": "https://api.buildkite.com/v2/organizations/acme-inc/repository_connections/12345678-9abc-def0-1234-56789abcdef0"
			}
		]`)
	})

	got, _, err := client.RepositoryConnections.List(context.Background(), "acme-inc")
	if err != nil {
		t.Fatalf("RepositoryConnections.List returned error: %v", err)
	}

	want := []RepositoryConnection{
		{
			ID:          "01234567-89ab-cdef-0123-456789abcdef",
			Type:        "github_code_access_app",
			DisplayName: "GitHub (acme-inc)",
			URL:         "https://api.buildkite.com/v2/organizations/acme-inc/repository_connections/01234567-89ab-cdef-0123-456789abcdef",
		},
		{
			ID:          "12345678-9abc-def0-1234-56789abcdef0",
			Type:        "gitlab_self_managed",
			DisplayName: "GitLab Self-Managed",
			URL:         "https://api.buildkite.com/v2/organizations/acme-inc/repository_connections/12345678-9abc-def0-1234-56789abcdef0",
		},
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("RepositoryConnections.List diff (-want +got):\n%s", diff)
	}
}

func TestRepositoryConnectionsService_Get(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	const id = "01234567-89ab-cdef-0123-456789abcdef"
	server.HandleFunc("/v2/organizations/acme-inc/repository_connections/"+id, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = fmt.Fprint(w, `{
			"id": "01234567-89ab-cdef-0123-456789abcdef",
			"type": "github_enterprise_app",
			"display_name": "GitHub Enterprise Server (acme-ghe)",
			"url": "https://api.buildkite.com/v2/organizations/acme-inc/repository_connections/01234567-89ab-cdef-0123-456789abcdef",
			"service_account": {"login": "acme-ghe"},
			"host": {
				"type": "github_enterprise_app",
				"url": "https://ghe.acme.test",
				"webhook_allowed_addresses": "192.0.2.0/24",
				"verify_server_cert": false,
				"authentication_strategy": "oauth2",
				"proxy_auth_strategy": "bearer"
			},
			"rate_limit": {
				"limit": 5000,
				"used": 42,
				"remaining": 4958,
				"reset_at": "2026-07-16T06:00:00Z"
			}
		}`)
	})

	got, _, err := client.RepositoryConnections.Get(context.Background(), "acme-inc", id)
	if err != nil {
		t.Fatalf("RepositoryConnections.Get returned error: %v", err)
	}

	webhookAllowedAddresses := "192.0.2.0/24"
	verifyServerCert := false
	authenticationStrategy := "oauth2"
	proxyAuthStrategy := "bearer"
	want := RepositoryConnection{
		ID:          id,
		Type:        "github_enterprise_app",
		DisplayName: "GitHub Enterprise Server (acme-ghe)",
		URL:         "https://api.buildkite.com/v2/organizations/acme-inc/repository_connections/" + id,
		ServiceAccount: &RepositoryConnectionServiceAccount{
			Login: "acme-ghe",
		},
		Host: &RepositoryConnectionHost{
			Type:                    "github_enterprise_app",
			URL:                     "https://ghe.acme.test",
			WebhookAllowedAddresses: &webhookAllowedAddresses,
			VerifyServerCert:        &verifyServerCert,
			AuthenticationStrategy:  &authenticationStrategy,
			ProxyAuthStrategy:       &proxyAuthStrategy,
		},
		RateLimit: &RepositoryConnectionRateLimit{
			Limit:     5000,
			Used:      42,
			Remaining: 4958,
			ResetAt:   time.Date(2026, time.July, 16, 6, 0, 0, 0, time.UTC),
		},
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("RepositoryConnections.Get diff (-want +got):\n%s", diff)
	}
}

func TestRepositoryConnectionsService_GetNullDetails(t *testing.T) {
	t.Parallel()

	server, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	const id = "12345678-9abc-def0-1234-56789abcdef0"
	server.HandleFunc("/v2/organizations/acme-inc/repository_connections/"+id, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = fmt.Fprint(w, `{
			"id": "12345678-9abc-def0-1234-56789abcdef0",
			"type": "cursor_origin_app",
			"display_name": "Origin",
			"url": "https://api.buildkite.com/v2/organizations/acme-inc/repository_connections/12345678-9abc-def0-1234-56789abcdef0",
			"service_account": null,
			"host": null,
			"rate_limit": null
		}`)
	})

	got, _, err := client.RepositoryConnections.Get(context.Background(), "acme-inc", id)
	if err != nil {
		t.Fatalf("RepositoryConnections.Get returned error: %v", err)
	}
	if got.ServiceAccount != nil || got.Host != nil || got.RateLimit != nil {
		t.Errorf("nullable details = (%v, %v, %v), want all nil", got.ServiceAccount, got.Host, got.RateLimit)
	}
}
