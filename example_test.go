package buildkite_test

import (
	"context"

	"github.com/prateek/go-buildkite"
)

// Example_clientOptComposition demonstrates how to programmatically compose client options
// based on conditional logic, leveraging the exported ClientOpt type.
func Example_clientOptComposition() {
	// Simulating command-line flags or config values
	baseURL := "https://api.buildkite.com/"
	token := "your-token"
	userAgentFlag := "custom-agent" // Could be empty in some cases
	debugEnabled := false

	// Build options programmatically
	var opts []buildkite.ClientOpt

	// Always add required options
	opts = append(opts, buildkite.WithTokenAuth(token))

	// Conditionally add other options
	if baseURL != buildkite.DefaultBaseURL {
		opts = append(opts, buildkite.WithBaseURL(baseURL))
	}

	if userAgentFlag != "" {
		opts = append(opts, buildkite.WithUserAgent(userAgentFlag))
	}

	if debugEnabled {
		opts = append(opts, buildkite.WithHTTPDebug(true))
	}

	// Create client with composed options
	client, err := buildkite.NewClient(opts...)
	if err != nil {
		// Handle error
		return
	}

	// Use the configured client
	_, _ = client.User.Get(context.Background())

	// Output:
}
