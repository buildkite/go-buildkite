package buildkite

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

// APIRequestLogExportsService handles communication with the API request log
// export related methods of the Buildkite API.
//
// API request log exports are in private preview and must be enabled for the
// organization. Requests need an organization administrator's API access
// token with the read_api_request_logs scope.
//
// Buildkite API docs: https://buildkite.com/docs/apis/rest-api/organizations/api-request-log-exports
type APIRequestLogExportsService struct {
	client *Client
}

// APIRequestLogExport is an export of an organization's REST or GraphQL API
// request logs. State is "pending", "running", "completed" or "failed".
// DownloadURL is set once the export completes and until it expires.
type APIRequestLogExport struct {
	UUID          string     `json:"uuid,omitempty"`
	URL           string     `json:"url,omitempty"`
	Dataset       string     `json:"dataset,omitempty"`
	Format        string     `json:"format,omitempty"`
	State         string     `json:"state,omitempty"`
	From          *Timestamp `json:"from,omitempty"`
	To            *Timestamp `json:"to,omitempty"`
	RowCount      *int       `json:"row_count,omitempty"`
	DownloadURL   *string    `json:"download_url,omitempty"`
	ErrorMessage  *string    `json:"error_message,omitempty"`
	CreatedByUUID string     `json:"created_by_uuid,omitempty"`
	CreatedAt     *Timestamp `json:"created_at,omitempty"`
	StartedAt     *Timestamp `json:"started_at,omitempty"`
	FinishedAt    *Timestamp `json:"finished_at,omitempty"`
	ExpiresAt     *Timestamp `json:"expires_at,omitempty"`
}

// APIRequestLogExportCreate is the request body for creating an export.
//
// Dataset is "rest" or "graphql". Format is "parquet" (the API default),
// "csv" or "ndjson". From is inclusive and To is exclusive; the period must
// cover one to 24 complete UTC hours within the last 30 days. When omitted,
// To defaults to today's UTC midnight and From to 24 hours before To.
type APIRequestLogExportCreate struct {
	Dataset string     `json:"dataset"`
	Format  string     `json:"format,omitempty"`
	From    *Timestamp `json:"from,omitempty"`
	To      *Timestamp `json:"to,omitempty"`
}

// Create requests a new API request log export. Poll it with Get until its
// state is "completed" or "failed".
func (s *APIRequestLogExportsService) Create(ctx context.Context, org string, ec APIRequestLogExportCreate) (APIRequestLogExport, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/api_request_log_exports", org)
	req, err := s.client.NewRequest(ctx, http.MethodPost, u, ec)
	if err != nil {
		return APIRequestLogExport{}, nil, err
	}

	var export APIRequestLogExport
	resp, err := s.client.Do(req, &export)
	if err != nil {
		return APIRequestLogExport{}, resp, err
	}

	return export, resp, nil
}

// Get returns an API request log export by UUID.
func (s *APIRequestLogExportsService) Get(ctx context.Context, org, id string) (APIRequestLogExport, *Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/api_request_log_exports/%s", org, id)
	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return APIRequestLogExport{}, nil, err
	}

	var export APIRequestLogExport
	resp, err := s.client.Do(req, &export)
	if err != nil {
		return APIRequestLogExport{}, resp, err
	}

	return export, resp, nil
}

// Download writes a completed API request log export's file to w. The API
// redirects to a short-lived signed link, which the HTTP client follows.
// Parquet files are zstd-compressed; CSV and NDJSON files are gzip-compressed.
func (s *APIRequestLogExportsService) Download(ctx context.Context, org, id string, w io.Writer) (*Response, error) {
	u := fmt.Sprintf("v2/organizations/%s/api_request_log_exports/%s/download", org, id)
	req, err := s.client.NewRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, w)
}
