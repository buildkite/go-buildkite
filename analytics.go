package buildkite

// The models for the Test Engine (Analytics) endpoints — suites, tests, runs,
// failed executions, and build tests — are generated from the Test Engine
// OpenAPI description into analytics_gen.go. Edit analytics-openapi.yaml, or
// the codegen config and its overlay, rather than the generated structs.
//
// Endpoints the description does not cover, such as flaky tests, keep their
// hand-written models (see flaky_tests.go).
//
//go:generate go tool oapi-codegen -config oapi-codegen-analytics.yaml analytics-openapi.yaml
