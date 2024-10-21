package buildkite

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/buildkite/go-buildkite/v3/internal/bkmultipart"
	"github.com/cenkalti/backoff"
	"github.com/google/go-querystring/query"
)

const (
	DefaultBaseURL   = "https://api.buildkite.com/"
	DefaultUserAgent = "go-buildkite/" + Version
)

// Client - A Client manages communication with the buildkite API.
type Client struct {
	// HTTP client used to communicate with the API.
	client *http.Client

	// Base URL for API requests.  Defaults to the public buildkite API. BaseURL should
	// always be specified with a trailing slash.
	BaseURL *url.URL

	// User agent used when communicating with the buildkite API.
	UserAgent string

	// Services used for talking to different parts of the buildkite API.
	AccessTokens             *AccessTokensService
	Agents                   *AgentsService
	Annotations              *AnnotationsService
	Artifacts                *ArtifactsService
	Builds                   *BuildsService
	Clusters                 *ClustersService
	ClusterQueues            *ClusterQueuesService
	ClusterTokens            *ClusterTokensService
	FlakyTests               *FlakyTestsService
	Jobs                     *JobsService
	Organizations            *OrganizationsService
	PackagesService          *PackagesService
	PackageRegistriesService *PackageRegistriesService
	Pipelines                *PipelinesService
	PipelineTemplates        *PipelineTemplatesService
	User                     *UserService
	Teams                    *TeamsService
	Tests                    *TestsService
	TestRuns                 *TestRunsService
	TestSuites               *TestSuitesService

	authHeader string
	httpDebug  bool
}

type clientOpt func(*Client) error

// WithHTTPClient configures the buildkite.Client to use the provided http.Client. This can be used to
// customise the client's transport (e.g. to use a custom TLS configuration) or to provide a mock client
func WithHTTPClient(client *http.Client) clientOpt {
	return func(c *Client) error {
		c.client = client
		return nil
	}
}

// WithBaseURL configures the buildkite.Client to use the provided base URL, instead of the default of https://api.buildkite.com/
func WithBaseURL(baseURL string) clientOpt {
	return func(c *Client) error {
		var err error
		c.BaseURL, err = url.Parse(baseURL)
		if err != nil {
			return fmt.Errorf("failed to parse baseURL: %w", err)
		}

		return nil
	}
}

// WithUserAgent configures the buildkite.Client to use the provided user agent string, instead of the default of "go-buildkite/<version>"
func WithUserAgent(userAgent string) clientOpt {
	return func(c *Client) error {
		c.UserAgent = userAgent
		return nil
	}
}

// WithTokenAuth configures the buildkite.Client to use the provided token for authentication.
// This is the recommended way to authenticate with the buildkite API
// Note that at least one of [WithTokenAuth] or [WithBasicAuth] must be provided to NewOpts
func WithTokenAuth(token string) clientOpt {
	return func(c *Client) error {
		c.authHeader = fmt.Sprintf("Bearer %s", token)
		return nil
	}
}

// WithHTTPDebug configures the buildkite.Client to print debug information about HTTP requests and responses as it makes them
func WithHTTPDebug(debug bool) clientOpt {
	return func(c *Client) error {
		c.httpDebug = debug
		return nil
	}
}

// NewClient returns a new buildkite API client with the provided options.
// Note that at [WithTokenAuth] must be provided for requests to the buildkite API to succeed.
// Otherwise, sensible defaults are used.
func NewClient(opts ...clientOpt) (*Client, error) {
	baseURL, _ := url.Parse(DefaultBaseURL)

	c := &Client{
		client:    http.DefaultClient,
		BaseURL:   baseURL,
		UserAgent: DefaultUserAgent,
	}

	for _, opt := range opts {
		err := opt(c)
		if err != nil {
			return nil, fmt.Errorf("failed to apply client option: %w", err)
		}
	}

	c.populateDefaultServices()

	return c, nil
}

// NewOpts is an alias for NewClient
func NewOpts(opts ...clientOpt) (*Client, error) {
	return NewClient(opts...)
}

func (c *Client) populateDefaultServices() {
	c.AccessTokens = &AccessTokensService{c}
	c.Agents = &AgentsService{c}
	c.Annotations = &AnnotationsService{c}
	c.Artifacts = &ArtifactsService{c}
	c.Builds = &BuildsService{c}
	c.Clusters = &ClustersService{c}
	c.ClusterQueues = &ClusterQueuesService{c}
	c.ClusterTokens = &ClusterTokensService{c}
	c.FlakyTests = &FlakyTestsService{c}
	c.Jobs = &JobsService{c}
	c.Organizations = &OrganizationsService{c}
	c.PackagesService = &PackagesService{c}
	c.PackageRegistriesService = &PackageRegistriesService{c}
	c.Pipelines = &PipelinesService{c}
	c.PipelineTemplates = &PipelineTemplatesService{c}
	c.User = &UserService{c}
	c.Teams = &TeamsService{c}
	c.Tests = &TestsService{c}
	c.TestRuns = &TestRunsService{c}
	c.TestSuites = &TestSuitesService{c}
}

// NewRequest creates an API request. A relative URL can be provided in urlStr,
// in which case it is resolved relative to the BaseURL of the Client.
// Relative URLs should always be specified without a preceding slash.  If
// specified, the value pointed to by body is JSON encoded and included as the
// request body.
func (c *Client) NewRequest(ctx context.Context, method, urlStr string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)

	headers := map[string]string{}

	if c.authHeader != "" {
		headers["Authorization"] = c.authHeader
	}

	if c.UserAgent != "" {
		headers["User-Agent"] = c.UserAgent
	}

	var reqBody io.Reader
	if body != nil {
		switch v := body.(type) {
		case *bkmultipart.Streamer:
			panic("bkmultipart.Streamer passed directly to NewRequest. Did you mean to pass bkstreamer.Streamer.Reader() instead?")

		case io.Reader: // If body is an io.Reader, use it directly, the caller is responsible for encoding
			reqBody = v

		default: // Otherwise, encode it as JSON
			buf := &bytes.Buffer{}
			headers["Content-Type"] = "application/json"

			err := json.NewEncoder(buf).Encode(body)
			if err != nil {
				return nil, err
			}

			reqBody = buf
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), reqBody)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	return req, nil
}

