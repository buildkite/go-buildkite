package buildkite

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type httpCall struct {
	Method string
	Path   string
}

type mockServer struct {
	mux   *http.ServeMux
	calls []httpCall
}

func (ms *mockServer) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	ms.mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		ms.calls = append(ms.calls, httpCall{
			Method: r.Method,
			Path:   r.URL.Path,
		})
		handler(w, r)
	})
}

// newMockServerAndClient sets up a test HTTP server along with a buildkite.Client that is
// configured to talk to that test server.  Tests should register handlers on
// mux which provide mock responses for the API method being tested.
func newMockServerAndClient(t *testing.T) (*mockServer, *Client, func()) {
	// test server
	mux := http.NewServeMux()
	// Fail test if unexpected request is received, "/" matches any request not matched by a more specific handler
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("unexpected %s request for %s", r.Method, r.URL.Path)
	})

	server := httptest.NewServer(mux)

	client, err := NewOpts(WithBaseURL(server.URL))
	if err != nil {
		t.Fatalf("unexpected NewOpts() error: %v", err)
	}

	ms := &mockServer{mux: mux, calls: make([]httpCall, 0, 10)}

	return ms, client, func() { server.Close() }
}

func testMethod(t *testing.T, r *http.Request, want string) {
	if got := r.Method; got != want {
		t.Errorf("Request method: %v, want %v", got, want)
	}
}

type (
	values     map[string]string
	valuesList []struct{ key, val string }
)

func testFormValues(t *testing.T, r *http.Request, values values) {
	want := url.Values{}
	for k, v := range values {
		want.Add(k, v)
	}

	err := r.ParseForm()
	if err != nil {
		t.Fatalf("parsing HTTP form body: %v", err)
	}

	if diff := cmp.Diff(r.Form, want); diff != "" {
		t.Errorf("Request parameters diff: (-got +want)\n%s", diff)
	}
}

func testFormValuesList(t *testing.T, r *http.Request, values valuesList) {
	want := url.Values{}
	for _, v := range values {
		want.Add(v.key, v.val)
	}

	err := r.ParseForm()
	if err != nil {
		t.Fatalf("parsing HTTP form body: %v", err)
	}

	if diff := cmp.Diff(r.Form, want); diff != "" {
		t.Errorf("Request parameters diff: (-got +want)\n%s", diff)
	}
}

func TestNewClient(t *testing.T) {
	c, err := NewClient()
	if err != nil {
		t.Fatalf("unexpected NewClient() error: %v", err)
	}

	if got, want := c.BaseURL.String(), DefaultBaseURL; got != want {
		t.Errorf("NewClient BaseURL is %v, want %v", got, want)
	}
	if got, want := c.UserAgent, DefaultUserAgent; got != want {
		t.Errorf("NewClient UserAgent is %v, want %v", got, want)
	}
}

func TestNewRequest(t *testing.T) {
	c, err := NewClient()
	if err != nil {
		t.Fatalf("unexpected NewClient() error: %v", err)
	}

	inURL, outURL := "/foo", DefaultBaseURL+"foo"
	inBody := User{ID: "123", Name: "Jane Doe", Email: "jane@doe.com"}
	outBody := `{"id":"123","name":"Jane Doe","email":"jane@doe.com"}` + "\n"

	req, _ := c.NewRequest(context.Background(), "GET", inURL, inBody)

	// test that relative URL was expanded
	if got, want := req.URL.String(), outURL; got != want {
		t.Errorf("NewRequest(%q) URL is %v, want %v", inURL, got, want)
	}

	// test that body was JSON encoded
	body, _ := io.ReadAll(req.Body)
	if got, want := string(body), outBody; got != want {
		t.Errorf("NewRequest(%v) Body is %v, want %v", inBody, got, want)
	}

	// test that content-type said it was JSON too
	if got, want := req.Header.Get("Content-Type"), "application/json"; got != want {
		t.Errorf("NewRequest() Content-Type is %v, want %v", got, want)
	}

	// test that default user-agent is attached to the request
	if got, want := req.Header.Get("User-Agent"), c.UserAgent; got != want {
		t.Errorf("NewRequest() User-Agent is %v, want %v", got, want)
	}
}

func TestNewRequest_WhenTokenAuthIsConfigured_AddsBearerTokenToHeaders(t *testing.T) {
	c, err := NewOpts(WithTokenAuth("hunter2"))
	if err != nil {
		t.Fatalf("unexpected NewOpts() error: %v", err)
	}
	req, _ := c.NewRequest(context.Background(), "GET", "/foo", nil)

	if got, want := req.Header.Get("Authorization"), "Bearer hunter2"; got != want {
		t.Errorf("NewRequest() Authorization is %v, want %v", got, want)
	}
}

