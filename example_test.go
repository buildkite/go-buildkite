package buildkite_test

import (
	"github.com/buildkite/go-buildkite/v4"
)

// Example_clientOptComposition demonstrates how to programmatically compose client options
// based on conditional logic, leveraging the exported ClientOpt type.
func Example_clientOptComposition() {
	// This is a test example only - we don't actually run the client
	// Simulating command-line flags or config values
	baseURL := "https://api.buildkite.com/"
	token := "your-token"
	userAgentFlag := "custom-agent" // Could be empty in some cases
	debugEnabled := false

	// Build options programmatically - we're testing the pattern more than actual execution
	var opts []buildkite.ClientOpt

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

	// Always add required options
	opts = append(opts, buildkite.WithTokenAuth(token))

	// Create the client using the composed options
	client, err := buildkite.NewClient(opts...)
	if err != nil {
		// In a real application, you'd handle this error appropriately
		// For this example, we'll just acknowledge the client was created
		_ = client
	}
}
