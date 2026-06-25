// Package buildkite provides a client for the Buildkite API.
package buildkite

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"net/http/httputil"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/buildkite/go-buildkite/v5/internal/bkmultipart"
	"github.com/buildkite/roko"
	"github.com/google/go-querystring/query"
)

const (
	DefaultBaseURL    = "https://api.buildkite.com/"
	DefaultMaxRetries = 3
)

// DefaultUserAgent is the default user agent used for API requests
var DefaultUserAgent = "go-buildkite/" + Version

// errRateLimited is returned from the roko callback to signal a retryable 429.
// It is never surfaced to callers — checkResponse produces the real *ErrorResponse.
var errRateLimited = errors.New("rate limited")

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
	BuildTests               *BuildTestsService
	Builds                   *BuildsService
	Clusters                 *ClustersService
	ClusterQueues            *ClusterQueuesService
	ClusterTokens            *ClusterTokensService
	ClusterSecrets           *ClusterSecretsService
	ClusterMaintainers       *ClusterMaintainersService
	Emojis                   *EmojisService
	FlakyTests               *FlakyTestsService
	Jobs                     *JobsService
	Members                  *MembersService
	Meta                     *MetaService
	Organizations            *OrganizationsService
	PackagesService          *PackagesService
	PackageRegistriesService *PackageRegistriesService
	Pipelines                *PipelinesService
	PipelineSchedules        *PipelineSchedulesService
	PipelineTemplates        *PipelineTemplatesService
	RateLimit                *RateLimitService
	Rules                    *RulesService
	User                     *UserService
	Teams                    *TeamsService
	TeamMember               *TeamMemberService
	TeamPipelines            *TeamPipelinesService
	TeamSuites               *TeamSuitesService
	Tests                    *TestsService
	TestRuns                 *TestRunsService
	TestSuites               *TestSuitesService

	authHeader      string
	httpDebug       bool
	rateLimitNotify RateLimitNotify
	maxRetries      int
	sleepFunc       func(time.Duration) // test-only override for retry backoff; nil uses roko's default
}

// RateLimitNotify is called each time a 429 response is received, including on
// the final exhausted attempt where no retry follows and when maxRetries is 0.
// attempt is 1-based; delay is the back-off before the next retry, or 0 if
// this is the final attempt and no sleep will follow.
type RateLimitNotify func(attempt int, delay time.Duration)

// ClientOpt is a function that configures a Client.
type ClientOpt func(*Client) error

// clientOpt is deprecated and will be removed in a future version.
// It is an alias for ClientOpt for backward compatibility.
//
//nolint:unused // used to ensure no major release needed
type clientOpt = ClientOpt

// WithHTTPClient configures the buildkite.Client to use the provided http.Client. This can be used to
// customise the client's transport (e.g. to use a custom TLS configuration) or to provide a mock client
func WithHTTPClient(client *http.Client) ClientOpt {
	return func(c *Client) error {
		c.client = client
		return nil
	}
}