// Response is a buildkite API response.  This wraps the standard http.Response
// returned from buildkite and provides convenient access to things like
// pagination links.
type Response struct {
	*http.Response

	// These fields provide the page values for paginating through a set of
	// results.  Any or all of these may be set to the zero value for
	// responses that are not part of a paginated set, or for which there
	// are no additional pages.

	NextPage  int
	PrevPage  int
	FirstPage int
	LastPage  int
}

// newResponse creats a new Response for the provided http.Response.
func newResponse(r *http.Response) *Response {
	response := &Response{Response: r}
	response.populatePageValues()
	return response
}

// populatePageValues parses the HTTP Link response headers and populates the
// various pagination link values in the Reponse.
func (r *Response) populatePageValues() {
	if links, ok := r.Response.Header["Link"]; ok && len(links) > 0 {
		for _, link := range strings.Split(links[0], ",") {
			segments := strings.Split(strings.TrimSpace(link), ";")

			// link must at least have href and rel
			if len(segments) < 2 {
				continue
			}

			// ensure href is properly formatted
			if !strings.HasPrefix(segments[0], "<") || !strings.HasSuffix(segments[0], ">") {
				continue
			}

			// try to pull out page parameter
			url, err := url.Parse(segments[0][1 : len(segments[0])-1])
			if err != nil {
				continue
			}
			page := url.Query().Get("page")
			if page == "" {
				continue
			}

			for _, segment := range segments[1:] {
				switch strings.TrimSpace(segment) {
				case `rel="next"`:
					r.NextPage, _ = strconv.Atoi(page)
				case `rel="prev"`:
					r.PrevPage, _ = strconv.Atoi(page)
				case `rel="first"`:
					r.FirstPage, _ = strconv.Atoi(page)
				case `rel="last"`:
					r.LastPage, _ = strconv.Atoi(page)
				}

			}
		}
	}
}

// Do sends an API request and returns the API response.  The API response is
// JSON decoded and stored in the value pointed to by v, or returned as an
// error if an API error has occurred.  If v implements the io.Writer
// interface, the raw response body will be written to v, without attempting to
// first decode it.
func (c *Client) Do(req *http.Request, v interface{}) (*Response, error) {
	respCh := make(chan *http.Response, 1)

	op := func() error {
		if c.httpDebug {
			if dump, err := httputil.DumpRequest(req, true); err == nil {
				fmt.Printf("DEBUG request uri=%s\n%s\n", req.URL, dump)
			}
		}

		resp, err := c.client.Do(req)
		if err != nil {
			return backoff.Permanent(err)
		}

		if c.httpDebug {
			if dump, err := httputil.DumpResponse(resp, true); err == nil {
				fmt.Printf("DEBUG response uri=%s\n%s\n", req.URL, dump)
			}
		}

		// Check for rate limiting response on idempotent requests
		if req.Method == http.MethodGet && resp.StatusCode == http.StatusTooManyRequests {
			errMsg := resp.Header.Get("Rate-Limit-Warning")
			if errMsg == "" {
				errMsg = "Too many requests, retry"
			}
			return errors.New(errMsg)
		}

		respCh <- resp
		return nil
	}

	notify := func(err error, delay time.Duration) {
		if c.httpDebug {
			fmt.Printf("DEBUG error %v, retry in %v\n", err, delay)
		}
	}

	if err := backoff.RetryNotify(op, backoff.NewExponentialBackOff(), notify); err != nil {
		return nil, err
	}

	resp := <-respCh

	defer resp.Body.Close()
	defer func() { _, _ = io.Copy(io.Discard, resp.Body) }()

	response := newResponse(resp)

	if err := checkResponse(resp); err != nil {
		// even though there was an error, we still return the response
		// in case the caller wants to inspect it further
		return response, err
	}

	var err error

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
		}
	}

	return response, err
}

// ErrorResponse provides a message.
type ErrorResponse struct {
	Response *http.Response // HTTP response that caused this error
	Message  string         `json:"message"` // error message
	RawBody  []byte         `json:"-"`       // Raw Response Body
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %v",
		r.Response.Request.Method, r.Response.Request.URL,
		r.Response.StatusCode, r.Message)
}

func checkResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}

	errorResponse := &ErrorResponse{Response: r}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("response failed with error %w, but reading response body failed with error %w", errorResponse, err)
	}
	errorResponse.RawBody = data

	err = json.Unmarshal(data, errorResponse)
	if err != nil {
		return fmt.Errorf("response failed with error %w, but parsing response body JSON failed with error: %w. Raw body of error was: %s", errorResponse, err, string(data))
	}
	return errorResponse
}

// addOptions adds the parameters in opt as URL query parameters to s.  opt
// must be a struct whose fields may contain "url" tags.
func addOptions(s string, opt interface{}) (string, error) {
	v := reflect.ValueOf(opt)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	qs, err := query.Values(opt)
	if err != nil {
		return s, err
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}

// ListOptions specifies the optional parameters to various List methods that
// support pagination.
type ListOptions struct {
	// For paginated result sets, page of results to retrieve.
	Page int `url:"page,omitempty"`

	// For paginated result sets, the number of results to include per page.
	PerPage int `url:"per_page,omitempty"`
}