func TestClient_resolveURL(t *testing.T) {
	tests := []struct {
		name        string
		baseURL     string
		relPath     string
		expectedURL string
		expectError bool
	}{
		// Backward compatibility - no path prefix
		{
			name:        "standard base URL without path",
			baseURL:     "https://api.buildkite.com/",
			relPath:     "v2/organizations/myorg/pipelines",
			expectedURL: "https://api.buildkite.com/v2/organizations/myorg/pipelines",
		},
		{
			name:        "standard base URL without trailing slash",
			baseURL:     "https://api.buildkite.com",
			relPath:     "v2/organizations/myorg/pipelines",
			expectedURL: "https://api.buildkite.com/v2/organizations/myorg/pipelines",
		},
		// New functionality - with path prefix
		{
			name:        "base URL with path prefix",
			baseURL:     "https://proxy.example.com/api/",
			relPath:     "v2/organizations/myorg/pipelines",
			expectedURL: "https://proxy.example.com/api/v2/organizations/myorg/pipelines",
		},
		{
			name:        "base URL with path prefix no trailing slash",
			baseURL:     "https://proxy.example.com/api",
			relPath:     "v2/organizations/myorg/pipelines",
			expectedURL: "https://proxy.example.com/api/v2/organizations/myorg/pipelines",
		},
		{
			name:        "base URL with deep path prefix",
			baseURL:     "https://gateway.example.com/internal/buildkite/",
			relPath:     "v2/organizations/myorg/builds",
			expectedURL: "https://gateway.example.com/internal/buildkite/v2/organizations/myorg/builds",
		},
		// Edge cases
		{
			name:        "relative path with leading slash",
			baseURL:     "https://proxy.example.com/api/",
			relPath:     "/v2/organizations/myorg/pipelines",
			expectedURL: "https://proxy.example.com/api/v2/organizations/myorg/pipelines",
		},
		{
			name:        "base URL with root path only",
			baseURL:     "https://proxy.example.com/",
			relPath:     "v2/organizations/myorg/pipelines",
			expectedURL: "https://proxy.example.com/v2/organizations/myorg/pipelines",
		},
		{
			name:        "empty relative path",
			baseURL:     "https://proxy.example.com/api/",
			relPath:     "",
			expectedURL: "https://proxy.example.com/api/",
		},
		// Query parameters and fragments
		{
			name:        "relative path with query parameters",
			baseURL:     "https://proxy.example.com/api/",
			relPath:     "v2/organizations/myorg/pipelines?page=2&per_page=50",
			expectedURL: "https://proxy.example.com/api/v2/organizations/myorg/pipelines?page=2&per_page=50",
		},
		{
			name:        "relative path with fragment",
			baseURL:     "https://proxy.example.com/api/",
			relPath:     "v2/organizations/myorg/pipelines#section",
			expectedURL: "https://proxy.example.com/api/v2/organizations/myorg/pipelines#section",
		},
		// Absolute URLs should be used as-is
		{
			name:        "absolute URL should be preserved",
			baseURL:     "https://proxy.example.com/api/",
			relPath:     "https://other.example.com/v2/test",
			expectedURL: "https://other.example.com/v2/test",
		},
		// Error cases
		{
			name:        "invalid relative path",
			baseURL:     "https://proxy.example.com/api/",
			relPath:     "://invalid-url",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewOpts(WithBaseURL(tt.baseURL))
			if err != nil {
				t.Fatalf("NewOpts() error = %v", err)
			}

			result, err := client.resolveURL(tt.relPath)

			if tt.expectError {
				if err == nil {
					t.Errorf("resolveURL() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("resolveURL() error = %v", err)
				return
			}

			if got := result.String(); got != tt.expectedURL {
				t.Errorf("resolveURL() = %v, want %v", got, tt.expectedURL)
			}
		})
	}
}

func TestNewRequest_WithPathPrefix(t *testing.T) {
	tests := []struct {
		name        string
		baseURL     string
		relPath     string
		expectedURL string
	}{
		{
			name:        "standard base URL maintains compatibility",
			baseURL:     DefaultBaseURL,
			relPath:     "v2/organizations/testorg/pipelines",
			expectedURL: DefaultBaseURL + "v2/organizations/testorg/pipelines",
		},
		{
			name:        "base URL with path prefix includes prefix",
			baseURL:     "https://proxy.example.com/api/",
			relPath:     "v2/organizations/testorg/pipelines",
			expectedURL: "https://proxy.example.com/api/v2/organizations/testorg/pipelines",
		},
		{
			name:        "base URL with deep path prefix",
			baseURL:     "https://gateway.example.com/internal/buildkite/",
			relPath:     "v2/organizations/testorg/builds/123",
			expectedURL: "https://gateway.example.com/internal/buildkite/v2/organizations/testorg/builds/123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewOpts(WithBaseURL(tt.baseURL))
			if err != nil {
				t.Fatalf("NewOpts() error = %v", err)
			}

			req, err := client.NewRequest(context.Background(), "GET", tt.relPath, nil)
			if err != nil {
				t.Fatalf("NewRequest() error = %v", err)
			}

			if got := req.URL.String(); got != tt.expectedURL {
				t.Errorf("NewRequest() URL = %v, want %v", got, tt.expectedURL)
			}
		})
	}
}

func TestResponse_populatePageValues(t *testing.T) {
	r := http.Response{
		Header: http.Header{
			"Link": {
				`<https://api.buildkite.com/?page=1>; rel="first",` +
					` <https://api.buildkite.com/?page=2>; rel="prev",` +
					` <https://api.buildkite.com/?page=4>; rel="next",` +
					` <https://api.buildkite.com/?page=5>; rel="last"`,
			},
		},
	}

	response := newResponse(&r)
	if got, want := response.FirstPage, 1; got != want {
		t.Errorf("response.FirstPage: %v, want %v", got, want)
	}
	if got, want := response.PrevPage, 2; want != got {
		t.Errorf("response.PrevPage: %v, want %v", got, want)
	}
	if got, want := response.NextPage, 4; want != got {
		t.Errorf("response.NextPage: %v, want %v", got, want)
	}
	if got, want := response.LastPage, 5; want != got {
		t.Errorf("response.LastPage: %v, want %v", got, want)
	}
}