// WithBaseURL configures the buildkite.Client to use the provided base URL, instead of the default of https://api.buildkite.com/
func WithBaseURL(baseURL string) ClientOpt {
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
func WithUserAgent(userAgent string) ClientOpt {
	return func(c *Client) error {
		c.UserAgent = userAgent
		return nil
	}
}

// WithTokenAuth configures the buildkite.Client to use the provided token for authentication.
// This is the recommended way to authenticate with the buildkite API
// Note that at least one of [WithTokenAuth] or [WithBasicAuth] must be provided to NewOpts
func WithTokenAuth(token string) ClientOpt {
	return func(c *Client) error {
		c.authHeader = fmt.Sprintf("Bearer %s", token)
		return nil
	}
}

// WithHTTPDebug configures the buildkite.Client to print debug information about HTTP requests and responses as it makes them
func WithHTTPDebug(debug bool) ClientOpt {
	return func(c *Client) error {
		c.httpDebug = debug
		return nil
	}
}

// WithRateLimitNotify registers a callback invoked on every 429 response,
// including when retries are exhausted or disabled (maxRetries=0). Use it to
// log, emit metrics, or display progress. The callback is not invoked when the
// request body is a raw io.Reader without GetBody set, since those requests
// cannot be retried and the 429 is surfaced directly as an *ErrorResponse.
// Passing nil clears any previously registered callback.
func WithRateLimitNotify(fn RateLimitNotify) ClientOpt {
	return func(c *Client) error {
		c.rateLimitNotify = fn
		return nil
	}
}

// WithMaxRetries sets the maximum number of retry attempts on rate-limited requests.
// Defaults to DefaultMaxRetries (3). Use 0 to disable retries entirely.
//
// There is no internal wall-clock time limit. With a high retry count and a
// server consistently returning RateLimit-Reset: 120, Do() can block for many
// minutes. Callers should set a context deadline to bound total wait time.
//
// All HTTP methods are retried, including POST, PUT, and DELETE, provided the
// request body is rewindable (i.e. created via NewRequest with a struct or
// bytes.Buffer body). Callers that cannot tolerate duplicate side-effects on
// non-idempotent requests should pass a context with an appropriate deadline
// or set WithMaxRetries(0) for those specific calls.
func WithMaxRetries(n int) ClientOpt {
	return func(c *Client) error {
		if n < 0 {
			return fmt.Errorf("max retries must be >= 0, got %d", n)
		}
		c.maxRetries = n
		return nil
	}
}

// NewClient returns a new buildkite API client with the provided options.
// Note that at [WithTokenAuth] must be provided for requests to the buildkite API to succeed.
// Otherwise, sensible defaults are used.
func NewClient(opts ...ClientOpt) (*Client, error) {
	baseURL, _ := url.Parse(DefaultBaseURL)

	c := &Client{
		client:     http.DefaultClient,
		BaseURL:    baseURL,
		UserAgent:  DefaultUserAgent,
		maxRetries: DefaultMaxRetries,
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
func NewOpts(opts ...ClientOpt) (*Client, error) {
	return NewClient(opts...)
}

func (c *Client) populateDefaultServices() {
	c.AccessTokens = &AccessTokensService{c}
	c.Agents = &AgentsService{c}
	c.Annotations = &AnnotationsService{c}
	c.Artifacts = &ArtifactsService{c}
	c.BuildTests = &BuildTestsService{c}
	c.Builds = &BuildsService{c}
	c.Clusters = &ClustersService{c}
	c.ClusterQueues = &ClusterQueuesService{c}
	c.ClusterTokens = &ClusterTokensService{c}
	c.ClusterSecrets = &ClusterSecretsService{c}
	c.ClusterMaintainers = &ClusterMaintainersService{c}
	c.Emojis = &EmojisService{c}
	c.FlakyTests = &FlakyTestsService{c}
	c.Jobs = &JobsService{c}
	c.Members = &MembersService{c}
	c.Meta = &MetaService{c}
	c.Organizations = &OrganizationsService{c}
	c.PackagesService = &PackagesService{c}
	c.PackageRegistriesService = &PackageRegistriesService{c}
	c.Pipelines = &PipelinesService{c}
	c.PipelineSchedules = &PipelineSchedulesService{c}
	c.PipelineTemplates = &PipelineTemplatesService{c}
	c.RateLimit = &RateLimitService{c}
	c.Rules = &RulesService{c}
	c.User = &UserService{c}
	c.Teams = &TeamsService{c}
	c.TeamMember = &TeamMemberService{c}
	c.TeamPipelines = &TeamPipelinesService{c}
	c.TeamSuites = &TeamSuitesService{c}
	c.Tests = &TestsService{c}
	c.TestRuns = &TestRunsService{c}
	c.TestSuites = &TestSuitesService{c}
}

// resolveURL properly handles base URLs with path prefixes by combining the base URL
// path with the relative path, unlike url.URL.ResolveReference which treats relative
// paths as starting from the host root.
func (c *Client) resolveURL(relPath string) (*url.URL, error) {
	rel, err := url.Parse(relPath)
	if err != nil {
		return nil, err
	}

	// If the relative path is absolute, use it as-is
	if rel.IsAbs() {
		return rel, nil
	}

	// Create a copy of the base URL
	result := *c.BaseURL

	// Extract just the path part from the relative URL for path combination
	basePath := strings.TrimSuffix(c.BaseURL.Path, "/")
	cleanRelPath := strings.TrimPrefix(rel.Path, "/")

	if basePath == "" {
		result.Path = "/" + cleanRelPath
	} else {
		result.Path = basePath + "/" + cleanRelPath
	}

	// Preserve query parameters and fragment from the relative URL
	if rel.RawQuery != "" {
		result.RawQuery = rel.RawQuery
	}
	if rel.Fragment != "" {
		result.Fragment = rel.Fragment
	}

	return &result, nil
}

// NewRequest creates an API request. A relative URL can be provided in urlStr,
// in which case it is resolved relative to the BaseURL of the Client.
// Relative URLs should always be specified without a preceding slash.  If
// specified, the value pointed to by body is JSON encoded and included as the
// request body.
func (c *Client) NewRequest(ctx context.Context, method, urlStr string, body interface{}) (*http.Request, error) {
	u, err := c.resolveURL(urlStr)
	if err != nil {
		return nil, err
	}

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
	if links, ok := r.Header["Link"]; ok && len(links) > 0 {
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

// retryDelay returns how long to wait before the next attempt after a 429.
// RateLimit-Reset is interpreted as delta-seconds per the Buildkite API spec.
// Falls back to capped exponential backoff when the header is absent or invalid.
// A random jitter of up to 1 second is added on all paths to spread out
// concurrent clients that receive the same reset time.
func retryDelay(resp *http.Response, attempt int) time.Duration {
	if s := resp.Header.Get("RateLimit-Reset"); s != "" {
		if secs, err := strconv.Atoi(s); err == nil && secs >= 0 {
			if secs > 120 {
				secs = 120
			}
			// 500ms buffer ensures we don't hit the API again in the same window
			// when the server sends RateLimit-Reset: 0.
			return time.Duration(secs)*time.Second + 500*time.Millisecond + time.Duration(rand.N(time.Second))
		}
	}
	// Cap the shift to prevent int64 overflow at high attempt counts.
	shift := attempt
	if shift > 5 {
		shift = 5
	}
	d := time.Duration(1<<uint(shift)) * time.Second
	if d > 30*time.Second {
		d = 30 * time.Second
	}
	return d + time.Duration(rand.N(time.Second))
}

// Do sends an API request and returns the API response.  The API response is
// JSON decoded and stored in the value pointed to by v, or returned as an
// error if an API error has occurred.  If v implements the io.Writer
// interface, the raw response body will be written to v, without attempting to
// first decode it.
func (c *Client) Do(req *http.Request, v interface{}) (*Response, error) {
	var resp *http.Response

	// roko requires a strategy, but the real delay is always driven by the
	// server response (RateLimit-Reset header or exponential fallback), so
	// Constant(0) is a required placeholder that is always overridden.
	retrierOpts := []roko.RetrierOpt{
		roko.WithMaxAttempts(c.maxRetries + 1),
		roko.WithStrategy(roko.Constant(0)),
	}
	if c.sleepFunc != nil {
		retrierOpts = append(retrierOpts, roko.WithSleepFunc(c.sleepFunc))
	}
	retrier := roko.NewRetrier(retrierOpts...)

	rokoErr := retrier.DoWithContext(req.Context(), func(rt *roko.Retrier) error {
		// GetBody is set automatically by http.NewRequestWithContext for
		// bytes.Buffer/bytes.Reader bodies, enabling body replay on retry.
		if rt.AttemptCount() > 0 && req.GetBody != nil {
			body, err := req.GetBody()
			if err != nil {
				rt.Break()
				return fmt.Errorf("rewinding request body for retry: %w", err)
			}
			req.Body = body
		}

		if c.httpDebug {
			if dump, err := httputil.DumpRequest(req, true); err == nil {
				fmt.Printf("DEBUG request uri=%s\n%s\n", req.URL, dump)
			}
		}

		var err error
		resp, err = c.client.Do(req)
		if err != nil {
			resp = nil
			rt.Break()
			return err
		}

		if c.httpDebug {
			if dump, err := httputil.DumpResponse(resp, true); err == nil {
				fmt.Printf("DEBUG response uri=%s\n%s\n", req.URL, dump)
			}
		}

		// canRewind is false for raw io.Reader bodies where GetBody is not set;
		// those requests are not retried since the body cannot be replayed.
		// When canRewind is false and the response is 429, returning nil causes
		// roko to treat the call as successful and exit the loop; execution then
		// falls through to checkResponse which surfaces a proper *ErrorResponse.
		// The caller receives a 429 *ErrorResponse with no WithRateLimitNotify
		// signal and no indication that the retry budget was unused.
		canRewind := req.Body == nil || req.GetBody != nil
		if resp.StatusCode != http.StatusTooManyRequests || !canRewind {
			return nil
		}

		// roko calls MarkAttempt() after this callback returns, so AttemptCount()
		// here is still the 0-based index of the current attempt. The last allowed
		// attempt is index c.maxRetries (= WithMaxAttempts(c.maxRetries+1) - 1).
		var delay time.Duration
		if rt.AttemptCount() < c.maxRetries {
			delay = retryDelay(resp, rt.AttemptCount())
			rt.SetNextInterval(delay)
			// More retries remaining — drain and close so the connection can be reused.
			_, _ = io.Copy(io.Discard, resp.Body)
			_ = resp.Body.Close()
			resp = nil
		}
		// delay is 0 on the final attempt; the callback can use this to distinguish
		// "rate limited and will retry" from "rate limited and giving up".
		if c.rateLimitNotify != nil {
			c.rateLimitNotify(rt.AttemptCount()+1, delay)
		}
		if c.httpDebug {
			fmt.Printf("DEBUG rate limited, retry %d in %v\n", rt.AttemptCount()+1, delay)
		}

		return errRateLimited
	})

	if resp == nil {
		return nil, rokoErr
	}

	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}()

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
